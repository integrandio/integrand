const endpointssHomeTemplate = document.createElement("template")
endpointssHomeTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/endpoints/endpoints-home/endpoints-home.css">
<div>
    <wc-page-heading>
        My Profile
    </wc-page-heading>
    <div>
        <button>Reset Password</button>
    </div>
</div>
`

class Profile extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(endpointssHomeTemplate.content.cloneNode(true))
    }
}

customElements.define("wc-profile", Profile)