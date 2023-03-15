// Listen for new messages from the server
socket.addEventListener("message", (event) => {
  //console.log("Received message from server: ", event.data);
  message = JSON.parse(event.data);

  if (message.type === "categories") {
    // console.log("Received categories from server", message.categories);
    // Update the categories list
    const categoriesUl = document.querySelector(".left-sidebar ul");
    // Clear the categories list before adding new categories
    categoriesUl.innerHTML = "<h2>Categories</h2>";

    if (message.categories) {
      const categoriesUl = document.querySelector(".left-sidebar ul");
      // Clear the categories list before adding new categories
      categoriesUl.innerHTML = "<h2>Categories</h2>";
    
      for (const category of message.categories) {
        const categoryLi = document.createElement("li");
        const categoryLink = document.createElement("a");
        categoryLink.textContent = category.categoryname;
        categoryLink.setAttribute("href", category.url);
        categoryLi.appendChild(categoryLink);
        categoriesUl.appendChild(categoryLi);
        
        // document.querySelector('.left-sidebar').style.display = 'block';
        // document.querySelector('.right-sidebar').style.display = 'block';
        // document.querySelector('.posts').style.display = 'block';
        // document.querySelector('.container').style.display = 'flex';
        // document.getElementById('login-form').style.display = 'none';
        // document.querySelector('.register-form').style.display = 'none';
        // document.getElementById('logged-in-message').style.display = 'block';
        // document.getElementById('logout-form').style.display = 'block';
      }
    }
    else {
      const categoriesUl = document.querySelector(".left-sidebar ul");
      // Clear the categories list before adding new categories
      categoriesUl.innerHTML = "<h2>Categories</h2>";
    
      const categoryLi = document.createElement("li");
      categoryLi.textContent = "Please login to see categories";
      categoriesUl.appendChild(categoryLi);
    }
    
  }
});

// Send a message to the server to request the list of categories
message = JSON.stringify({ type: "get_categories" });
socket.send(message);

console.log("Categories.js loaded Getting categories from server...");
