const pageHeadingButtonTemplate = document.createElement("template")
pageHeadingButtonTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/page-heading-button/page-heading-button.css">
  <div class="page-heading-button-container">
    <div class="page-heading-button-wrapper">
      <div class="page-heading-button-sub-wrapper">
        <div class="title-container">
          <h2 class="title">
            <slot></slot>
          </h2>
        </div>
        <div class="button-container">
          <button
            type="button"
            class="button"
          >
            Placeholder
          </button>
        </div>
      </div>
    </div>
  </div>
`

class PageHeadingButton extends HTMLElement {
  constructor(){
    super()
    this.shawdow = this.attachShadow({mode: "open"});
    this.shawdow.append(pageHeadingButtonTemplate.content.cloneNode(true));
    this.buttonText = ""
    this.button = this.shawdow.querySelector(".button");
    this._button_action = function() { return };
  }

  set buttonText(text) {
    this.button_text = text
  }

  set buttonFunction(funcValue) {
    this._button_action = funcValue
  }

  buttonClick() {
    this._button_action()
  }

  connectedCallback() {
    this.button.innerText = this.button_text
    this.button.addEventListener('click', this.buttonClick.bind(this))
  }
}

customElements.define("wc-page-heading-button", PageHeadingButton)