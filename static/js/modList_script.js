document.addEventListener("DOMContentLoaded", async function () {
    let modList = await setupModList();
    createTable(modList);
});

function createTable(modList) {
    modList = modList.sort((a,b) => a.Rank > b.Rank ? 1 : (a.Rank < b.Rank) ? -1 : 0);

    let table = document.createElement("table");
    table.innerHTML = "<tr> <td>Rank</td> <td>Display Name</td> <td>Internal Name</td> <td>Downloads Total</td> <td>Downloads Today</td> <td>Downloads Yesterday</td> <td>tModloader Version</td> </tr>"
    modList.forEach(el => {
        let tr = document.createElement("tr");

        // extremely bad code lol
        let td = document.createElement("td");
        td.innerHTML = el.Rank;
        tr.appendChild(td);

        let td1 = document.createElement("td");
        td1.innerHTML = `<a href="/stats?mod=${el.ModName}">${el.DisplayName.length >= 50 ?  el.DisplayName.substr(0, 50) + "..." : el.DisplayName}</a>`;
        tr.appendChild(td1);

        let td2 = document.createElement("td");
        td2.innerHTML = `<a href="/stats?mod=${el.ModName}">${el.ModName}</a>`;
        tr.appendChild(td2);

        let td3 = document.createElement("td");
        td3.innerHTML = el.DownloadsTotal;
        tr.appendChild(td3);

        let td4 = document.createElement("td");
        td4.innerHTML = el.DownloadsToday;
        tr.appendChild(td4);

        let td5 = document.createElement("td");
        td5.innerHTML = el.DownloadsYesterday;
        tr.appendChild(td5);

        let td6 = document.createElement("td");
        td6.innerHTML = el.TModLoaderVersion;
        tr.appendChild(td6);

        table.appendChild(tr)
    });

    document.getElementById("loading-text").style.display = "none";
    document.getElementById("content").appendChild(table);
}