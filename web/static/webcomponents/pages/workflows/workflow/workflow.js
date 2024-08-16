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

    async deleteWorkflowAction(evt) {
        const workflow_id = evt.currentTarget.workflow_param_id
        const url = `/api/v1/workflow/${workflow_id}`
        await fetch(url, {
            method: "Delete",
            headers: {"Content-Type": "application/json",}
        })
        window.location.replace("/workflows");
    }

    newDeleteModal() {
        const modalMarkup = `
            <wc-modal id="modalThing">
                <wc-title>Confirm Deletion</wc-title>
                <p>Are you sure you want to delete workflow ${this.workflow_id}?<p>
                <form id="myForm">
                  <input class="submit-button" type="submit" value="Confirm">
                </form>
            </wc-modal>`
        const modal_element = fromHTML(modalMarkup);
        this.shawdow.append(modal_element)
        const formComponent = this.shawdow.querySelector('#myForm');
        formComponent.addEventListener('submit', this.deleteWorkflowAction);
        formComponent.workflow_param_id = this.workflow_id
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