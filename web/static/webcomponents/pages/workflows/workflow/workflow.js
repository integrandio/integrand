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

    deleteWorkflowAction() {
        const url = `/api/v1/workflow/${this.workflow_id}`
        fetch(url, {
            method: "Delete",
            headers: {"Content-Type": "application/json",}
        }).then(res => {
            // Check response to see if it's bad
            res.json().then((endpointResponseData) => {
                console.log(endpointResponseData)
                window.location.replace("/workflows");
            });
        })
    }

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
                    <p>${workflow.Offset}</p>
                </li>
                <li>
                    <p class="titler">Topic Name:</p>
                    <a class="link" href=${topic_link}>${endpoint.topicName}</a>
                </li>
                <li>
                    <p class="titler">Sink URL:</p>
                    <p>${workflow.SinkURL}</p>
                </li>
                <li>
                    <p class="titler">Enabled:</p>
                    <p>${workflow.Enabled}</p>
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
        pageTitleElement.buttonFunction = this.deleteWorkflowAction.bind(this);
        this.shawdow.append(pageTitleElement)

        this.shawdow.append(contentTemplate.content.cloneNode(true))
    }
}

customElements.define("wc-workflow", WorkflowPage)