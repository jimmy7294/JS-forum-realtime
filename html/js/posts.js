function requestComments(postId) {
  //save post title to global variable as a string
  console.log("Requesting comments for post", postId);
  // save postID to session storage
  sessionStorage.setItem("postID", postId);
  // create a new message object
  const message = {
    type: "get_comments",
    data: {
      content: postId,
      username: document.getElementById('username').value,
    },
  };

 
  // send the message object as a JSON string to the server
  console.log("Sending message to server", message);
  socket.send(JSON.stringify(message));
}

function openCommentsModal(comments) {
  console.log("Opening comments modal for post", comments);
  // Create a modal to display the comments and the input field for new comments
  const modal = document.createElement("div");
  modal.classList.add("comments-modal");
  modal.style.display = 'block';

  // Add a close button to the modal
  const closeButton = document.createElement("button");
  closeButton.textContent = "X";
  closeButton.classList.add("close-btn");
  closeButton.onclick = () => {
    modal.remove();
  };
  modal.appendChild(closeButton);

  // Display the comments
  const commentsDiv = document.createElement("div");
  commentsDiv.classList.add("comments");
  modal.appendChild(commentsDiv);

  if (comments) {
    for (const comment of comments) {
      const commentDiv = document.createElement("div");
      commentDiv.classList.add("comment");

      const commentContent = document.createElement("p");
      commentContent.classList.add("comment-content");
      commentContent.textContent = comment.content;

      const commentAuthor = document.createElement("span");
      commentAuthor.classList.add("comment-author");
      commentAuthor.textContent = `Comment by ${comment.username}`;

      commentDiv.appendChild(commentContent);
      commentDiv.appendChild(commentAuthor);

      commentsDiv.appendChild(commentDiv);
    }
  } else {
    const noCommentsDiv = document.createElement("div");
    noCommentsDiv.classList.add("no-comments");
    noCommentsDiv.textContent = "No comments yet.";
    commentsDiv.appendChild(noCommentsDiv);
  }

  // Add an input field for new comments
  const newCommentForm = document.createElement("form");
  newCommentForm.classList.add("new-comment-form");
  const newCommentInput = document.createElement("input");
  newCommentInput.classList.add("new-comment-input");
  newCommentInput.type = "text";
  newCommentInput.placeholder = "Write a comment...";
  newCommentForm.appendChild(newCommentInput);

  newCommentForm.onsubmit = (event) => {
    event.preventDefault();
    // Send the new comment to the server
    // create a new message object
    const message = {
      type: "new_comment",
      data: {
        content: newCommentInput.value,
        username: sessionStorage.getItem("username"),
        postid: sessionStorage.getItem("postID"),
      },
    };

    // send the message object as a JSON string to the server
    socket.send(JSON.stringify(message));
    console.log("Sending new comment to server", message);

    // update the comments list
    const commentDiv = document.createElement("div");
    commentDiv.classList.add("comment");

    const commentContent = document.createElement("p");
    commentContent.classList.add("comment-content");
    commentContent.textContent = newCommentInput.value;

    const commentAuthor = document.createElement("span");
    commentAuthor.classList.add("comment-author");
    commentAuthor.textContent = `Comment by ${document.getElementById('username').value}`;

    commentDiv.appendChild(commentContent);
    commentDiv.appendChild(commentAuthor);

    commentsDiv.appendChild(commentDiv);

    // Reset the input field
    newCommentInput.value = "";
  };

  modal.appendChild(newCommentForm);

  // Add the modal to the body
  document.body.appendChild(modal);
}

    // Listen for new messages from the server
    socket.addEventListener("message", (event) => {
      message = JSON.parse(event.data);
      // console.log("Received message from server: ", message);

      if (message.type === "comments") {
        console.log("Received comments from server", message.comment);
        const postId = message.comment.content;
        const comments = message.comment;
        openCommentsModal(comments);
      }    
  
      if (message.type === "posts") {
        // Update the posts list
        const postsDiv = document.querySelector(".posts");
        console.log("Received posts from server", message.posts);
        // Clear the posts list before adding new posts
        postsDiv.innerHTML = "<h1 style=\"font-family: 'Roboto Mono', monospace; font-weight: bold; font-size: 30px;\">Recent Posts</h1>";

        if (message.posts) {
          for (const post of message.posts) {
            const postDiv = document.createElement("div");
            postDiv.classList.add("post");

            const hr2 = document.createElement("hr");
            postDiv.appendChild(hr2);
  
            const postTitle = document.createElement("h4");
            postTitle.classList.add("post-title");
            postTitle.textContent = post.title;
  
            const postContent = document.createElement("p");
            postContent.classList.add("post-content");
            postContent.textContent = post.content;
  
            const postMeta = document.createElement("div");
            postMeta.classList.add("post-meta");
            postMeta.textContent = `Category: ${post.category}`;
  
            const postAuthor = document.createElement("span");
            postAuthor.classList.add("post-author");
            postAuthor.textContent = `Posted by: ${post.username}`;

            //convert timestamp to date
            const date = new Date(post.createdat);
  
            const postDate = document.createElement("span");
            postDate.classList.add("post-date");
            postDate.textContent = `Post created: ${date.toDateString()}`;

            const commentButton = document.createElement("button");
            commentButton.classList.add("comment-btn");
            commentButton.textContent = "Comment";
            commentButton.onclick = () => requestComments(post.title);
  
            postMeta.appendChild(postAuthor);
            postMeta.appendChild(postDate);
            postMeta.appendChild(commentButton);
            

            postDiv.appendChild(postTitle);
            postDiv.appendChild(postContent);
            postDiv.appendChild(postMeta);
            
            //add hr after each post
            const br = document.createElement("br");
            postDiv.appendChild(br);
            const hr = document.createElement("hr");
            postDiv.appendChild(hr);
  
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
    console.log("Sending get_posts to server", message)

    console.log("Posts.js loaded Getting posts from server...");