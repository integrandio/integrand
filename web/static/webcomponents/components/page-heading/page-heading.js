const pageHeadingTemplate = document.createElement("template")
pageHeadingTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/page-heading/page-heading.css">
<div class="page-heading-container">
    <div class="page-heading-wrapper">
      <div>
        <div class="heading-content-container">
          <div class="heading-container">
            <h2 class="heading-title">
                <slot></slot>
            </h2>
          </div>
        </div>
      </div>
    </div>
  </div>
`

class PageHeading extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(pageHeadingTemplate.content.cloneNode(true))
    }
}

customElements.define("wc-page-heading", PageHeading)