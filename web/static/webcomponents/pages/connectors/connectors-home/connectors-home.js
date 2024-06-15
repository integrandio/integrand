import { fromHTML } from '../../../utils.js'

const jobsHomeTemplate = document.createElement("template")
jobsHomeTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/connectors/connectors-home/connectors-home.css">
`

class ConnectorsHome extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(jobsHomeTemplate.content.cloneNode(true))
    }

    newConnectionAction() {
        let cardsContainer = this.shawdow.querySelector("#cardContainer");
        fetch('/api/v1/glue', {
            method: "POST",
            mode: 'cors',
            headers: {"Content-Type": "application/json",},
            body: JSON.stringify({"id": "", "topicName": ""}),
        }).then(res => {
            // TODO: add element to the page 
            res.json().then((glueResponseData) => {
                const element = this.generateEndpointCard(glueResponseData.data)
                cardsContainer.prepend(element)
            });
        })
    }
    
    generateEndpointCard(endpoint) {
        const endpoint_link = `/app/connectors/${endpoint.id}`
        let endpoint_markup = `<div class="jobCard">
        <h1><span class="titler">ID:</span> ${endpoint.id}</h1>
        <h2><span class="titler">Connection Key:</span> ${endpoint.connectionKey}</h2>
        <p><span class="titler">Topic Name:</span> ${endpoint.topicName}</p>
        <a class="jobLink" href="${endpoint_link}"> View Endpoint Details </a>
        </div>`
        const card_element = fromHTML(endpoint_markup);
        return card_element;
    }

    async generateEndpointsContainer () {
        const glueResponse = await fetch('/api/v1/glue');
        const glueJsonData = await glueResponse.json();

        const endpoint_card_container = document.createElement("data-cards-container")
        endpoint_card_container.id = "cardContainer"
        for (const endpoint of glueJsonData.data) {
            const card_element = this.generateEndpointCard(endpoint)
            endpoint_card_container.appendChild(card_element)
        }
        return endpoint_card_container
    }

    async connectedCallback(){
        const pageTitleElement = document.createElement("wc-page-heading-button")
        pageTitleElement.innerText = "All Connectors";
        pageTitleElement.buttonText = 'New Connector';
        pageTitleElement.buttonFunction = this.newConnectionAction.bind(this);
        this.shawdow.appendChild(pageTitleElement)

        const endpoint_card_container = await this.generateEndpointsContainer()
        this.shawdow.appendChild(endpoint_card_container)
    }
}

customElements.define("connectors-home", ConnectorsHome)