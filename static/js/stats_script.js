// get data incase site is reloaded or accessed through the url
document.addEventListener("DOMContentLoaded", async function () {
  let queriedMod = new URLSearchParams(window.location.search).get('mod');
  document.getElementById('mod-search').value = queriedMod;
  getData(queriedMod);

  let modList = await setupModList();
});

function parseChatTags(str) {
  let lt = str.replace(/</g, "&lt;");
  let gt = lt.replace(/>/g, "&gt;");
  let linebr = gt.replace(/\\r\\n|\\n/g, "<br>"); // replace \r\n and \n with <br> tags
  let tab = linebr.replace(/\\t/g, "    "); // replace \t with 4 spaces
  let apo = tab.replace(/\\'/g, "'"); // replace \' with '
  let backslash = apo.replace(/\\/g, "\\"); // replace \\ with \
  let quot = backslash.replace(/\\"/g, "&quot;"); // replace " with &quot;
  let itemtag = quot.replace(/\[i(.*?):(\w+)\]/g, `<img src="https://tmlapis.repl.co/img/Item_$2.png" id="item-icon">`); // replace item tags with the correct image
  let colortag = itemtag.replace(/\[c\/(\w+):([\s\S]+?)\]/g, `<span style="color: #$1;">$2</span>`); // replace the color tags with <span>

  return colortag;
}

function linkedModRefs(str) {
  if (str == "") return "no mods";

  let stre = "";
  let mods = str.split(', ');
  for (let item in mods) {
    // add links
    stre += `<a href="?mod=${mods[item]}">${mods[item]}</a>`;
    // seperate with commas
    if (item != mods.length - 1)
      stre += ', ';
  }
  return stre;
}

async function getData(modName) {
  // convert to internal name if it isn't already
  let resp = await fetch("/api/getInternalName?displayname=" + encodeURIComponent(modName));
  let name = await resp.json(); // name is empty if it already is a internal name or if it doesn't exist
  if (name != "") {
    modName = name;
  }

  window.history.pushState({}, null, '?mod=' + modName);
  
  var response = await fetch(`/api/getModInfo?modname=${modName}`);
  if (response.status === 200) {
    let modData = await response.json();
    displayData(modData);
  }
  else {
    document.getElementById('mod-search').value = 'Invalid Request';
    document.getElementById('oopsText').style.display = "block";
    document.getElementById("content").style.display = "none";
    document.getElementById("title-text").innerHTML = 'Invalid';
  }
}

async function displayData(modData) {
  let html = `<div>
  <div id="mod-info">
    <img src="https://mirror.sgkoi.dev/direct/${modData.InternalName}.png" id="icon" width="160px" height="160px" style="display: ${modData.Icon !== "" ? "block" : "none"}"></img>
    <p>Display name: <span id="displayName">${parseChatTags(modData.DisplayName)}</span></p>
    <p>Internal name: <span id="internalName">${modData.InternalName}</span></p>
    <p>Version: <span id="version">${modData.Version} (${modData.TModLoaderVersion})</span></p>
    <p>Author: <span id="author">${modData.Author}</span></p>
    <p>Homepage: <span id="homepage">${modData.Homepage != "" ? `<a href="${modData.Homepage}" target="_blank">${modData.Homepage}</a>` : "no homepage"}</span></p>
    <p>Last updated: <span id="updated">${modData.LastUpdated}</span></p>
    <p>Widget url: <span id="widget"><a href="https://bettermodwidget.javidpack.repl.co/?mod=${modData.InternalName}" target="_blank">https://bettermodwidget.javidpack.repl.co/?mod=${modData.InternalName}</a></span></p>
    <p>Mod dependencies: <span>${linkedModRefs(modData.ModDependencies)}</span>
    <p>Mod Side: <span>${modData.ModSide}</span></p>
  </div>
  <div id="description-container">
    <h1>Description</h1>
    <p id="description">Loading...</p>
  </div>
  <div id="download-info">
    <h1>Downloads: </h1>
    <p>Download link: <span id="dl-link">${modData.DownloadLink}</span></p>
    <p>Downloads total: <span id="dl-total">${modData.DownloadsTotal}</span></p>
    <p>Downloads today: <span id="dl-today">${modData.DownloadsToday}</span></p>
    <p>Downloads yesterday: <span id="dl-yesterday">${modData.DownloadsYesterday}</span></p>
    <br>
    <p>Rank total: <span id="total-rank">${modData.Rank}</span></p>
  </div>
  <div id="dl-history">
    <h1>Version history</h1>
    <div id="history-table">
      <table>
        <tbody id="history-tbody">
          <tr>
            <td>Version</td><td>Downloads</td><td>TModLoaderVersion</td><td>PublishDate</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</div>`;

  document.getElementById("content").innerHTML = html;

  document.getElementById("content").style.display = "block";
  document.getElementById('oopsText').style.display = "none";
  document.getElementById("title-text").innerHTML = parseChatTags(modData.DisplayName);

  let description = JSON.stringify(modData.Description);
  document.getElementById("description").innerHTML = parseChatTags(description.substr(1, description.length - 2));

  let response = await fetch(`/api/getVersionHistory?modname=${modData.InternalName}`);
  let versionHistory = await response.json();
  makeHistoryGraph(versionHistory);
}

function makeHistoryGraph(versionHistory) {
  let table = document.getElementById("history-tbody");

  for (let i = 0; i < versionHistory.length; i++) {
    let tr = document.createElement("tr");
    tr.innerHTML = `<td>${versionHistory[i].Version}</td><td>${versionHistory[i].Downloads}</td><td>${versionHistory[i].TModLoaderVersion}</td><td>${versionHistory[i].PublishDate}</td>`
    table.appendChild(tr);
  }
}