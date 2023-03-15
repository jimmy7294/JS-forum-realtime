    // Listen for new messages from the server
    socket.addEventListener("message", (event) => {
        //console.log("Message received from server: " + event.data)
        message = JSON.parse(event.data);
    
        if (message.type === "offline_users") {
          // Debugging to see the list of users
/*           for (value of message.users) {
            console.log("offline User >>>> ",value);
          } */

          // Update the user list
          const usersDiv = document.getElementById("offline-list");
          // Clear the user list before adding new users
          usersDiv.innerHTML = "<h2>Logged In Users</h2>";
        
          if (message.type === "offline_users") {
            // Update the user list
            const usersDiv = document.getElementById("offline-list");
            // Clear the user list before adding new users
            usersDiv.innerHTML = "<h2>Offline Users</h2>";
          
            if (message.users) {
              let count = 0;
              for (const user of message.users) {
                if (count >= 10) {
                  break; // Limit the list to a maximum of 10 users
                }
                const userElem = document.createElement("div");
                //userElem.textContent = user.username;
                
                // Create a green circle element to indicate that the user is logged in
                const greenCircle = document.createElement("span");
                greenCircle.style.backgroundColor = "red";
                greenCircle.style.width = "10px";
                greenCircle.style.height = "10px";
                greenCircle.style.borderRadius = "50%";
                greenCircle.style.display = "inline-block";
                greenCircle.style.marginRight = "10px";
                greenCircle.style.marginLeft = "10px";
                
                // Add the green circle element to the user element
                userElem.appendChild(greenCircle);
                userElem.appendChild(document.createTextNode(user.username));

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
                
                // Add a click event listener to each user element that opens a chat window when clicked
                userElem.addEventListener("click", () => {
                  openChatWindow(user);
                  console.log("Clicked on user: " + user.username);
                });
          
                usersDiv.appendChild(userElem);
                count++;
              }
            } 
          }       
        }
    
      });
    
      // Send a message to the server to request the list of users
      message = JSON.stringify({ type: "get_offline_users" });
      socket.send(message);
    
      console.log("Display-offline.js loaded Getting users from server...");