const reportsHomeTemplate = document.createElement("template")
reportsHomeTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<div>
    <wc-page-heading>
        Reports
    </wc-page-heading>
</div>
` 

class ReportsHome extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(reportsHomeTemplate.content.cloneNode(true))
    }
}

customElements.define("reports-home", ReportsHome)