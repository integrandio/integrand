import { fromHTML } from '../../../utils.js'

const topicsHomeTemplate = document.createElement("template")
topicsHomeTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/topics/topics-home/topics-home.css">
`

class TopicsHome extends HTMLElement {
    constructor(){
        super()
        this.shawdow = this.attachShadow({mode: "open"})
        this.shawdow.append(topicsHomeTemplate.content.cloneNode(true))
    }

    newTopicAction() {
        let cardsContainer = this.shawdow.querySelector("#cardContainer");
        fetch('/api/v1/topic', {
            method: "POST",
            headers: {"Content-Type": "application/json",}
        }).then(res => {
            // TODO: add element to the page 
            res.json().then((topicResponseData) => {
                console.log(topicResponseData)
                const element = this.generateEndpointCard(topicResponseData.data)
                cardsContainer.prepend(element)
            });
        })
    }


    generateEndpointCard(topic) {
        const topic_link = `/app/topics/${topic.topicName}`
        let topic_markup = `
        <div class="jobCard">
            <h1><span class="titler">Name:</span> ${topic.topicName}</h1>
            <h2><span class="titler">Latest Offset:</span> ${topic.latestOffset}</h2>
            <p><span class="titler">Oldest Offset: </span> ${topic.oldestOffset}</p>
            <a class="jobLink" href="${topic_link}"> View Topic Details </a>
        </div>`
        const card_element = fromHTML(topic_markup);
        return card_element;
    }

    async generateEndpointsContainer () {
        const topicResponse = await fetch('/api/v1/topic');
        const topicJsonData = await topicResponse.json();
        console.log(topicJsonData)

        const endpoint_card_container = document.createElement("data-cards-container")
        endpoint_card_container.id = "cardContainer"
        for (const topic of topicJsonData.data) {
            const card_element = this.generateEndpointCard(topic)
            endpoint_card_container.appendChild(card_element)
        }
        //let button_markup = this.generateNewEndpointButton()
        //endpoint_card_container.appendChild(button_markup)
        return endpoint_card_container
    }

    async connectedCallback(){
        const pageTitleElement = document.createElement("wc-page-heading-button")
        pageTitleElement.innerText = "All Topics";
        pageTitleElement.buttonText = 'New Topic';
        pageTitleElement.buttonFunction = this.newTopicAction.bind(this);
        this.shawdow.appendChild(pageTitleElement)

        const topics_card_container = await this.generateEndpointsContainer()
        this.shawdow.appendChild(topics_card_container)
    }
}

customElements.define("topics-home", TopicsHome)