import { fromHTML } from '../../../utils.js'

const jobTemplate = document.createElement("template")
jobTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/workflows/workflow/workflow.css">
`

class WorkflowPage extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.workflow_id = this.getAttribute("workflow_id")
        this.shawdow.append(jobTemplate.content.cloneNode(true))
    }

    async deleteWorkflowAction(workflow_id = this.workflow_id) {
        const url = `/api/v1/workflow/${workflow_id}`
        await fetch(url, {
            method: "Delete",
            headers: {"Content-Type": "application/json",}
        })
        window.location.replace("/workflows");
    }

    newDeleteModal() {
        const modalContainer = document.createElement("wc-modal")
        modalContainer.id = "modalThing"
        const confirmDeletionContainer = document.createElement("wc-delete-alert")
        confirmDeletionContainer.titleText = "Delete Workflow";
        confirmDeletionContainer.descriptionText = `Are you sure you want to delete workflow ${this.workflow_id}?`;
        confirmDeletionContainer.buttonFunction = this.deleteWorkflowAction.bind(this);
        modalContainer.appendChild(confirmDeletionContainer)
        this.shawdow.appendChild(modalContainer)
    };

    generateMarkup(workflow) {
        const topic_link = `/topics/${workflow.topicName}`
        let workflow_markup = `
            <ul class="endpointContainerCard">
                <li>
                    <p class="titler">Function Name:</p>
                    <p>${workflow.functionName}</p>
                </li>
                <li>
                    <p class="titler">Offset:</p>
                    <p>${workflow.offset}</p>
                </li>
                <li>
                    <p class="titler">Topic Name:</p>
                    <a class="link" href=${topic_link}>${workflow.topicName}</a>
                </li>
                <li>
                    <p class="titler">Sink URL:</p>
                    <p>${workflow.sinkURL}</p>
                </li>
                <li>
                    <p class="titler">Enabled:</p>
                    <p>${workflow.enabled}</p>
                </li>
            </ul>`
        const div = fromHTML(workflow_markup);
        return div;
    }

    async connectedCallback(){
        const response = await fetch(`/api/v1/workflow/${this.workflow_id}`);
        const jsonData = await response.json();
        let element = this.generateMarkup(jsonData.data)
        const contentTemplate = document.createElement("template")
        contentTemplate.content.append(element)

        // Create the title and the delete button
        const pageTitleElement = document.createElement("wc-page-heading-button")
        pageTitleElement.innerText = `Workflow ${this.workflow_id}`;
        pageTitleElement.buttonText = 'Delete';
        pageTitleElement.buttonFunction = this.newDeleteModal.bind(this);
        this.shawdow.append(pageTitleElement)

        this.shawdow.append(contentTemplate.content.cloneNode(true))
    }
}

customElements.define("wc-workflow", WorkflowPage)