const notFoundTemplate = document.createElement("template")
notFoundTemplate.innerHTML = `
<style>
</style>
<div>
<h1>Page not Found</h1>
</div>

`

class NotFound extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(notFoundTemplate.content.cloneNode(true))
    }
}

customElements.define("wc-notfound", NotFound)