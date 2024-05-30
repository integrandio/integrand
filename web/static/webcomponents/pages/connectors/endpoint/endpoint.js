import { fromHTML } from '../../../utils.js'

const jobTemplate = document.createElement("template")
jobTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/connectors/endpoint/endpoint.css">
`

class EndpointPage extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.endpoint_id = this.getAttribute("endpoint_id")
        this.shawdow.append(jobTemplate.content.cloneNode(true))
    }

    generateMarkup(endpoint) {
        const endpoint_link = `/api/v1/glue/f/${endpoint.id}`
        var date = new Date(endpoint.lastModified);
        let job_markup = `
        <div>
            <wc-page-heading>Endpoint Details: ${endpoint.id}</wc-page-heading>
            <ul class="endpointContainerCard">
                <li>
                    <p class="titler">Endpoint URL:<p>
                    <a href=${endpoint_link}>${endpoint_link}</a>
                </li>
                <li>
                    <p class="titler">Connection Key:</p>
                    <p>${endpoint.connectionKey}</p>
                </li>
                <li>
                    <p class="titler">Topic Name:</p>
                    <p>${endpoint.topicName}</p>
                </li>
                <li>
                    <p class="titler">Last Modified:</p>
                    <p>${date.toDateString()}</p>
                </li>
            </div>
        </div>
        `
        const div = fromHTML(job_markup);
        return div;
    }

    async connectedCallback(){
        const response = await fetch(`/api/v1/glue/${this.endpoint_id}`);
        const jsonData = await response.json();
        let element = this.generateMarkup(jsonData.data)
        const contentTemplate = document.createElement("template")
        contentTemplate.content.append(element)
        this.shawdow.append(contentTemplate.content.cloneNode(true))
    }
}

customElements.define("wc-endpoint", EndpointPage)