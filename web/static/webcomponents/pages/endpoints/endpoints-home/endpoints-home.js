import { fromHTML } from '../../../utils.js'

const endpointssHomeTemplate = document.createElement("template")
endpointssHomeTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/endpoints/endpoints-home/endpoints-home.css">
`

class EndpointsHome extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(endpointssHomeTemplate.content.cloneNode(true))
        this.pendingRequest = false
    }

    newConnector(e) {
        e.preventDefault()
        // Check if a request is being processed by this same client
        if (this.pendingRequest) {
            return
        } else {
            this.pendingRequest = true
        }
        
        const modal = this.shawdow.querySelector('#modalThing');
        const cardsContainer = this.shawdow.querySelector("#cardContainer");
        const data = new FormData(e.target);
        const value = Object.fromEntries(data.entries());
        fetch('/api/v1/connector', {
            method: "POST",
            mode: 'cors',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(value),
        }).then(res => {
            if (res.ok) {
                return res.json();
            }
            this.pendingRequest = false
            throw new Error('Something went wrong');
        }).then(glueResponseData => {
            const element = this.generateEndpointCard(glueResponseData.data)
            cardsContainer.prepend(element)
            const successMessage = `<h1>Connector Successfully Created</h1>`
            const modal_element = fromHTML(successMessage);
            modal.innerHTML = '';
            modal.appendChild(modal_element)
            this.pendingRequest = false
        }).catch((error) => {
            const errorMessage = `<h1>Unable to Create Connector</h1>`
            const modal_element = fromHTML(errorMessage);
            modal.innerHTML = '';
            modal.appendChild(modal_element)
            this.pendingRequest = false
        });
    }

    newConnectionAction() {
        const modalMarkup = `
            <wc-modal id="modalThing">
                <wc-title>Create New Connector</wc-title>
                <form id="myForm">
                  <label for="id">id:</label><br>
                  <input type="text" id="id" name="id" value=""><br>
                  <label for="topicName">Topic Name:</label><br>
                  <input type="text" id="topicName" name="topicName" value="">
                  <br>
                  <input class="submit-button" type="submit" value="Create">
                </form>
            </wc-modal>`
        const modal_element = fromHTML(modalMarkup);
        this.shawdow.append(modal_element)
        const formComponent = this.shawdow.querySelector('#myForm');
        formComponent.addEventListener('submit', this.newConnector.bind(this));
    };
    
    generateEndpointCard(endpoint) {
        const endpoint_link = `/endpoints/${endpoint.id}`
        let endpoint_markup = `<div class="jobCard">
        <h1><span class="titler">ID:</span> ${endpoint.id}</h1>
        <h2><span class="titler">Security Key:</span> ${endpoint.securityKey}</h2>
        <p><span class="titler">Topic Name:</span> ${endpoint.topicName}</p>
        <a class="jobLink" href="${endpoint_link}"> View Endpoint Details </a>
        </div>`
        const card_element = fromHTML(endpoint_markup);
        return card_element;
    }

    async generateEndpointsContainer () {
        const glueResponse = await fetch('/api/v1/connector');
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
        pageTitleElement.innerText = "All Endpoints";
        pageTitleElement.buttonText = 'New Endpoint';
        pageTitleElement.buttonFunction = this.newConnectionAction.bind(this);
        this.shawdow.appendChild(pageTitleElement)

        const endpoint_card_container = await this.generateEndpointsContainer()
        this.shawdow.appendChild(endpoint_card_container)
    }
}

customElements.define("endpoints-home", EndpointsHome)