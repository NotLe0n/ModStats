class DateTime extends HTMLElement {
    constructor() {
        super();

        this.attachShadow({ mode: "open" });

        setTimeout(() => {
            let time = document.createElement('time')
            time.setAttribute('datetime', "2008-02-14 20:00");
            time.innerHTML = new Date(this.innerHTML * 1000).toLocaleString('de-DE')

            this.shadowRoot.append(time);
        })



    }
}

customElements.define("date-time", DateTime);