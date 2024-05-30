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

    deleteTopicAction() {
        console.log("deleting topic")
        const url = `/api/v1/topic/${this.topic_id}`
        fetch(url, {
            method: "Delete",
            headers: {"Content-Type": "application/json",}
        }).then(res => {
            // Check response to see if it's bad
            res.json().then((topicResponseData) => {
                console.log(topicResponseData)
                window.location.replace("/app/topics");
            });
        })
    }

    generateMarkup(topic) {
        let job_markup = `
        <ul class="endpointContainerCard">
            <li>
                <p class="titler">Topic Name:<p>
                <p >${topic.topicName}</p>
            </li>
            <li>
                <p class="titler">Latest Offset:</p>
                <p>${topic.latestOffset}</p>
            </li>
            <li>
                <p class="titler">Oldest Offset:</p>
                <p>${topic.oldestOffset}</p>
            </li>
            <li>
                <p class="titler">Retention Bytes:</p>
                <p>${topic.retentionBytes}</p>
            </li>
        </ul>`
        const div = fromHTML(job_markup);
        return div;
    }

    async connectedCallback(){
        const response = await fetch(`/api/v1/topic/${this.topic_id}`);
        const jsonData = await response.json();
        console.log(jsonData)
        let element = this.generateMarkup(jsonData.data)
        const contentTemplate = document.createElement("template")
        contentTemplate.content.append(element);

        const pageTitleElement = document.createElement("wc-page-heading-button")
        pageTitleElement.innerText = `Topic ${this.topic_id}`;
        pageTitleElement.buttonText = 'Delete';
        pageTitleElement.buttonFunction = this.deleteTopicAction.bind(this);
        this.shawdow.append(pageTitleElement)

        this.shawdow.append(contentTemplate.content.cloneNode(true))
    }
}

customElements.define("wc-topic", TopicPage)