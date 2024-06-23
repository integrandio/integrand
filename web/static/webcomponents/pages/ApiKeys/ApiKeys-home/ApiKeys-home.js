import { fromHTML } from "../../../utils.js";

const apikeysHomeTemplate = document.createElement("template");
apikeysHomeTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/apikeys/apikeys-home/apikeys-home.css">
`;

class ApiKeysHome extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.shadow.append(apikeysHomeTemplate.content.cloneNode(true));
    this.apiKeysData = [];
  }

  async fetchApiKeys() {
    try {
      const response = await fetch("/api/v1/apikeys");
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const data = await response.json();
      this.apiKeysData = data.data.apiKeys.reverse(); // Reverse the order of the API keys
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
      const newApiKey = { key: data.data.apiKey };
      this.apiKeysData.unshift(newApiKey); // Add new key to the top of the list
      this.updateApiKeysComponent();
      console.log(`New API key created: ${newApiKey.key}`); // Log creation
      this.showSuccessMessage(newApiKey.key);
    } catch (error) {
      console.error("Error creating API key:", error);
      this.showErrorMessage();
    }
  }

  updateApiKeysComponent() {
    const apiKeysComponent = this.shadow.querySelector("api-keys");
    if (apiKeysComponent) {
      apiKeysComponent.apiKeys = this.apiKeysData;
    }
  }

  showSuccessMessage(apiKey) {
    const apiKeysComponent = this.shadow.querySelector("api-keys");
    if (apiKeysComponent) {
      apiKeysComponent.displayNewApiKey(apiKey);
    }
  }

  showErrorMessage() {
    // Display an error message (if needed) for creating a new API key
  }

  connectedCallback() {
    const pageTitleElement = document.createElement("wc-page-heading-button");
    pageTitleElement.innerText = "All API Keys";
    pageTitleElement.buttonText = "New API Key";
    pageTitleElement.buttonFunction = this.createApiKey.bind(this);
    this.shadow.appendChild(pageTitleElement);

    const apiKeysComponent = document.createElement("api-keys");
    apiKeysComponent.addEventListener("create-api-key", () =>
      this.createApiKey(),
    );
    apiKeysComponent.addEventListener("delete-api-key", (event) =>
      this.deleteApiKey(event.detail),
    );
    this.shadow.appendChild(apiKeysComponent);

    this.fetchApiKeys();
  }
}

customElements.define("apikeys-home", ApiKeysHome);
