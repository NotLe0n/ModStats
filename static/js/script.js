document.addEventListener("DOMContentLoaded", async function () {
  let modList = await setupModList();
  updateTMLStats(modList);
});

function updateTMLStats(modList) {
    modList = modList.sort((a, b) => (a.Rank > b.Rank) ? 1 : (a.Rank < b.Rank) ? -1 : 0)

    document.getElementById("modcount").innerHTML = modList.length;
    document.getElementById("deadmods").innerHTML = modList.filter(x => x.DownloadsYesterday < 5).length;
    document.getElementById("median").innerHTML = modList[Math.floor(modList.length / 2)].DownloadsTotal;
    document.getElementById("combined").innerHTML = combinedDownloads(modList);
    document.getElementById("percent").innerHTML = ((combinedDownloads(modList.filter(x => x.Rank <= 10)) / combinedDownloads(modList)) * 100).toPrecision(4);
    document.getElementById("contribs").innerHTML = 80;

    let top10modshtml = "";
    modList.filter(x => x.Rank <= 10).forEach(x => top10modshtml += `<li><a href="stats?mod=${x.ModName}">${x.DisplayName}</a></li>`);
    document.getElementById("Top10Mods").innerHTML = top10modshtml;

    modList = modList.sort((a, b) => (a.DownloadsYesterday < b.DownloadsYesterday) ? 1 : (a.DownloadsYesterday > b.DownloadsYesterday) ? -1 : 0);
    let hotmodshtml = "";
    for (let i = 0; i < 10; i++) {
        hotmodshtml += `<li><a href="stats?mod=${modList[i].ModName}">${modList[i].DisplayName}</a></li>`
    }
    document.getElementById("hotmods").innerHTML = hotmodshtml;
}

async function randomMod() {
  var response = await fetch('/api/getRandomMod');
  let modName = await response.json();
  window.location.href = `/stats?mod=${modName}`
}

/* Set the width of the side navigation to 250px and the left margin of the page content to 250px */
function openNav() {
  document.getElementById("mySidenav").style.width = "250px";
  document.getElementById("main").style.marginLeft = "250px";
  document.getElementById("header").querySelector("h1").style.marginLeft = "210px";
}

/* Set the width of the side navigation to 0 and the left margin of the page content to 0 */
function closeNav() {
  document.getElementById("mySidenav").style.width = "0";
  document.getElementById("main").style.marginLeft = "0";
  document.getElementById("header").querySelector("h1").style.marginLeft = "20px";
}

function combinedDownloads(modList) {
    let count = 0;
    for (let i = 0; i < modList.length; i++) {
        count += modList[i].DownloadsTotal;
    }

    return count;
}