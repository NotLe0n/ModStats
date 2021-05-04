// get data incase site is reloaded or accessed through the url
document.addEventListener("DOMContentLoaded", async function () {
    let queriedMod = new URLSearchParams(window.location.search).get('author');
    document.getElementById('mod-search').value = "author: " + queriedMod;
    getData(queriedMod);

    let modList = await setupModList();
});
  
async function getData(steamid64) {
    let resp = await fetch("/api/getAuthorInfo?steamid64=" + encodeURIComponent(steamid64));

    if (resp.status === 200) {
        let json = await resp.json();
        displayData(json);
    }
    else {
        document.getElementById('mod-search').value = 'Invalid Request';
        document.getElementById('oopsText').style.display = "block";
        document.getElementById("content").style.display = "none";
        document.getElementById("title-text").innerHTML = 'Invalid';
    }
}

function displayData(data) {
    console.log(data)
    let html = `<div>
    <p>Steam Name: <span>${data.SteamName}</span></p>
    <p>Downloads Total: <span>${data.DownloadsTotal}</span></p>
    <p>Downloads Yesterday: <span>${data.DownloadsYesterday}</span></p>
    <p>Mods: <span>${data.Mods.length != 0 ? arrayToTable(data.Mods) : "no mods"}</span></p>
    <p>Maintained Mods: <span>${data.MaintainedMods.length != 0 ? arrayToTable2(data.MaintainedMods) : "no mods"}</span></p>
    </div>`;

    document.getElementById("content").innerHTML = html;

    document.getElementById("content").style.display = "block";
    document.getElementById('oopsText').style.display = "none";
    document.getElementById("title-text").innerHTML = data.SteamName;
}

function arrayToTable(array) {
    let table = document.createElement("table");
    table.innerHTML = "<tr> <td>Rank</td> <td>Display Name</td> <td>Downloads Total</td> <td>Downloads Yesterday</td> </tr>"
    array.forEach(el => {
        let tr = document.createElement("tr");

        // extremely bad code lol
        let td = document.createElement("td");
        td.innerHTML = el.RankTotal;
        tr.appendChild(td);

        let td1 = document.createElement("td");
        td1.innerHTML = `<a href="/stats?mod=${encodeURIComponent(el.DisplayName)}">${el.DisplayName}</a>`;
        tr.appendChild(td1);

        let td2 = document.createElement("td");
        td2.innerHTML = el.DownloadsTotal;
        tr.appendChild(td2);

        let td3 = document.createElement("td");
        td3.innerHTML = el.DownloadsYesterday;
        tr.appendChild(td3);

        table.appendChild(tr);
    });
    return table.outerHTML;
}

function arrayToTable2(array) {
    let table = document.createElement("table");
    table.innerHTML = "<tr> <td>Display Name</td> <td>Downloads Total</td> <td>Downloads Yesterday</td> </tr>"
    array.forEach(el => {
        let tr = document.createElement("tr");

        // extremely bad code lol
        let td1 = document.createElement("td");
        td1.innerHTML = `<a href="/stats?mod=${encodeURIComponent(el.ModName)}">${el.ModName}</a>`;
        tr.appendChild(td1);

        let td2 = document.createElement("td");
        td2.innerHTML = el.DownloadsTotal;
        tr.appendChild(td2);

        let td3 = document.createElement("td");
        td3.innerHTML = el.DownloadsYesterday;
        tr.appendChild(td3);

        table.appendChild(tr);
    });
    return table.outerHTML;
}