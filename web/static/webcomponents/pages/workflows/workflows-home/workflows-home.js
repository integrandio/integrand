import { fromHTML } from '../../../utils.js'

const jobsHomeTemplate = document.createElement("template")
jobsHomeTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/workflows/workflows-home/workflows-home.css">
`

class WorkflowsHome extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(jobsHomeTemplate.content.cloneNode(true))
    }

    newWorkflow(e) {
        e.preventDefault()
        const modal = this.shawdow.querySelector('#modalThing');
        const cardsContainer = this.shawdow.querySelector("#cardContainer");
        const data = new FormData(e.target);
        const value = Object.fromEntries(data.entries());
        fetch('/api/v1/workflow', {
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
            throw new Error('Something went wrong');
        }).then(glueResponseData => {
            const element = this.generateEndpointCard(glueResponseData.data)
            cardsContainer.prepend(element)
            const successMessage = `<h1>Connector Successfully Created</h1>`
            const modal_element = fromHTML(successMessage);
            modal.innerHTML = '';
            modal.appendChild(modal_element)
        }).catch((error) => {
            const errorMessage = `<h1>Unable to Create Connector</h1>`
            const modal_element = fromHTML(errorMessage);
            modal.innerHTML = '';
            modal.appendChild(modal_element)
        });
    }

    newWorkflowAction() {
        const modalMarkup = `
            <wc-modal id="modalThing">
                <wc-title>Create New Workflow</wc-title>
                <form id="myForm">
                  <label for="id">Topic Name:</label><br>
                  <input type="text" id="topicName" name="topicName" value=""><br>
                  <label for="topicName">Function Name:</label><br>
                  <input type="text" id="functionName" name="functionName" value="">
                  <br>
                  <input type="submit" value="Create">
                </form>
            </wc-modal>`
        const modal_element = fromHTML(modalMarkup);
        this.shawdow.append(modal_element)
        const formComponent = this.shawdow.querySelector('#myForm');
        formComponent.addEventListener('submit', this.newConnector.bind(this));
    };
    
    generateWorkflowCard(workflow) {
        const endpoint_link = `/workflows/${endpoint.id}`
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
        const workflowResponse = await fetch('/api/v1/workflow');
        const workflowResponseData = await glueResponse.json();

        const workflow_card_container = document.createElement("data-cards-container")
        workflow_card_container.id = "cardContainer"
        for (const workflow of workflowResponseData.data) {
            const card_element = this.generateWorkflowCard(workflow)
            workflow_card_container.appendChild(card_element)
        }
        return workflow_card_container
    }

    async connectedCallback(){
        const pageTitleElement = document.createElement("wc-page-heading-button")
        pageTitleElement.innerText = "All Workflows";
        pageTitleElement.buttonText = 'New Workflow';
        pageTitleElement.buttonFunction = this.newWorkflowAction.bind(this);
        this.shawdow.appendChild(pageTitleElement)

        const endpoint_card_container = await this.generateEndpointsContainer()
        this.shawdow.appendChild(endpoint_card_container)
    }
}

customElements.define("workflows-home", WorkflowsHome)