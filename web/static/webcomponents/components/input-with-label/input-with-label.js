const inputWithLabelTemplate = document.createElement("template")
inputWithLabelTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/input-with-label/input-with-lable.css">
<div class="basic-input-container">
    <label for="email" class="labeler">Email</label>
    <div
      class="input-wrapper"
    >
      <input
        id="email"
        name="email"
        type="email"
        placeholder="you@example.com"
        class="inputer"
      />
    </div>
</div>
`

class InputWithLabel extends HTMLElement {
  constructor(){
    super()
    this.shawdow = this.attachShadow({mode: "open"});
    this.shawdow.append(inputWithLabelTemplate.content.cloneNode(true));
    // this._delete_action = function() {return};
  }

  // set titleText(text) {
  //   this._titleText = text
  // }

  // set descriptionText(text) {
  //   this._descriptionText = text
  // }

  // set buttonFunction(funcValue) {
  //   console.log(funcValue)
  //   this._delete_action = funcValue
  // }

  // buttonClick() {
  //   this._delete_action()
  // }

  connectedCallback() {
    // const heading_text = this.shawdow.querySelector(".heading-content");
    // heading_text.innerText = this._titleText;
    // const heading_description = this.shawdow.querySelector(".details-content");
    // heading_description.innerText = this._descriptionText;

    // const button = this.shawdow.querySelector(".deactivate-button");
    // button.addEventListener('click', this.buttonClick.bind(this))
  }
}

customElements.define("wc-input-label", InputWithLabel)