    // Listen for new messages from the server
    socket.addEventListener("message", (event) => {
      message = JSON.parse(event.data);
      // console.log("Received message from server: ", message);
  
      if (message.type === "posts") {
        // Update the posts list
        const postsDiv = document.querySelector(".posts");
        // Clear the posts list before adding new posts
        postsDiv.innerHTML = "<h1 style=\"font-family: 'Roboto Mono', monospace; font-weight: bold; font-size: 30px;\">Recent Posts</h1><hr><br>";
  
        if (message.posts) {
          for (const post of message.posts) {
            const postDiv = document.createElement("div");
            postDiv.classList.add("post");
  
            const postTitle = document.createElement("h4");
            postTitle.classList.add("post-title");
            postTitle.textContent = post.title;
  
            const postContent = document.createElement("p");
            postContent.classList.add("post-content");
            postContent.textContent = post.content;
  
            const postMeta = document.createElement("div");
            postMeta.classList.add("post-meta");
  
            const postAuthor = document.createElement("span");
            postAuthor.classList.add("post-author");
            postAuthor.textContent = `Posted by ${post.username}`;
  
            const postDate = document.createElement("span");
            postDate.classList.add("post-date");
            postDate.textContent = post.createdat;
  
            const commentButton = document.createElement("button");
            commentButton.classList.add("comment-btn");
            commentButton.textContent = "Comment";
  
            postMeta.appendChild(postAuthor);
            postMeta.appendChild(postDate);
            postMeta.appendChild(commentButton);
            

            postDiv.appendChild(postTitle);
            postDiv.appendChild(postContent);
            postDiv.appendChild(postMeta);
  
            postsDiv.appendChild(postDiv);
            // console.log("Received posts from server", message.posts);
            

            document.querySelector('.left-sidebar').style.display = 'block';
            document.querySelector('.right-sidebar').style.display = 'block';
            document.querySelector('.posts').style.display = 'block';
            document.querySelector('.container').style.display = 'flex';
            document.getElementById('login-form').style.display = 'none';
            document.querySelector('.register-form').style.display = 'none';
            document.getElementById('logged-in-message').style.display = 'block';
            document.getElementById('logout-form').style.display = 'block';
            
          }
        }
        else {
          const noPostsDiv = document.createElement("div");
          noPostsDiv.textContent = "Please login to see posts";
          postsDiv.appendChild(noPostsDiv);
          // hide the new-post-form 
            document.querySelector(".new-post-form").style.display = "none";
        }
      }
    });
  
    // Send a message to the server to request the list of posts
    message = JSON.stringify({ type: "get_posts" });
    socket.send(message);

    console.log("Posts.js loaded Getting posts from server...");