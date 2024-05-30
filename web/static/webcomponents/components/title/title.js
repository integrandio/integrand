const titleTemplate = document.createElement("template")
titleTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/title/title.css">
<div class="heading-container">
    <div class="title-wrapper">
        <h3 class="title">
            <slot></slot>
        </h3>
    </div>
</div>
`

class Title extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(titleTemplate.content.cloneNode(true))
    }
}

customElements.define("wc-title", Title)