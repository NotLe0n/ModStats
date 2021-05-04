document.addEventListener("DOMContentLoaded", async function () {
  // get all mod names
  var response = await fetch('/api/getModlist', { method: "GET" });
  let modList = await response.json();

  // add mods to search dropdown
  modList.forEach(el => {
    let option = document.createElement("option");
    option.innerHTML = el.DisplayName;
    document.getElementById("modlist").appendChild(option);
  })
  document.getElementById("mod-search").setAttribute("list", "modlist");

  updateTMLStats(modList);
});

function search(element) {
  // has enter been pressed?
  if (event.keyCode === 13) {
    window.location.href = `/stats?mod=${encodeURIComponent(element.value)}`;
  }
}

function updateTMLStats(modList) {
  document.getElementById("modcount").innerHTML = modList.length;
  document.getElementById("deadmods").innerHTML = modList.filter(x => x.DownloadsYesterday < 5).length;
  document.getElementById("median").innerHTML = 4896;
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