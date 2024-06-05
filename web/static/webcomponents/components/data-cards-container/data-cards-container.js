const dccTemplate = document.createElement("template")
dccTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/data-cards-container/data-cards-container.css">
`

class DataCardsContainer extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(dccTemplate.content.cloneNode(true))
    }

    getHtmlElement() {
        return `
        <div class="datasource-container">
            <ul role="list" class="database-cards-list">
                <slot></slot>
            </ul>
        </div>`
    }

    connectedCallback(){
        const contentTemplate = document.createElement("template")
        const markup = this.getHtmlElement()
        contentTemplate.innerHTML = markup;
        this.shawdow.append(contentTemplate.content.cloneNode(true))
    }
}

customElements.define("data-cards-container", DataCardsContainer)