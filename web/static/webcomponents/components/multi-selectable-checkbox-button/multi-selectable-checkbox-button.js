const msc_button = document.createElement("template");
msc_button.innerHTML = `
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/webcomponents/components/multi-selectable-checkbox-button/multi-selectable-checkbox-button.css">
<div class="multiselect">
  <a id="opener" class="tableHeaderLink">
    <slot></slot>
    <span class="svgSpan">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" class="svger">
        <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd">
        </path>
      </svg>
    </span>
  </a>
  <div id="checkboxes">
  </div>
</div>
`;

class MultiSelectButton extends HTMLElement {
  constructor() {
    super();
    this._colname = ""
    this._filtered_state = new Set();
    this.expanded = false;
    this._shadowRoot = this.attachShadow({mode: "open"})
    this._shadowRoot.append(msc_button.content.cloneNode(true))
    this.openButton = this._shadowRoot.querySelector("#opener");
    this.checkboxes = this._shadowRoot.querySelector('#checkboxes');
    this._avaliable_filterable_values = new Set();
    this._notifyer = function() { return };
  }

  get filteredState() {
    return this._filtered_state
  }

  set notifyer(funcValue) {
    this._notifyer = funcValue
  }

  get avaliableFilters() {
    return this._avaliable_filterable_values
  }

  set avaliableFilters(setValue) {
    this._avaliable_filterable_values = setValue
  }

  set colName(value) {
    this._colname = value
  }

  get colName() {
    return this._colname
  }

  showCheckboxes(e) {
    if (!this.expanded) {
      this.checkboxes.style.display = "block";
      this.expanded = true;
    } else {
      this.checkboxes.style.display = "none";
      this.expanded = false;
    }
  }

  toggleSelection(event) {
    let element_id = event.target.id
    if (this._filtered_state.has(element_id)) {
      this._filtered_state.delete(element_id)
    } else {
      this._filtered_state.add(element_id)
    }
    this._notifyer(
      {
        field : this._colname,
        value: element_id
      }
    )
  }

  createCheckboxElement(name) {
    let label_el = document.createElement('label');
    label_el.for = name
    label_el.innerText = name
    let input_el = document.createElement('input');
    input_el.type = "checkbox"
    input_el.className = "myInput"
    input_el.id = name
    label_el.prepend(input_el)
    return label_el
  }

  connectedCallback() {
    this._avaliable_filterable_values.forEach(value => {
      let check_el = this.createCheckboxElement(value)
      this.checkboxes.appendChild(check_el)
    })

    if (!this.openButton) return;
    this.openButton.addEventListener('click', this.showCheckboxes.bind(this));

    let selectionElements = this._shadowRoot.querySelectorAll('.myInput');
    if (!selectionElements) return;
    selectionElements.forEach(filter =>
      filter.addEventListener('change', this.toggleSelection.bind(this))
    )
  }

  disconnectedCallback () {
    console.log('disconnecting!!')
    // This wont work, we're going to have an object leak
    // if (!this.openButton) return;
    // openButton.removeEventListener('click', this.showCheckboxes);

    // if (!this.selectionElements) return;
    // selectionElements.forEach(filter =>
    //   filter.removeEventListener('change', this.toggle)
    // );
  }
}

customElements.define("multi-select-button", MultiSelectButton);
