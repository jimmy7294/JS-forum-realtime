package backend

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var time_format = "2006-01-02 15:04:05"

// Opening the database
func OpenDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "database/forum.db?parseTime=true")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	return db
}

// Closing the database
func CloseDatabase(db *sql.DB) {
	db.Close()
}

// Function to check if user exist in database, return user privilege
func GetUserPrivilege(db *sql.DB, username string) int {
	var expected_user User

	row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where username = ?", username)
	err := row.Scan(&expected_user.ID, &expected_user.Privilege, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	if err != nil {
		return 0
	}

	return expected_user.Privilege
}

// Function checks if the email exist in the database, return true if exist, false if not
func CheckIfEmailExist(db *sql.DB, email string) bool {
	fmt.Println("Checking if email exist")
	var expected_user User

	row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where email = ?", email)
	err := row.Scan(&expected_user.ID, &expected_user.Privilege, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	return err == nil
}

// Function to check if user exist in database, return true if exist, false if not
func CheckIfUserExist(db *sql.DB, username string) bool {
	var expected_user User

	row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where username = ?", username)
	err := row.Scan(&expected_user.ID, &expected_user.Privilege, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	return err == nil
}

// Function to register a user in the database
func RegisterUser(db *sql.DB, username string, age int, gender string, fname string, lname string, email string, password string) string {
	// Set privilege to 1 CHANGE LATER!!!
	prev := 1

	// Insert the user into the database
	_, err := db.Exec("insert into user(privilege,username,passwrd,email,fname,lname,age,gender,created_at) values(?,?,?,?,?,?,?,?,?)", prev, username, password, email, fname, lname, age, gender, time.Now().Local().Format(time_format))
	if err != nil {
		log.Fatal(err)
	}
	return "200 OK"
}

// Function to check if the password is correct
func CheckIfPasswordIsCorrect(db *sql.DB, username string, password string) bool {
	var expected_user User

	row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where username = ?", username)
	err := row.Scan(&expected_user.ID, &expected_user.Privilege, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	if err != nil {
		return false
	}

	if expected_user.Password == password {
		return true
	}

	return false
}

func GetPrivilegeLevel(db *sql.DB, username string) int {
	var expected_user User

	row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where username = ?", username)
	err := row.Scan(&expected_user.ID, &expected_user.Privilege, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	if err != nil {
		return 0
	}

	return expected_user.Privilege
}

// Getting all the categories from the database
func GetAllCategories(db *sql.DB) ([]Category, string) {

	// Creating a variable to hold the categories
	var category Category
	var categories []Category

	// Selecting all fields from the category table
	rows, err := db.Query("select id,category_name,descript,created_at from category")
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return nil, "500 INTERNAL SERVER ERROR: CATEGORY DATA CORRUPTED"
	}
	defer rows.Close()

	// Looping through each row and saving the returned category
	for rows.Next() {
		err := rows.Scan(&category.ID, &category.CategoryName, &category.Description, &category.CreatedAt)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: CATEGORY DATA CORRUPTED"
		}

		// Appending the returned category to the categories list
		categories = append(categories, category)
	}

	return categories, "200 OK"
}

func GetLatestPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at, p.updated_at, p.liked_no, p.disliked_no, p.img_url, p.approved, p.dummy, CASE WHEN p.updated_at > p.created_at THEN 1 ELSE 0 END AS is_edited FROM post p INNER JOIN user u ON p.user_id = u.id ORDER BY p.created_at DESC LIMIT 5")
	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var post Post
		var createdAt, updatedAt time.Time
		err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.Title, &post.Content, &createdAt, &updatedAt, &post.LikedNumber, &post.DislikedNumber, &post.ImgUrl, &post.Approved, &post.Dummy, &post.IsEdited)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt.UTC()
		post.UpdatedAt = updatedAt.UTC()
		post.Date = post.CreatedAt.Format("January 2, 2006 at 3:04pm")
		post.URL = "/posts/" + strconv.Itoa(post.ID)
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// update posts with correct username
	for i := 0; i < len(posts); i++ {
		temp, _ := GetUser(db, posts[i].UserID)
		temp2 := string(temp.Username)
		posts[i].UserName = temp2
	}

	return posts, nil
}

// Get user depending on userID
func GetUser(db *sql.DB, userID int) (User, error) {
	var user User
	row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where id = ?", userID)
	err := row.Scan(&user.ID, &user.Privilege, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserID(db *sql.DB, username string) string {
	// Prepare the SQL query to retrieve the user ID based on the username
	query := "SELECT id FROM user WHERE username = ?"

	// Execute the query and retrieve the user ID
	var userID string
	err := db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		fmt.Printf(Red+"Server >> Error getting user ID: %s"+Reset, err)
	}

	return userID
}

// Get chat history for a user from database
func GetChatHistory(user string) []Message {
	fmt.Println(Green + "Server >> Getting chat history for user: " + user + Reset)
	// Connect to database
	db := OpenDatabase()
	defer db.Close()
	// Get messages from database/persistent storage
	rows, err := db.Query("SELECT from_user, to_user, is_read, message FROM message WHERE (from_user = ? OR to_user = ?)", user, user)

	if err != nil {
		fmt.Printf(Red+"Server >> Error getting chat history: %s"+Reset, err)
	}

	messages := []Message{}
	for rows.Next() {
		// Read row data
		var fromUser, toUser, isread int
		var message string
		err = rows.Scan(&fromUser, &toUser, &isread, &message)
		if err != nil {
			fmt.Printf(Red+"Server >> Error reading chat history: %s"+Reset, err)
		}

		// Create Message object
		msg := Message{
			From: fromUser,
			To:   toUser,
			Read: isread,
			Text: message,
		}
		messages = append(messages, msg)
	}

	return messages
}
func AddMessageToHistory(fromUser string, toUser string, messageText string) {
	isread := 0
	//DEBUG PRINT
	// fmt.Println("Adding message to database")
	// fmt.Printf("From: %s To: %s Message: %s IsRead: %d ", fromUser, toUser, messageText, isread)

	// Connect to database
	db := OpenDatabase()
	defer db.Close()
	// Insert message into database
	_, err := db.Exec("INSERT INTO message (from_user, to_user, is_read, message,created_at) VALUES (?, ?, ?, ?,?)", fromUser, toUser, isread, messageText, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		// Handle error
		fmt.Printf(Red+"Server >> Error adding message to database: %s "+Reset, err)
	}
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id,privilege,username,passwrd,email,created_at FROM user")
	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var user User
		var createdAt time.Time
		err := rows.Scan(&user.ID, &user.Privilege, &user.Username, &user.Password, &user.Email, &createdAt)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = time.Now().Format(time_format)
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
