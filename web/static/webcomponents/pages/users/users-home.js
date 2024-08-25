import { fromHTML } from "../../utils.js";

const usersHomeTemplate = document.createElement("template");
usersHomeTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/users/users-home.css">
`;

class UsersHome extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.shadow.append(usersHomeTemplate.content.cloneNode(true));
    this.usersData = [];
  }

  renderUsers() {
    const userCount = this.usersData.length;
    const usersList = this.shadow.getElementById("users-list");
    usersList.innerHTML = "";
    this.usersData.forEach((user, index) => {
      const listItem = document.createElement("li");
      // Do not render delete button for the root user
      if (index === userCount - 1) {
        listItem.innerHTML = `
          <span>${userCount - index}. ${user.email}</span>
        `;
      } else {
        listItem.innerHTML = `
          <span>${userCount - index}. ${user.email}</span>
          <button class="delete-button" data-key="${user.id}">Delete</button>
        `;
      }
      usersList.appendChild(listItem);
    });
    this.attachDeleteHandlers();
  }

  attachDeleteHandlers() {
    console.log('Attaching delete handler')
    const deleteButtons = this.shadow.querySelectorAll(".delete-button");
    deleteButtons.forEach((button) => {
      button.addEventListener("click", (event) => {
        const deleteEvent = new CustomEvent("delete-user", {
          detail: event.target.dataset.key,
        });
        // TODO: Change this from the shawdow DOM
        this.shadow.dispatchEvent(deleteEvent);
      });
    });
  }

  async fetchUsers() {
    try {
      const response = await fetch("/api/v1/user");
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const data = await response.json();
      this.usersData = data.data.reverse();;
      this.updateUsersComponent();
    } catch (error) {
      console.error("Error fetching API keys:", error);
    }
  }

  async deleteUser(id) {
    console.log(this.usersData)
    try {
      const response = await fetch(`/api/v1/user/${id}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      });
      if (!response.ok) {
        throw new Error(`Failed to delete User: ${response.statusText}`);
      }
      const numericId = parseInt(id, 10);
      this.usersData = this.usersData.filter((user) => user.id !== numericId);
      console.log(this.usersData)
      this.updateUsersComponent();
    } catch (error) {
      console.error("Error deleting User:", error);
    }
  }
  

  createUser(e) {
    e.preventDefault()
    // Check if a request is being processed by this same client
    if (this.pendingRequest) {
        return
    } else {
        this.pendingRequest = true
    }
    const modal = this.shadow.querySelector('#userModal');
    const data = new FormData(e.target);
    const value = Object.fromEntries(data.entries());
    fetch('/api/v1/user/', {
        method: "POST",
        mode: 'cors',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(value),
    }).then(res => {
        if (res.ok) {
            return res.json();
        }
        throw new Error('Something went wrong');
    }).then(userResponse => {
        const newUser = userResponse.data
        this.usersData.unshift(newUser); // Add new key to the top of the list
        this.updateUsersComponent()
        console.log(`New User with email ${newUser.email} is created.`); // Log creation
        const successMessage = `<h1>User Successfully Created</h1>`
        const modal_element = fromHTML(successMessage);
        modal.innerHTML = '';
        modal.appendChild(modal_element)
    }).catch((error) => {
      const errorMessage = `<h1>Unable to Create User</h1>`
      const modal_element = fromHTML(errorMessage);
      modal.innerHTML = '';
      modal.appendChild(modal_element)
    })
    .finally(() => {
      this.pendingRequest = false; // Set pendingRequest to false when the request is done
    });;
  }

  updateUsersComponent() {
    this.renderUsers()
  }

  connectedCallback() {
    const pageTitleElement = document.createElement("wc-page-heading-button");
    pageTitleElement.innerText = "Users";
    pageTitleElement.buttonText = "New User";
    pageTitleElement.buttonFunction = this.createUserAction.bind(this);
    this.shadow.appendChild(pageTitleElement);

    const userContainerElement = fromHTML(`<ul id="users-list"></ul>`)
    this.shadow.appendChild(userContainerElement);

    this.fetchUsers();
    // TODO: Change this from the shawdow DOM
    this.shadow.addEventListener("delete-user", (event) => {
      this.deleteUser(event.detail)
      }
    );
  }

  createUserAction() {
    let modal = this.shadow.querySelector('#userModal');
    if (modal) {
      modal.remove();
    }
    
    // Create a new modal if it doesn't exist
    const modalMarkup = `
        <wc-modal id="userModal">
            <wc-title>Create New User</wc-title>
            <form id="myForm">
              <label for="email">Email:</label><br>
              <input type="text" id="email" name="email" value=""><br>
              <label for="password">Password:</label><br>
              <input type="text" id="password" name="password" value=""><br>
              <br>
              <input type="submit" value="Create">
            </form>
        </wc-modal>`;
    const modalElement = fromHTML(modalMarkup);
    this.shadow.append(modalElement);

    // Attach the form submit handler
    const formComponent = this.shadow.querySelector('#myForm');
    formComponent.addEventListener('submit', this.createUser.bind(this));


    // Optionally, ensure the modal is visible if hidden
    modal.style.display = 'block';
  }
}

customElements.define("users-home", UsersHome);