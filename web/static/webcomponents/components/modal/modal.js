const modalTemplate = document.createElement("template")
modalTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/modal/modal.css">

<div id="modal">
    <div class="modal-underlay"></div>
    <div id="modalContent" class="modal-content">
        <slot></slot>
    </div>
</div>
`

class Modal extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(modalTemplate.content.cloneNode(true))
        this.modal = this.shawdow.querySelector("#modal")
        this.contentContainer = this.shawdow.querySelector("#modalContent")
    }

    stopClose(event) {
        event.stopPropagation()
    }

    removeButton() {
        const modal = this.shawdow.querySelector("#modal")
        modal.remove()
    }

    connectedCallback(){
        this.modal.addEventListener('click', this.removeButton.bind(this))
        this.contentContainer.addEventListener('click', this.stopClose)
        console.log("loaded in the modal")
    }
}

customElements.define("wc-modal", Modal)