package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Global variables for the WebSocket server
var Broadcast = make(chan ServerMessage)
var users []ServerUser
var categories []ServerCategory
var posts []ServerPost
var clients = make(map[*websocket.Conn]*Session)

// function to start the WebSocket server and go routine for broadcasting messages to all clients
func StartWebSocketServer() {
	upgrader := configureUpgrader()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		cookieValue := r.Header.Get("Cookie")
		fmt.Println(Purple+"Server >> NEW websocket connection with Cookie value:", cookieValue+Reset)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		// Handle the new WebSocket connection
		handleWebSocketConnection(conn, cookieValue)
	})

	// Start the broadcast goroutine
	go func() {
		for {
			message := <-Broadcast
			for client := range clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Println(err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()
}

// function to configure the upgrader for the WebSocket server
func configureUpgrader() *websocket.Upgrader {
	upgrader := &websocket.Upgrader{}
	//Set buffer size for the messages
	upgrader.ReadBufferSize = 1024
	upgrader.WriteBufferSize = 1024
	// Allow requests from all origins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// return the upgrader
	return upgrader
}

// function to handle the WebSocket connection and messages
func handleWebSocketConnection(conn *websocket.Conn, cookieValue string) {
	// SessionsMutex is a mutex to lock the sessions map when adding or removing sessions from it
	var sessionsMutex sync.Mutex
	fmt.Println("New WebSocket Connection", conn.RemoteAddr().String())

	if cookieValue != "" {
		for _, v := range LoggedInUsers {
			if v.Cookie == cookieValue[14:] {
				fmt.Println(Purple+"Server >> User is logged in with cookie value: ", cookieValue+Reset)
				fmt.Println(Purple+"Server >> Updating connection ID for user: ", v.Username+Reset)
				v.WebSocketConn = conn.RemoteAddr().String()
			}
		}
	} else {
		fmt.Println(Purple + "Server >> User is not logged in" + Reset)
	}

	// Create a new session for the WebSocket client
	session := Session{
		WebSocketConn: conn.RemoteAddr().String(),
		UserID:        0,
	}

	// Add the session to the clients map
	// Note: you need to synchronize access to the map using a mutex
	sessionsMutex.Lock()
	clients[conn] = &session
	sessionsMutex.Unlock()

	// Send initial data to the client
	sendInitialData(conn)

	// Listen for new messages from the client
	for {
		var message ServerMessage
		err := conn.ReadJSON(&message)
		//fmt.Println("Message received: ", message.Type)
		if err != nil {
			log.Println(err)
			delete(clients, conn)
			break
		}

		// Handle the message
		handleWebSocketMessage(conn, message)
	}
}

// function to send the initial data to the client when it connects to the server
func sendInitialData(conn *websocket.Conn) {
	// Send the list of users to the new client
	message := ServerMessage{Type: "users", Users: users}
	fmt.Println("Sending initial data to the client")
	conn.WriteJSON(message)

	// Send the list of categories to the new client
	message = ServerMessage{Type: "categories", Categories: categories}
	conn.WriteJSON(message)

	// Send the list of recent posts to the new client
	message = ServerMessage{Type: "posts", Posts: posts}
	conn.WriteJSON(message)
}

// function to handle the messages received from the client
func handleWebSocketMessage(conn *websocket.Conn, message ServerMessage) {
	//fmt.Println("Message received: ", message.Type)

	switch message.Type {
	case "new_user":
		handleNewUserMessage(message)
	case "new_category":
		handleNewCategoryMessage(message)
	case "new_post":
		handleNewPostMessage(message)
	case "get_posts":
		handleGetPostsMessage(conn, message)
	case "get_chat_history":
		handleGetChatHistoryMessage(conn, message)
	case "message":
		handleMessageMessage(conn, message)
	case "login":
		handleLoginMessage(conn, message)
	case "loginResponse":
		handleLoginResponseMessage(conn, message)
	case "logout":
		handleLogoutMessage(conn, message)
	case "register":
		handleRegisterMessage(conn, message)
	case "registerResponse":
		handleRegisterResponseMessage(conn, message)
	case "get_categories":
		handleGetCategoriesMessage(conn, message)
	case "get_users":
		handleGetUsersMessage(conn, message)
	case "get_offline_users":
		handleGetOfflineUsersMessage(conn, message)
	}
}

// function to handle the new user message (when a new user joins the chat
func handleNewUserMessage(message ServerMessage) {
	// Add the new user to the list of users
	users = append(users, ServerUser{Name: message.Users[0].Name})
	// Send the updated list of users to all clients
	Broadcast <- ServerMessage{Type: "users", Users: users}
}

// Send the updated list of categories to all clients
func handleGetCategoriesMessage(conn *websocket.Conn, message ServerMessage) {

	// open the database
	db := OpenDatabase()
	defer CloseDatabase(db)

	// Get all the categories from the database and send them to the client
	tempUserlist, _ := GetAllCategories(db)
	tempUserlist2 := make([]ServerCategory, len(tempUserlist))
	for i, v := range tempUserlist {
		tempUserlist2[i] = ServerCategory{ID: v.ID, CategoryName: v.CategoryName, Description: v.Description, CreatedAt: v.CreatedAt}
	}

	fmt.Println(Blue + "Server >> Sending Everything to clients" + Reset)

	Broadcast <- ServerMessage{Type: "categories", Categories: tempUserlist2}
}

// function to handle the new category message (when a new category is created)
func handleNewCategoryMessage(message ServerMessage) {
	fmt.Println(Red + "Category Handler called" + Reset)
	// Add the new category to the list of categories
	categories = append(categories, ServerCategory{Name: message.Categories[0].Name, URL: message.Categories[0].URL})

	// Send the updated list of categories to all clients
	Broadcast <- ServerMessage{Type: "categories", Categories: categories}
}

// function to handle the new post message (when a new post is created)
func handleNewPostMessage(message ServerMessage) {
	// Add the new post to the list of posts
	posts = append(posts, ServerPost{Title: message.Posts[0].Title, Content: message.Posts[0].Content, Author: message.Posts[0].Author, Date: message.Posts[0].Date})

	// Send the updated list of posts to all clients
	Broadcast <- ServerMessage{Type: "posts", Posts: posts}
}

// function to handle the get chat history message (when a user wants to see the chat history)
func handleGetChatHistoryMessage(conn *websocket.Conn, message ServerMessage) {
	db := OpenDatabase()
	defer db.Close()
	// Get the conversation history between the two users
	jsonStr, _ := json.Marshal(message.User)
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println(Red, err, Reset)
	}
	username := data["username"].(string)
	// Get userID from username
	userID := GetUserID(db, username)

	chatHistory := GetChatHistory(userID)

	// Send the conversation history to the client
	conn.WriteJSON(ServerMessage{Type: "chat_history", ChatHistory: chatHistory})
}

// function to handle the message message (when a user sends a message to another user)
func handleMessageMessage(conn *websocket.Conn, message ServerMessage) {
	for _, v := range sessions {
		fmt.Println(v.WebSocketConn)
		fmt.Println(v.Username)
	}
	fmt.Println("Message received from ", message.From, " to ", message.To, ": ", message.Text)

	// Add the message to the conversation history
	//open db
	db := OpenDatabase()
	defer db.Close()

	historyTo := GetUserID(db, message.To)
	historyFrom := GetUserID(db, message.From)
	AddMessageToHistory(historyFrom, historyTo, message.Text)

	for _, value := range LoggedInUsers {
		fmt.Println(value.Username, value.WebSocketConn)
		fmt.Println(message.To, message.From)
		if value.Username == message.To {
			message.To = value.WebSocketConn
		}
		if value.Username == message.From {
			message.From = value.WebSocketConn
		}
	}

	// Send the message to the recipient
	for client := range clients {
		fmt.Printf("Client: %s To: %s From: %s, \n", client.RemoteAddr().String(), message.To, message.From)
		if client.RemoteAddr().String() == message.To {
			fmt.Println("Sending message to ", message.To)
			err := client.WriteJSON(message)
			if err != nil {
				log.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleLoginMessage(conn *websocket.Conn, message ServerMessage) {
	//open db
	db := OpenDatabase()
	defer db.Close()
	username := message.Data["username"]
	password := message.Data["password"]

	// Check if the user exist in the database
	if !CheckIfUserExist(db, username) {
		// If the user exist, check if the password is correct
		fmt.Printf(Red+"Server >> User %s tried to login but! but does not exist!!!!!\n"+Reset, username)
		conn.WriteJSON(ServerMessage{Type: "loginResponse", Data: map[string]string{"status": "error", "message": "User Exist!"}})
		conn.WriteJSON(ServerMessage{Type: "loginResponse", Data: map[string]string{"login": "false"}})
		return
	}

	if !CheckIfPasswordIsCorrect(db, username, password) {
		// If the user exist but the password is incorrect
		fmt.Printf(Red+"Server >> User %s tried to login but the password is incorrect!\n"+Reset, username)
		conn.WriteJSON(ServerMessage{Type: "loginResponse", Data: map[string]string{"status": "error", "message": "Incorrect password"}})
		conn.WriteJSON(ServerMessage{Type: "loginResponse", Data: map[string]string{"login": "false"}})
		return
	}

	if UserLoggedIn(username) {
		fmt.Println(Red + "Server >> User already logged in" + Reset)

		// FIX so the current user is logged out and the new user is logged in
		// TODO: FIX THIS

		conn.WriteJSON(ServerMessage{Type: "loginResponse", Data: map[string]string{"status": "error", "message": "User already logged in"}})
		conn.WriteJSON(ServerMessage{Type: "loginResponse", Data: map[string]string{"login": "false"}})
		return
	}

	// Set the session and get the session token
	sessionToken := SetClientCookieWithSessionToken(conn, db, username)

	// Add the user session to the loggedInUsers map
	session := &Session{
		Username:      username,
		Privilege:     GetPrivilegeLevel(db, username),
		Cookie:        sessionToken,
		ExpiredTime:   time.Now().Add(1 * time.Hour), // set expiration time to 1 hour from now WE MAY WANT TO CHANGE THIS!!!!!!
		WebSocketConn: conn.RemoteAddr().String(),
	}
	LoggedInUsers[sessionToken] = session

	conn.WriteJSON(ServerMessage{Type: "loginResponse", Data: map[string]string{"login": "true"}})
	fmt.Printf(Green+"Server >> User %s has logged in!\n"+Reset, username)

}

func handleLoginResponseMessage(conn *websocket.Conn, message ServerMessage) {

	fmt.Println("When am I called?")
}

func handleLogoutMessage(conn *websocket.Conn, message ServerMessage) {
	modifiedCookie := message.Data["cookie"]
	// ugly fix to get the cookie
	if len(message.Data["cookie"]) > 9 {
		modifiedCookie = modifiedCookie[14:]
	}

	//debug print
	fmt.Println("Server >> Cookie: " + modifiedCookie)

	tempUser := ""

	//loop current sesssions
	for _, v := range LoggedInUsers {
		//fmt.Println("v is: " + v.Username + " and v.Cookie is: " + v.Cookie + " and modifiedCookie is: " + modifiedCookie + "")
		if v.Cookie == modifiedCookie {
			tempUser = v.Username
			delete(LoggedInUsers, v.Cookie)

		}
	}

	var templist []ServerUser
	// add all users in LoggedInUsers to the templist
	for _, v := range LoggedInUsers {
		templist = append(templist, ServerUser{Username: v.Username, Privilege: v.Privilege})
	}

	// Get username from session
	var username string
	for _, session := range sessions {
		if strings.Compare(session.Cookie, modifiedCookie) == 0 {
			username = session.Username
			fmt.Println(Yellow + "Server >> User " + username + " has logged out!" + Reset)
		}
	}

	// Remove the user from the LoggedInUsers and sessions map
	delete(LoggedInUsers, username)
	delete(sessions, modifiedCookie)

	//loop current sesssions
	for _, v := range LoggedInUsers {
		fmt.Println("FUCK!!! ", v.Username)
	}

	conn.WriteJSON(ServerMessage{Type: "logoutResponse", Data: map[string]string{"logout": "true"}})
	fmt.Println(Yellow + "Server >> User " + tempUser + " has logged out!" + Reset)
	Broadcast <- ServerMessage{Type: "users", Users: templist}

}

func handleRegisterMessage(conn *websocket.Conn, message ServerMessage) { // we know this works
	//open db
	db := OpenDatabase()
	defer db.Close()

	username := message.Data["username"]
	password := message.Data["password"]
	confirmpassword := message.Data["cfpassword"]
	email := message.Data["email"]
	age := message.Data["age"]
	gender := message.Data["gender"]
	firsnName := message.Data["firstname"]
	lastName := message.Data["lastname"]

	// Debug all the data
	/* 	fmt.Println("Username:", username)
	   	fmt.Println("password:", password)
	   	fmt.Println("confirmpassword:", confirmpassword)
	   	fmt.Println("email:", email)
	   	fmt.Println("age:", age)
	   	fmt.Println("gender", gender)
	   	fmt.Println("firsnName", firsnName)
	   	fmt.Println("lastName", lastName) */

	//convert age to int
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		fmt.Println("Error converting age to int")
	}

	// Check if the user exist in the database
	if CheckIfUserExist(db, username) {
		// If the user exist, redirect to the login page
		fmt.Printf(Red+"Server >> User %s tried to register but already exist!\n"+Reset, username)
		// send to frontend that the user already exist
		if conn != nil {
			conn.WriteJSON(ServerMessage{
				Type: "registerResponse",
				Data: map[string]string{
					"register": "false",
					"status":   "User already exist!",
				},
			})
		}
		return
	}
	if CheckIfEmailExist(db, email) {
		// If the user exist, redirect to the login page
		fmt.Printf(Red+"Server >> User %s tried to register but already exist!\n"+Reset, username)
		// send to frontend that the user already exist
		if conn != nil {
			conn.WriteJSON(ServerMessage{
				Type: "registerResponse",
				Data: map[string]string{
					"register": "false",
					"status":   "Email already exist!",
				},
			})
		}
		return
	}
	if password != confirmpassword {
		// Check if the password and the password confirm are the same
		// If the password and the password confirm are not the same, redirect to the login page
		fmt.Printf(Red+"Server >> User %s tried to register but the password and the password confirm are not the same!\n"+Reset, username)
		if conn != nil {
			conn.WriteJSON(ServerMessage{
				Type: "registerResponse",
				Data: map[string]string{
					"register": "false",
					"status":   "Password mismatched!",
				},
			})
		}
		return
	} else {
		// If the user doesn't exist, Register the user and redirect to the chat page
		fmt.Printf(Green+"Server >> User %s has registered!\n"+Reset, username)

		// Insert the user into the database
		RegisterUser(db, username, ageInt, gender, firsnName, lastName, email, password)

		// Do something so frontend knows that the registration was successful
	}

	conn.WriteJSON(ServerMessage{Type: "registerResponse", Data: map[string]string{"register": "true"}})
}

func handleRegisterResponseMessage(conn *websocket.Conn, message ServerMessage) {
	//open db
	db := OpenDatabase()
	defer db.Close()
	if message.Data["register"] == "true" {
		fmt.Println("Registration successful")
		// signupHandler(conn, message)
	} else {
		fmt.Println("Registration failed")
	}
}

// function handlePostMessage handles the post message from the client
func handleGetPostsMessage(conn *websocket.Conn, message ServerMessage) {
	//open db
	db := OpenDatabase()
	defer db.Close()
	// Get all the posts from the database and send them to the client
	posts, _ := GetLatestPosts(db)
	postList := make([]ServerPost, len(posts))
	for i, v := range posts {
		postList[i] = ServerPost{ID: v.ID, Title: v.Title, Content: v.Content, CreatedAt: v.CreatedAt, UserID: v.UserID, UserName: v.UserName}
	}

	// send the posts to the client
	Broadcast <- ServerMessage{Type: "posts", Posts: postList}
}

func handleGetUsersMessage(conn *websocket.Conn, message ServerMessage) {
	//open db
	db := OpenDatabase()
	defer db.Close()

	var templist []ServerUser
	// add all users in LoggedInUsers to the templist
	for _, v := range LoggedInUsers {
		templist = append(templist, ServerUser{Username: v.Username, Privilege: v.Privilege})
	}
	Broadcast <- ServerMessage{Type: "users", Users: templist}
}

func handleGetOfflineUsersMessage(conn *websocket.Conn, message ServerMessage) {
	//open db
	db := OpenDatabase()
	defer db.Close()

	// Get all the users from the database
	offusers, err := GetAllUsers(db)
	if err != nil {
		fmt.Println("Error getting all users")
	}

	var templist []ServerUser

	// add all users in offusers to the templist and remove the users that are already in LoggedInUsers
	for _, v := range offusers {
		var found bool
		for _, v2 := range LoggedInUsers {
			if v.Username == v2.Username {
				found = true
				break
			}
		}
		if !found {
			templist = append(templist, ServerUser{Username: v.Username, Privilege: v.Privilege})
		}
	}

	fmt.Println("Offline users: ", len(templist))

	Broadcast <- ServerMessage{Type: "offline_users", Users: templist}
}
