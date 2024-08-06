import { fromHTML } from '../../../utils.js'

const jobTemplate = document.createElement("template")
jobTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/endpoints/endpoint/endpoint.css">
`

class EndpointPage extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.endpoint_id = this.getAttribute("endpoint_id")
        this.shawdow.append(jobTemplate.content.cloneNode(true))
    }

    deleteConnectorAction() {
        const url = `/api/v1/connector/${this.endpoint_id}`
        fetch(url, {
            method: "Delete",
            headers: {"Content-Type": "application/json",}
        }).then(res => {
            // Check response to see if it's bad
            res.json().then((endpointResponseData) => {
                console.log(endpointResponseData)
                window.location.replace("/endpoints");
            });
        })
    }

    generateMarkup(endpoint) {
        const endpoint_link = `/api/v1/connector/f/${endpoint.id}?apikey=${endpoint.securityKey}`;
        const topic_link = `/topics/${endpoint.topicName}`
        var date = new Date(endpoint.lastModified);
        let job_markup = `
            <ul class="endpointContainerCard">
                <li>
                    <p class="titler">Connection Key:</p>
                    <p>${endpoint.securityKey}</p>
                </li>
                <li>
                    <p class="titler">Topic Name:</p>
                    <a class="link" href=${topic_link}>${endpoint.topicName}</a>
                </li>
                <li>
                    <p class="titler">Last Modified:</p>
                    <p>${date.toDateString()}</p>
                </li>
                <li>
                    <p class="titler">Endpoint URL:<p>
                    <a class="link" href=${endpoint_link}>${endpoint_link}</a>
                </li>
            </ul>`
        const div = fromHTML(job_markup);
        return div;
    }

    async connectedCallback(){
        const response = await fetch(`/api/v1/connector/${this.endpoint_id}`);
        const jsonData = await response.json();
        let element = this.generateMarkup(jsonData.data)
        const contentTemplate = document.createElement("template")
        contentTemplate.content.append(element)

        // Create the title and the delete button
        const pageTitleElement = document.createElement("wc-page-heading-button")
        pageTitleElement.innerText = `Endpoint ${this.endpoint_id}`;
        pageTitleElement.buttonText = 'Delete';
        pageTitleElement.buttonFunction = this.deleteConnectorAction.bind(this);
        this.shawdow.append(pageTitleElement)

        this.shawdow.append(contentTemplate.content.cloneNode(true))
    }
}

customElements.define("wc-endpoint", EndpointPage)