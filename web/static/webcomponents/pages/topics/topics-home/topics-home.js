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

    newTopic(e) {
        e.preventDefault()
        const modal = this.shawdow.querySelector('#modalThing');
        const cardsContainer = this.shawdow.querySelector("#cardContainer");
        const data = new FormData(e.target);
        const value = Object.fromEntries(data.entries());
        fetch('/api/v1/topic', {
            method: "POST",
            mode: 'cors',
            headers: {"Content-Type": "application/json",},
            body: JSON.stringify(value),
        }).then(res => {
            if (res.ok) {
                return res.json();
              }
            throw new Error('Something went wrong');
        }).then(topicResponseData => {
            const card_element = this.generateEndpointCard(topicResponseData.data)
            cardsContainer.prepend(card_element)
            const successMessage = `<h1>Topic Successfully Created</h1>`
            const modal_element = fromHTML(successMessage);
            modal.innerHTML = '';
            modal.appendChild(modal_element)
        }).catch((error) => {
            const errorMessage = `<h1>Unable to Create Topic</h1>`
            const modal_element = fromHTML(errorMessage);
            modal.innerHTML = '';
            modal.appendChild(modal_element)
        });
    }

    newTopicAction() {
        const modalMarkup = `
        <wc-modal id="modalThing">
            <wc-title>Create New Connector</wc-title>
            <form id="myForm">
              <label for="topicName">Topic Name:</label><br>
              <input type="text" id="topicName" name="topicName" value="">
              <br>
              <input type="submit" value="Create">
            </form>
        </wc-modal>`

        const modal_element = fromHTML(modalMarkup);
        this.shawdow.append(modal_element)
        const formComponent = this.shawdow.querySelector('#myForm');
        formComponent.addEventListener('submit', this.newTopic.bind(this));
    }


    generateEndpointCard(topic) {
        const topic_link = `/topics/${topic.topicName}`
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