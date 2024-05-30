import { randomColorPicker } from '../../utils.js'

const dataCardCssTemplate = document.createElement("template")
dataCardCssTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/data-element-card/data-element-card.css">
`

class DataElementCard extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(dataCardCssTemplate.content.cloneNode(true))
        this.src_link = this.getAttribute("src-link")
        this.card_title = this.getAttribute("card-title")
        this.card_sub_title = this.getAttribute("card-sub-title")
        const ext_atr = this.getAttribute("external")
        this.external = ext_atr == null ? false : true
        this.image_link = this.getAttribute("image-link");
    }

    generateDataCardHtml() {
        const target_string = this.external ? "target=\"_blank\"" : "";
        let image_inner;
        if (this.image_link != null) {
            image_inner = `<img class="imger" src="${this.image_link}"/>`
        } else {
            const titleNameSplit = this.card_title.split(" ");
            const titleInitials = titleNameSplit.map((word) => word[0]).join('')
            const color = randomColorPicker()
            image_inner = `<div class="card-boxer" style="background-color: rgb(${color} / var(--tw-bg-opacity));">${titleInitials}</div>` 
        }
        
        let image = `<img class="imger" src="${this.image_link}"/>`
        let dataCardHtml = `
        <li class="database-card">
            <div class="image-wrapper">
                ${image_inner}
            </div>
            <div class="card-title-container">
                <a href="${this.src_link}" ${target_string} class="linker">
                    <span class="spaner" aria-hidden="true"></span>
                    <p class="title">${this.card_title}</p>
                    <p class="subtitle">${this.card_sub_title}</p>
                <a>
            </div>
        </li>`
        return dataCardHtml;
    }

    connectedCallback(){
        const contentTemplate = document.createElement("template")
        const dataCardMarkup = this.generateDataCardHtml()
        contentTemplate.innerHTML = dataCardMarkup;
        this.shawdow.append(contentTemplate.content.cloneNode(true))
    }
}

customElements.define("data-element-card", DataElementCard)