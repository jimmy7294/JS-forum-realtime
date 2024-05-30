# Real-Time Forum

This project is an upgraded version of a previous forum implementation. It includes features such as private messages, real-time actions, and live sharing video. The project is divided into five parts:

- SQLite: Used to store data, similar to the previous forum.
- Golang: Handles data and Websockets (Backend).
- JavaScript: Handles frontend events and client Websockets.
- HTML: Organizes the elements of the page.
- CSS: Stylizes the elements of the page.

## Objectives

The project focuses on the following points:

### Registration and Login

- Users must register and login to access the forum.
- Registration form should include fields like nickname, age, gender, first name, last name, email, and password.
- Users can connect using either their nickname or email combined with the password.
- Users can log out from any page on the forum.

### Posts and Comments

- Users can create posts with categories.
- Users can create comments on posts.
- Posts are displayed in a feed.
- Comments are visible when a user clicks on a post.

### Private Messages

- Users can send private messages to each other.
- The chat section shows who are online/offline and available to talk to.
- The chat section is organized by the last message sent.
- Users can send private messages to online users.
- Users can view previous messages with other users.
- Messages have a specific format including the date and username.
- Messages work in real-time using WebSockets.

## Installation

To install the project, follow these steps:

1. Install the required packages:
  - [github.com/gofrs/uuid v4.4.0+incompatible](https://github.com/gofrs/uuid)
  - [github.com/gorilla/websocket v1.5.0](https://github.com/gorilla/websocket)
  - [github.com/mattn/go-sqlite3 v1.14.16](https://github.com/mattn/go-sqlite3)

2. Run the project:
  - Open the root folder of the project.
  - Run `go run .` to start the server.

## Usage

To use the project, follow these steps:

1. Open your browser and navigate to `localhost`.
2. Go to `http://localhost:8080` and start exploring.


## License

This project is licensed under the [CC BC-NC] Creative Commons Attribution-NonCommercial 2.0 Generic License.

## Screenshots
[cmd](./img/cmd.png)
[login](./img/login.png)
[chat](./img/chat.png)

