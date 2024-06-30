const searchCssTemplate = document.createElement("template")
searchCssTemplate.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/search/search.css">
<div class="searchBarContainer">
    <label for="search" class="searchLabel">Search</label>
    <div class="searchBarWrapper">
    <div class="svgWrapper">
        <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        aria-hidden="true"
        class="svger"
        >
        <path
            fill-rule="evenodd"
            d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z"
            clip-rule="evenodd"
        ></path>
        </svg>
    </div>
    <input
        id="search"
        name="search"
        class="searchBarInput"
        placeholder="Search"
        type="search"
    />
    </div>
</div>
<div id="resultsDiv" class="resultsDiv">
</div>
`

function isNumeric(value) {
    return /^-?\d+$/.test(value);
}

class Searchbar extends HTMLElement {
    constructor() {
        super()
        const shadow = this.attachShadow({mode: "open"})
        shadow.append(searchCssTemplate.content.cloneNode(true))
        this.expanded = false;
        this.input_content = '';
        this.topic_name = this.getAttribute("topic_name")
        this.resultsContainer = shadow.querySelector('#resultsDiv')
        this.inputSelector = shadow.querySelector('#search')
    }

    async run_search() {
        let isNumber = isNumeric(this.input_content)
        if (!isNumber) {
            return "must be an int"
        }
        const endpoint = `/api/v1/topic/${this.topic_name}/events?offset=${this.input_content}&limit=1`
        const response = await fetch(endpoint);
        if (!response.ok){
            return "not found"
        }
        const jsonData = await response.json();
        console.log(jsonData)
        return 0;
    }

    async showResults() {
        if (this.expanded) {
            let link = await this.run_search()
            this.resultsContainer.style.display = "block";
            this.resultsContainer.innerHTML = link;
        } else {
            this.resultsContainer.style.display = "none";
            this.resultsContainer.innerHTML = "";
        }
    }

    async updateValue(e) {
        this.input_content = e.target.value;
        // This is janky let's clean this up
        if (this.input_content !== "") {
            this.expanded = true;
            this.showResults();
            
        } else {
            this.expanded = false;
            this.showResults();
        }
    }

    async connectedCallback() {
        //Hate having to bind this, let's see if we can clean this up.
        //see this: https://stackoverflow.com/questions/11565471/removing-event-listener-which-was-added-with-bind
        this.inputSelector.addEventListener("input", this.updateValue.bind(this));
    }

    disconnectedCallback() {
    }
}

customElements.define("search-bar", Searchbar)