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
        this.pendingRequest = false;
        this.workflowfunctions = []
        this.topicNames = []
    }

    newWorkflow(e) {
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
            this.pendingRequest = false
            throw new Error('Something went wrong');
        }).then(workflowResponse => {
            const element = this.generateWorkflowCard(workflowResponse.data)
            cardsContainer.prepend(element)
            const successMessage = `<h1>Workflow Successfully Created</h1>`
            const modal_element = fromHTML(successMessage);
            modal.innerHTML = '';
            modal.appendChild(modal_element)
            this.pendingRequest = false
        }).catch((error) => {
            const errorMessage = `<h1>Unable to Create Workflow</h1>`
            const modal_element = fromHTML(errorMessage);
            modal.innerHTML = '';
            modal.appendChild(modal_element)
            this.pendingRequest = false
        });
    }

    newWorkflowAction() {
        const functionOptions = this.workflowfunctions.map((func) => {return `<option value="${func}">${func}</option>`})
        const topicOptions = this.topicNames.map((topic) => {return `<option value="${topic}">${topic}</option>`})

        const modalMarkup = `
            <wc-modal id="modalThing">
                <wc-title>Create New Workflow</wc-title>
                <form id="myForm">
                  <label for="topicName">Topic Name:</label><br>
                  <select id="topicName" name="topicName">${topicOptions.join()}</select><br>
                  <label for="functionName">Function Name:</label><br>
                  <select id="functionName" name="functionName">${functionOptions.join()}</select><br>
                  <label for="sinkURL">Sink Url:</label><br>
                  <input type="text" id="sinkURL" name="sinkURL" value=""><br>
                  <br>
                  <input type="submit" value="Create">
                </form>
            </wc-modal>`
        const modal_element = fromHTML(modalMarkup);
        this.shawdow.append(modal_element)
        const formComponent = this.shawdow.querySelector('#myForm');
        formComponent.addEventListener('submit', this.newWorkflow.bind(this));
    };
    
    generateWorkflowCard(workflow) {
        const workflow_link = `/workflows/${workflow.id}`
        let workflow_markup = `<div class="jobCard">
        <h1><span class="titler">ID:</span> ${workflow.id}</h1>
        <p><span class="titler">Topic Name:</span> ${workflow.topicName}</p>
        <h2><span class="titler">Function Name:</span> ${workflow.functionName}</h2>
        <a class="jobLink" href="${workflow_link}"> View Workflow Details </a>
        </div>`
        const card_element = fromHTML(workflow_markup);
        return card_element;
    }

    async generateWorkflowsContainer () {
        const workflowResponse = await fetch('/api/v1/workflow');
        const workflowResponseData = await workflowResponse.json();

        const workflow_card_container = document.createElement("data-cards-container")
        workflow_card_container.id = "cardContainer"
        for (const workflow of workflowResponseData.data) {
            const card_element = this.generateWorkflowCard(workflow)
            workflow_card_container.appendChild(card_element)
        }
        return workflow_card_container
    }

    async getWorkflowFunctions() {
        const workflowFunctionsResponse = await fetch('/api/v1/workflow/functions');
        const workflowFunctionsResponseData = await workflowFunctionsResponse.json();
        this.workflowfunctions = workflowFunctionsResponseData.data
    }

    async getTopicNames() {
        const topicsResponse = await fetch('/api/v1/topic');
        const topicsResponseData = await topicsResponse.json();
        const topicNames = topicsResponseData.data.map((topic) => {return topic.topicName})
        this.topicNames = topicNames
    }

    async connectedCallback(){
        // Async process workflow functions and topics
        this.getWorkflowFunctions()
        this.getTopicNames()
        const pageTitleElement = document.createElement("wc-page-heading-button")
        pageTitleElement.innerText = "All Workflows";
        pageTitleElement.buttonText = 'New Workflow';
        pageTitleElement.buttonFunction = this.newWorkflowAction.bind(this);
        this.shawdow.appendChild(pageTitleElement)

        const workflow_card_container = await this.generateWorkflowsContainer()
        this.shawdow.appendChild(workflow_card_container)
    }
}

customElements.define("workflows-home", WorkflowsHome)