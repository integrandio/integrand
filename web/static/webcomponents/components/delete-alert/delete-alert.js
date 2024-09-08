const deleteAlertTemplate = document.createElement("template")
deleteAlertTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/delete-alert/delete-alert.css">
<div class="delete-alert-container">
  <div class="text-icon-container">
    <div class="icon-wrapper">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        aria-hidden="true"
        class="svger"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
        ></path>
      </svg>
    </div>
    <div class="text-wrapper">
      <h3 class="heading-content">
      </h3>
      <div class="details-content-wrapper">
        <p class="details-content">
        </p>
      </div>
    </div>
  </div>
  <div class="buttons-container">
    <button type="button" class="deactivate-button">
      Delete</button
    ><button
      type="button"
      data-autofocus="true"
      class="cancel-button"
    >
      Cancel
    </button>
  </div>
</div>
`

class DeleteAlert extends HTMLElement {
  constructor(){
    super()
    this.shawdow = this.attachShadow({mode: "open"});
    this.shawdow.append(deleteAlertTemplate.content.cloneNode(true));
    this._delete_action = function() {return};
  }

  set titleText(text) {
    this._titleText = text
  }

  set descriptionText(text) {
    this._descriptionText = text
  }

  set buttonFunction(funcValue) {
    console.log(funcValue)
    this._delete_action = funcValue
  }

  buttonClick() {
    this._delete_action()
  }

  connectedCallback() {
    const heading_text = this.shawdow.querySelector(".heading-content");
    heading_text.innerText = this._titleText;
    const heading_description = this.shawdow.querySelector(".details-content");
    heading_description.innerText = this._descriptionText;

    const button = this.shawdow.querySelector(".deactivate-button");
    button.addEventListener('click', this.buttonClick.bind(this))
  }
}

customElements.define("wc-delete-alert", DeleteAlert)