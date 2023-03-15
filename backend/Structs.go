package backend

import (
	"time"
)

// Global maps to store all users, posts, comments, categories, and sessions
var LoggedInUsers = make(map[string]*Session)

// Struct to define a session
type Session struct {
	Username      string
	UserID        int
	Privilege     int
	Cookie        string
	ExpiredTime   time.Time
	WebSocketConn string
}

// Struct to define a user account
type User struct {
	ID        int
	Username  string
	Fname     string
	Lname     string
	Age       int
	Gender    string
	Password  string
	Email     string
	CreatedAt string
	Privilege int
	Send      chan []byte
}

// Struct to define a post
type Category struct {
	ID           int
	CategoryName string
	Description  string
	CreatedAt    string
}

type Comment struct {
	ID             int
	UserID         int
	UserName       string
	PostID         int
	Content        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Date           string
	LikedNumber    int
	DislikedNumber int
	IsEdited       bool
}

type Post struct {
	ID             int
	UserID         int
	UserName       string
	Title          string
	Content        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Date           string
	LikedNumber    int
	DislikedNumber int
	ImgUrl         string
	URL            string
	Approved       int
	Dummy          int
	IsEdited       bool
}

type Reaction struct {
	ID        int
	UserID    int
	PostID    int
	IsLiked   int
	CreatedAt string
}

type Relation struct {
	ID         int
	CategoryID int
	PostID     int
}

// Struct to define a message
type Message struct {
	ID       int
	Username string
	// The username of the recipient of the message (if it is a private message)
	RecipientUsername string
	Content           string
	CreatedAt         string
	From              int       `json:"from"`
	Text              string    `json:"text"`
	ChatHistory       []Message `json:"chathistory"`
	To                int       `json:"to,omitempty"`
	Read              int       `json:"isread"`
}

// Structs to define a hub (for websockets)
type ServerUser struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Privilege int    `json:"privilege"`
	Password  string `json:"password"`
}

type ServerMessage struct {
	Type        string            `json:"type"`
	Users       []ServerUser      `json:"users"`
	Categories  []ServerCategory  `json:"categories"`
	Posts       []ServerPost      `json:"posts"`
	User        ServerUser        `json:"user"`
	Post        ServerPost        `json:"post"`
	Category    ServerCategory    `json:"category"`
	To          string            `json:"to"`
	From        string            `json:"from"`
	Text        string            `json:"text"`
	Username    string            `json:"username"`
	ChatHistory []Message         `json:"chathistory"`
	Password    string            `json:"password"`
	Data        map[string]string `json:"data"`
}

type ServerCategory struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	ID           int    `json:"id"`
	CategoryName string `json:"categoryname"`
	Description  string `json:"description"`
	CreatedAt    string `json:"createdat"`
}
type ServerPost struct {
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	Author         string    `json:"author"`
	Date           string    `json:"date"`
	ID             int       `json:"id"`
	UserID         int       `json:"userid"`
	UserName       string    `json:"username"`
	CreatedAt      time.Time `json:"createdat"`
	UpdatedAt      time.Time `json:"updatedat"`
	LikedNumber    int       `json:"likednumber"`
	DislikedNumber int       `json:"dislikednumber"`
	ImgUrl         string    `json:"imgurl"`
	URL            string    `json:"url"`
	Approved       int       `json:"approved"`
	Dummy          int       `json:"dummy"`
	IsEdited       bool      `json:"isedited"`
}
