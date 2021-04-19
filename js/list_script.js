document.addEventListener("DOMContentLoaded", async function () {
    let table = document.createElement("table");
    table.id = "table"
    document.getElementById("content").appendChild(table);

    let response = await fetch('/getList', { method: 'POST' });
    const mods = await response.json();

    console.log(mods);
    for (let item in mods) {
        let row = document.createElement("tr");
        row.classList.add("row");

        let modname = document.createElement("td");
        modname.innerHTML = linkedMod(mods[item].name);

        let dlCount = document.createElement("td");
        dlCount.innerHTML = mods[item].downloads;

        row.appendChild(modname);
        row.appendChild(dlCount);
        table.appendChild(row);
    }
});

function linkedMod(str) {
    return `<a href="https://modstats.repl.co/stats.html?mod=${str}">${str}</a>`;
}