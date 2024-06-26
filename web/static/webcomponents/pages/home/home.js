
class Home extends HTMLElement {
    constructor(){
        super()
        this.innerHTML = `
        <div>
            <wc-page-heading>
                Home
            </wc-page-heading>
            <div>
                <div>
                    <wc-title>Connectors</wc-title>
                </div>
                <div>
                    <wc-title>Topics</wc-title>
                </div>
                <div>
                    <wc-title>API Keys</wc-title>
                </div>
            </div>
        </div>
        `
    }
}

customElements.define("wc-home", Home)