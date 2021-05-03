// get data incase site is reloaded or accessed through the url
document.addEventListener("DOMContentLoaded", async function() {
  let queriedMod = new URLSearchParams(window.location.search).get('mod');
  document.getElementById('mod-search').value = queriedMod;
  getData(queriedMod);

  // get all mod names from the database
  var response = await fetch('/api/getModlist', { method: "GET" });
  let modData = await response.json();
  modData.forEach(el => {
    let option = document.createElement("option");
    option.innerHTML = el.DisplayName;
    document.getElementById("modlist").appendChild(option);
  })
  document.getElementById("mod-search").setAttribute("list", "modlist");
});

//redirection from searchbar
function search(element) {
  var opts = document.getElementById('modlist').childNodes;
  for (var i = 0; i < opts.length; i++) {
    if (opts[i].value === element.value) {
      getData(opts[i].value);
      break;
    }
  }
}

function parseChatTags(str) {
  let linebr = str.replace(/\\r\\n|\\n/g, "<br>");
  let quot = linebr.replace(/\\"/g, "&quot;");
  let itemtag = quot.replace(/\[i(.*?):(\w+)\]/g, `<img src="assets/Item_$2.png" id="item-icon">`);
  let colortag = itemtag.replace(/\[c\/(\w+):([\s\S]+?)\]/g, `<span style="color: #$1;">$2</span>`);

  return colortag;
}
function linkedModRefs(str) {
  if(str == "") return "no mods";

  let stre = "";
  let mods = str.split(', ');
  for(let item in mods) {
    stre += `<a href="?mod=${mods[item]}">${mods[item]}</a>`;
    if (item != mods.length - 1)
      stre += ', ';
  }
  return stre;
}

async function getData(modName) {
  window.history.pushState({}, null, '?mod=' + modName);

  var response = await fetch(`/api/getModInfo?modname=${modName}`);
  let modData = await response.json();

  if (response.status == 200) {
    let html = `<div>
      <div id="mod-info">
        <img src="https://mirror.sgkoi.dev/direct/${modData.InternalName}.png" id="icon" width="160px" height="160px" style="display: ${modData.Icon !== "" ? "block" : "none"}"></img>
        <p>Display name: <span id="displayName">${parseChatTags(modData.DisplayName)}</span></p>
        <p>Internal name: <span id="internalName">${modData.InternalName}</span></p>
        <p>Version: <span id="version">${modData.Version} (${modData.TModLoaderVersion})</span></p>
        <p>Author: <span id="author">${modData.Author}</span></p>
        <p>Homepage: <span id="homepage">${modData.Homepage != "no homepage" ? `<a href="${modData.Homepage}" target="_blank">${modData.Homepage}</a>` : `${modData.Homepage}`}</span></p>
        <p>Last updated: <span id="updated">${modData.LastUpdated}</span></p>
        <p>Widget url: <span id="widget">${'<a href="https://bettermodwidget.javidpack.repl.co/?mod=' + modData.name + '" target="_blank">' + 'https://bettermodwidget.javidpack.repl.co/?mod=' + modData.name + '</a>'}</span></p>
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
        <p>Downloads today: <span id="dl-today">no Data</span></p>
        <p>Downloads yesterday: <span id="dl-yesterday">${modData.DownloadsYesterday}</span></p>
        <br>
        <p>Popularity rank: <span id="pop-rank">no Data</span></p>
      </div>
      <div id="dl-history">
        <embed src="http://javid.ddns.net/tModLoader/tools/moddownloadhistory.php?modname=${modData.InternalName}" style="width:800px; height: 400px;">
      </div>
    </div>`;

    document.getElementById("content").innerHTML = html;

    document.getElementById("content").style.display = "block";
    document.getElementById('oopsText').style.display = "none";
    document.getElementById("title-text").innerHTML = parseChatTags(modData.DisplayName);

    let description = JSON.stringify(modData.Description);
    document.getElementById("description").innerHTML = parseChatTags(description.substr(1, description.length - 2));
  }
  else {
    document.getElementById('mod-search').value = 'Invalid Request';
    document.getElementById('oopsText').style.display = "block";
    document.getElementById("content").style.display = "none";
    document.getElementById("title-text").innerHTML = 'Invalid';
  }
}