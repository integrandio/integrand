class ReportTable extends HTMLTableElement {
    constructor() {
        super()
        this.base_api_endpoint = '/static/'
        this.filterState = {};
        this.innerHTML = `
        <link rel="stylesheet" type="text/css" href="/static/reset.css">
        <link rel="stylesheet" type="text/css" href="/static/webcomponents/pages/reports/report-table/report-table.css">
        `
    }

    dataFromRow = (row, headers) => {
        return Object.fromEntries([...row.querySelectorAll('td')]
            .map((td, index) => [headers[index], td.textContent]));
    }
  
    matchesCriteria = (rowData, filters) =>
        filters.every(([key, value]) => rowData[key] === value);
    
    noti = (e) => {
        console.log(e)
        this.refresh()
    }
    
    refresh = () => {
        // const allButtons = this.querySelectorAll('multi-select-button');
        // allButtons.forEach((but) => {
        //     console.log(but.filteredState)
        // })
        const headers = [...this.querySelectorAll('thead th')].map(th => th.textContent),
            filters = Object.entries(this.filterState),
            showAll = filters.length === 0;
        this.querySelectorAll('tbody tr').forEach(row => {
            const show = showAll || this.matchesCriteria(this.dataFromRow(row, headers), filters);
            row.classList.toggle('hidden-row', !show);
        });
    };
  
    handleFilterChange = (e) => {
        console.log(e)
        const field = e.field,
            value = e.value;
        if (value) {
            this.filterState[field] = value;
        } else {
            delete this.filterState[field];
        }
        console.log(this.filterState)
        this.refresh();
    };

    newHeaderRow(jsonData) {
        // Get the keys (column names) of the first object in the JSON data
        let cols = Object.keys(jsonData[0]);
        //Create our row element
        let row_el = document.createElement('tr')
        row_el.className = 'dataTableRow';
        cols.forEach(col => {
            let row_filter_set = new Set()
            jsonData.forEach((item) => {
                row_filter_set.add(item[col])
            })
            let table_head_el = document.createElement('th')
            table_head_el.className = "dataTableHeadMiddleCol"
            let multi_select_el = document.createElement('multi-select-button');
            multi_select_el.innerText = col
            multi_select_el.colName = col
            multi_select_el.avaliableFilters = row_filter_set
            multi_select_el.notifyer = this.handleFilterChange

            table_head_el.appendChild(multi_select_el)
            row_el.appendChild(table_head_el)
        })
        let thead_el = document.createElement('thead')
        thead_el.appendChild(row_el)
        return thead_el
    }
  
    headerRowString(jsonData) {
        // Get the keys (column names) of the first object in the JSON data
        let cols = Object.keys(jsonData[0]);
        var first_row = "<tr class='dataTableRow'>"
        cols.forEach((item) => {
            let table_header_string = `
            <th scope="col" class="dataTableHeadMiddleCol">${item}</th>`
            first_row = first_row.concat(table_header_string)
        })
        first_row = first_row.concat("</tr>")

        var second_row = "<tr>"
        cols.forEach((col) => {
            let setter = new Set()
            jsonData.forEach((item) => {
                setter.add(item[col])
            })
            var table_header_string = `<th><select class="filter" data-field="${col}"><option value="">None</option>`
            setter.forEach((item) => {
                let option_string = `<option value="${item}">${item}</option>`
                table_header_string = table_header_string.concat(option_string)
            })
            table_header_string = table_header_string.concat("</select></th>")
            second_row = second_row.concat(table_header_string)
        });
        second_row  = second_row.concat("</tr>")
        
        let all_rows = "".concat("<thead>", first_row, second_row, "</thead>")
        return all_rows
    }

    newTableBody(jsonData) {
        let tbody_el = document.createElement("tbody");
        tbody_el.className="dataTableBody";
        jsonData.forEach((item) => {
            let tr_el = document.createElement("tr");
            tr_el.className = "dataTableRow";
            
            let vals = Object.values(item);
            vals.forEach((elem) => {
                let td_el = document.createElement("td");
                td_el.className = "dataTableDataMiddleCol";
                td_el.innerText = elem;
                tr_el.appendChild(td_el)
            })
            tbody_el.appendChild(tr_el)
        })
        return tbody_el
    }

    tableBodyString(jsonData) {
        var table_body_rows_string  = "<tbody class='dataTableBody'>"
        jsonData.forEach((item) => {
            table_body_rows_string = table_body_rows_string.concat("<tr class='dataTableRow'>")
            let vals = Object.values(item);
            vals.forEach((elem) => {
                let table_data_string = `<td class="dataTableDataMiddleCol">${elem}</td>`
                table_body_rows_string = table_body_rows_string.concat(table_data_string)
            })
            table_body_rows_string = table_body_rows_string.concat("</tr>")
        })
        table_body_rows_string = table_body_rows_string.concat("</tbody>")
        return table_body_rows_string
    }

    async connectedCallback(){
        let data_endpoint = this.base_api_endpoint + this.getAttribute("data_url")
        const response = await fetch(data_endpoint);
        const jsonData = await response.json();
        const table_body_rows_string = this.tableBodyString(jsonData)
        await customElements.whenDefined('multi-select-button');
        const headerRowElement = await this.newHeaderRow(jsonData)
        const bodyRowsElement = this.newTableBody(jsonData);
        this.className = "dataTable"
        this.appendChild(headerRowElement)
        this.appendChild(bodyRowsElement)

        // document.querySelectorAll('.filter').forEach(filter =>
        //     filter.addEventListener('change', this.handleFilterChange));
    }
}

customElements.define("report-table", ReportTable, {
    extends: "table"
})