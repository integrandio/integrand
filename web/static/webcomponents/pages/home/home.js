

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
                    <wc-title>Jobs</wc-title>
                </div>
                <div>
                    <wc-title>Datasources</wc-title>
                </div>
                <div>
                    <wc-title>Report</wc-title>
                </div>
            </div>
        </div>
        `
    }
}

customElements.define("wc-home", Home)