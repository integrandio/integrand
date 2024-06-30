import { fromHTML } from "../../utils.js";

const apikeysHomeTemplate = document.createElement("template");
apikeysHomeTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/apiKeys/apiKeys-home.css">
`;

class ApiKeysHome extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.shadow.append(apikeysHomeTemplate.content.cloneNode(true));
    this.apiKeysData = [];
  }

  renderApiKeys() {
    const apiKeyCount = this.apiKeysData.length;
    const apiKeysList = this.shadow.getElementById("api-keys-list");
    apiKeysList.innerHTML = "";
    this.apiKeysData.forEach((key, index) => {
      const listItem = document.createElement("li");
      listItem.innerHTML = `
        <span>${apiKeyCount - index}. ${key.key}</span>
        <button class="delete-button" data-key="${key.key}">Delete</button>
      `;
      apiKeysList.appendChild(listItem);
    });
    this.attachDeleteHandlers();
  }

  attachDeleteHandlers() {
    console.log('Attaching delete handler')
    const deleteButtons = this.shadow.querySelectorAll(".delete-button");
    deleteButtons.forEach((button) => {
      button.addEventListener("click", (event) => {
        const deleteEvent = new CustomEvent("delete-api-key", {
          detail: event.target.dataset.key,
        });
        // TODO: Change this from the shawdow DOM
        this.shadow.dispatchEvent(deleteEvent);
      });
    });
  }

  async fetchApiKeys() {
    try {
      const response = await fetch("/api/v1/apikey");
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const data = await response.json();
      this.apiKeysData = data.data.reverse(); // Reverse the order of the API keys
      this.updateApiKeysComponent();
    } catch (error) {
      console.error("Error fetching API keys:", error);
    }
  }

  async deleteApiKey(apiKey) {
    try {
      const response = await fetch(`/api/v1/apikey/${apiKey}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      });
      if (!response.ok) {
        throw new Error(`Failed to delete API key: ${response.statusText}`);
      }
      this.apiKeysData = this.apiKeysData.filter((key) => key.key !== apiKey);
      this.updateApiKeysComponent();
      console.log(`API key deleted: ${apiKey}`); // Log deletion
    } catch (error) {
      console.error("Error deleting API key:", error);
    }
  }

  async createApiKey() {
    try {
      const response = await fetch("/api/v1/apikey", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });
      if (!response.ok) {
        throw new Error(`Failed to create API key: ${response.statusText}`);
      }
      const data = await response.json();
      const newApiKey = { key: data.data };
      this.apiKeysData.unshift(newApiKey); // Add new key to the top of the list
      this.updateApiKeysComponent();
      console.log(`New API key created: ${newApiKey.key}`); // Log creation
      this.showSuccessMessage(newApiKey.key);
      this.renderApiKeys()
    } catch (error) {
      console.error("Error creating API key:", error);
      this.showErrorMessage();
    }
  }

  updateApiKeysComponent() {
    this.renderApiKeys()
  }

  showSuccessMessage(apiKey) {
    const modalElement = document.createElement("wc-modal");
    modalElement.innerHTML = `
      <div class="modal-content">
          <p>New API Key: <span id="new-api-key">${apiKey}</span></p>
      </div>
    `;
    this.shadow.appendChild(modalElement);
  }

  showErrorMessage() {
    // Display an error message (if needed) for creating a new API key
  }

  connectedCallback() {
    const pageTitleElement = document.createElement("wc-page-heading-button");
    pageTitleElement.innerText = "API Keys";
    pageTitleElement.buttonText = "New API Key";
    pageTitleElement.buttonFunction = this.createApiKey.bind(this);
    this.shadow.appendChild(pageTitleElement);

    const keyContainerElement = fromHTML(`<ul id="api-keys-list"></ul>`)
    this.shadow.appendChild(keyContainerElement);

    this.fetchApiKeys();
    // TODO: Change this from the shawdow DOM
    this.shadow.addEventListener("delete-api-key", (event) => {
      this.deleteApiKey(event.detail)
      }
    );
  }
}

customElements.define("apikeys-home", ApiKeysHome);