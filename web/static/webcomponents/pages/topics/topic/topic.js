import { fromHTML } from '../../../utils.js'

const topicPageTemplate = document.createElement("template")
topicPageTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/topics/topic/topic.css">
`

class TopicPage extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.topic_id = this.getAttribute("topic_id")
        this.shawdow.append(topicPageTemplate.content.cloneNode(true))
    }

    createTopicExplorer(oldestOffset, latestOffset) {
        //TODO, paginate this somehow....
        var list_items = ''
        for (let i = oldestOffset; i < latestOffset; i++) {
            list_items= list_items.concat(`<li class="messageItem">
                <span>Offset ${i}</span>
                <button class="select-button" data-key="${i}">View</button>
            </li>`);
        };
        let markup = `<div id="topicExplorer">
        <wc-title>Topic Messages</wc-title>
        <ul class="messageList">
        ${list_items}
        </ul>
        </div>`;
        const divElementContainer = fromHTML(markup)
        const selectButtons = divElementContainer.querySelectorAll(".select-button");
        selectButtons.forEach((button) => {
          button.addEventListener("click", (event) => {
            const selectEvent = new CustomEvent("select-message", {
              detail: event.target.dataset.key,
            });
            // TODO: Change this from the shawdow DOM
            this.shawdow.dispatchEvent(selectEvent);
          });
        });
        this.shawdow.append(divElementContainer)
    }

    async getTopicMessage(offset) {
        const endpoint = `/api/v1/topic/${this.topic_id}/events?offset=${offset}&limit=1`
        const response = await fetch(endpoint);
        const jsonData = await response.json();
        let thingData = JSON.stringify(jsonData.data[0], undefined, 2);
        let markup = `<wc-modal id="modalThing">
        <wc-title>Offset ${offset}</wc-title>
        <div class="dataContainer"><pre><code>${thingData}</code></pre></div>
        </wc-modal>`
        let divElementContainer = fromHTML(markup)
        this.shawdow.append(divElementContainer)
    }

    async deleteTopicAction(topic_id = this.topic_id) {
        const url = `/api/v1/topic/${topic_id}`
        await fetch(url, {
            method: "Delete",
            headers: {"Content-Type": "application/json",}
        })
        window.location.replace("/topics");
    }

    newDeleteModal() {
        const modalContainer = document.createElement("wc-modal")
        modalContainer.id = "modalThing"
        const confirmDeletionContainer = document.createElement("wc-delete-alert")
        confirmDeletionContainer.titleText = "Delete Topic";
        confirmDeletionContainer.descriptionText = `Are you sure you want to delete topic ${this.topic_id}?`;
        confirmDeletionContainer.buttonFunction = this.deleteTopicAction.bind(this);
        modalContainer.appendChild(confirmDeletionContainer)
        this.shawdow.appendChild(modalContainer)
    };

    generateMarkup(topic) {
        let job_markup = `
        <div>
        <wc-title>Topic Details</wc-title>
        <ul class="endpointContainerCard">
            <li>
                <p class="titler">Topic Name:<p>
                <p >${topic.topicName}</p>
            </li>
            <li>
                <p class="titler">Next Offset:</p>
                <p>${topic.nextOffset}</p>
            </li>
            <li>
                <p class="titler">Oldest Offset:</p>
                <p>${topic.oldestOffset}</p>
            </li>
            <li>
                <p class="titler">Retention Bytes:</p>
                <p>${topic.retentionBytes}</p>
            </li>
        </ul>
        <div>`
        const div = fromHTML(job_markup);
        return div;
    }

    async connectedCallback(){
        const response = await fetch(`/api/v1/topic/${this.topic_id}`);
        const jsonData = await response.json();

        let element = this.generateMarkup(jsonData.data)
        const contentTemplate = document.createElement("template")
        contentTemplate.content.append(element);

        // Create the title and the delete button
        const pageTitleElement = document.createElement("wc-page-heading-button")
        pageTitleElement.innerText = `Topic ${this.topic_id}`;
        pageTitleElement.buttonText = 'Delete';
        pageTitleElement.buttonFunction = this.newDeleteModal.bind(this);
        this.shawdow.append(pageTitleElement)
        this.shawdow.append(contentTemplate.content.cloneNode(true))

        this.createTopicExplorer(jsonData.data.oldestOffset, jsonData.data.nextOffset)
        this.shawdow.addEventListener("select-message", (event) => {
            this.getTopicMessage(event.detail)
          }
        );
        
        
    }
}

customElements.define("wc-topic", TopicPage)