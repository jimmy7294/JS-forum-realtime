function addCategoryButton() {
  const categoriesDiv = document.querySelector(".left-sidebar");
  const createCategoryButton = document.createElement("button");
  createCategoryButton.textContent = "Create Category";
  createCategoryButton.className = "create-category-button";
  categoriesDiv.appendChild(createCategoryButton);

  createCategoryButton.addEventListener("click", () => {
    showPopup();
  });
}

function showPopup() {
  const popup = document.createElement("div");
  popup.className = "popup";
  popup.innerHTML = `
<div class="popup-container">
  <div class="popup-content">
    <span class="close">&times;</span>
    <h2>Create Category</h2>
    <form id="create-category-form">
      <label for="category-name">Category Name:</label>
      <input type="text" id="category-name" name="category-name" required>
      <button type="submit">Create</button>
    </form>
  </div>
</div>
  `;

  document.body.appendChild(popup);

  const closeButton = document.querySelector(".close");
  closeButton.addEventListener("click", closePopup);

  const createCategoryForm = document.querySelector("#create-category-form");
  createCategoryForm.addEventListener("submit", (event) => {
    event.preventDefault();

    const categoryName = document.querySelector("#category-name").value;
    // Send the category to the server to be created
    //const message = JSON.stringify({ type: "create_category", name: categoryName });


  // create a new message object
  const message = {
    type: "create_category",
    data: {
      category: categoryName,
      username: document.getElementById('username').value,
    },
  };

  socket.send(JSON.stringify(message));
  console.log("Category created", message);

    closePopup();
  });
}

function closePopup() {
  const popup = document.querySelector(".popup");
  if (popup) {
    document.body.removeChild(popup);
  }
}

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
      //clear the category select before adding new categories
      const categorySelect = document.querySelector("#category");
      categorySelect.innerHTML = "<option value=''>Select a category</option>";
    
      for (const category of message.categories) {
        const categoryLi = document.createElement("li");
        const categoryLink = document.createElement("a");
        categoryLink.textContent = category.categoryname;
        categoryLink.setAttribute("href", category.url);
        categoryLi.appendChild(categoryLink);
        categoriesUl.appendChild(categoryLi);

        const option = document.createElement("option");
        option.value = category.id;
        option.text = category.categoryname;
        categorySelect.appendChild(option);
      }
      // Add the create category if button doesn't already exist
      if (!document.querySelector(".create-category-button"))
        addCategoryButton();

      
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
