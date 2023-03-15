// Create a chat window element and append it to the DOM
const chatWindow = document.createElement("div");
chatWindow.id = "chat-window";
document.body.appendChild(chatWindow);

function openChatWindow(user) {
  // Populate the chat window with the conversation history between the two users (if it exists)
  getChatHistory(user, (chatHistory) => {
    const history = JSON.parse(chatHistory);

    // Create a header container for the chat window
    const headerContainer = document.createElement("div");
    headerContainer.className = "header-container";
    headerContainer.style.position = 'sticky';
    // Set the top position of the header container to 0 so that it sticks to the top of the chat window
    headerContainer.style.top = '0';
    headerContainer.style.zIndex = '1';

    // Create a header for the chat window
    const chatHeader = document.createElement("div");
    chatHeader.className = "chat-header";
    chatHeader.textContent = "Chatting with " + user.username;

    headerContainer.appendChild(chatHeader);

    // Create a close button for the chat window
    const closeButton = document.createElement("button");
    closeButton.textContent = "X";
    closeButton.addEventListener("click", () => {
      chatWindow.innerHTML = "";
      // Remove the chat-window class from the chat window
      chatWindow.classList.remove("chat-window");
    });

    // Append the close button to the chat header
    chatHeader.appendChild(closeButton);

    // Create a container for the chat messages
    const chatContainer = document.createElement("div");
    chatWindow.classList.add("chat-window"); // Add the chat-window class
    chatContainer.className = "chat-container";

    // Array to store received messages
    const receivedMessages = [];

    // Loop through each message in the chat history and display it in the chat window
    history.chathistory.forEach((message) => {
      const messageContainer = document.createElement("div");
      messageContainer.className = message.from === user.id ? "message sent" : "message received";
      
      const messageText = document.createElement("p");
      messageText.textContent = message.text;
      messageContainer.appendChild(messageText);
      
      chatContainer.appendChild(messageContainer);

      // Add received messages to the array
      if (message.from !== user.id) {
        receivedMessages.push(message);
      }
    });
    
    // Create a container for the message input and send button
    const inputContainer = document.createElement("div");
    inputContainer.className = "input-container";
    
    // Create the message input field
    const messageInput = document.createElement("input");
    messageInput.type = "text";
    messageInput.placeholder = "Type your message...";
    inputContainer.appendChild(messageInput);

    // Add an event listener to the message input field to detect when the Enter key is pressed
    messageInput.addEventListener("keydown", (event) => {
      if (event.key === "Enter") {
        // Prevent the default behavior of the Enter key, which is to add a new line
        event.preventDefault();
        // Call the sendMessage function to send the message
        sendMessage(user, messageInput.value);
        //console.log("Sending message: " + messageInput.value);
        // Clear the message input field
        messageInput.value = "";
      }
    });

    // Create the send button
    const sendButton = document.createElement("button");
sendButton.textContent = "Send";
sendButton.addEventListener("click", () => {
sendMessage(user, messageInput.value);
console.log("Sending message: " + messageInput.value);
// Clear the message input field
messageInput.value = "";
});
inputContainer.appendChild(sendButton);

// Add the header container, chat container, and input container to the chat window
chatWindow.innerHTML = "";
chatWindow.appendChild(headerContainer);
chatWindow.appendChild(chatContainer);
chatWindow.appendChild(inputContainer);

document.querySelector('.left-sidebar').style.display = 'block';
document.querySelector('.right-sidebar').style.display = 'block';
document.querySelector('.posts').style.display = 'block';
document.querySelector('.container').style.display = 'flex';
document.querySelector('.container').style.justifyContent = 'space-between';
document.querySelector('.container').style.marginTop = '20px';
document.querySelector('.container').style.flexWrap = 'wrap';
document.getElementById('login-form').style.display = 'none';
document.querySelector('.register-form').style.display = 'none';
document.getElementById('logged-in-message').style.display = 'block';
document.getElementById('logout-form').style.display = 'block';


});
}

// Function to get the conversation history between two users
function getChatHistory(user, callback) {
  // Send a message to the server to request the conversation history
  const message = JSON.stringify({ type: "get_chat_history", user: user });
  socket.send(message);

  // Listen for the response from the server
  socket.addEventListener("message", (event) => {
    const message = JSON.parse(event.data);
    //console.log("Message received from backend!: " + JSON.stringify(message));
    if (message.type === "chat_history") {
      console.log("Chat history received from server: " + JSON.stringify(message));
      const chatHistory = JSON.stringify(message);
      callback(chatHistory);
    }
  });
}

function sendMessage(user, messageText) {
  // Construct the message object
  const message = {
    type: "message",
    from: sessionStorage.getItem("username"),
    to: user.username,
    text: messageText,
  };
  
  // Send the message to the server
  socket.send(JSON.stringify(message));
  //console.log("Message sent to server: " + JSON.stringify(message));

  // Update the chat window with the new message
  const chatContainer = document.querySelector(".chat-container");
  const messageContainer = document.createElement("div");
  messageContainer.className = "message sent";

  if (message.from != user.username) {
  const messageTextElem = document.createElement("p");
  messageTextElem.textContent = messageText;
  messageContainer.appendChild(messageTextElem);
  }

  chatContainer.appendChild(messageContainer);

  // Scroll the chat window to the bottom
  chatContainer.scrollTop = chatContainer.scrollHeight;
  
  // Clear the message input field
  const messageInput = document.querySelector(".input-container input");
  messageInput.value = "";
}

// Listen for received messages from the server
socket.addEventListener("message", (event) => {
  //console.log("Message received from server: " + event.data);
  const message = JSON.parse(event.data);
  if (message.type === "message") {
    //console.log("Message received from server: " + JSON.stringify(message));
    // Add the received message to the chat window
    const chatContainer = document.querySelector(".chat-container");
    const receivedMessageContainer = document.createElement("div");
    receivedMessageContainer.className = "message received";

    const receivedMessageTextElem = document.createElement("p");
    receivedMessageTextElem.textContent = message.text;
    receivedMessageContainer.appendChild(receivedMessageTextElem);

    chatContainer.appendChild(receivedMessageContainer);

    // Scroll the chat window to the bottom
    chatContainer.scrollTop = chatContainer.scrollHeight;
  }
});



