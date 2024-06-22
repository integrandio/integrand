import { fromHTML } from "../../utils.js";

const apiKeysTemplate = document.createElement("template");
apiKeysTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/ApiKeys/ApiKeys.css">
<div class="api-keys-container">
    <ul id="api-keys-list"></ul>
</div>
`;

class ApiKeys extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.shadow.appendChild(apiKeysTemplate.content.cloneNode(true));
    this.apiKeysData = [];
    this.apiKeyCount = 0;
  }

  set apiKeys(data) {
    this.apiKeysData = data;
    this.apiKeyCount = data.length;
    this.renderApiKeys();
  }

  renderApiKeys() {
    const apiKeysList = this.shadow.getElementById("api-keys-list");
    apiKeysList.innerHTML = "";
    this.apiKeysData.forEach((key, index) => {
      const listItem = document.createElement("li");
      listItem.innerHTML = `
        <span>${this.apiKeyCount - index}. ${key.key}</span>
        <button class="delete-button" data-key="${key.key}">Delete</button>
      `;
      apiKeysList.appendChild(listItem);
    });
    this.attachDeleteHandlers();
  }

  attachDeleteHandlers() {
    const deleteButtons = this.shadow.querySelectorAll(".delete-button");
    deleteButtons.forEach((button) => {
      button.addEventListener("click", (event) => {
        const deleteEvent = new CustomEvent("delete-api-key", {
          detail: event.target.dataset.key,
        });
        this.dispatchEvent(deleteEvent);
      });
    });
  }

  displayNewApiKey(apiKey) {
    const modalElement = document.createElement("wc-modal");
    modalElement.innerHTML = `
          <div class="modal-content">
              <p>Newest API Key: <span id="new-api-key">${apiKey}</span></p>
          </div>
      `;
    this.shadow.appendChild(modalElement);
  }

  connectedCallback() {
    const createButton = this.shadow.getElementById("create-api-key-button");
    createButton.addEventListener("click", () => {
      const createEvent = new CustomEvent("create-api-key");
      this.dispatchEvent(createEvent);
    });
  }
}

customElements.define("api-keys", ApiKeys);
