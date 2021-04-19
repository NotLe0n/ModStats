document.addEventListener("DOMContentLoaded", async function () {
  // get all mod names from the database
  var response = await fetch('/getModNames', { method: "POST" });
  let modData = await response.json();
  for (let element in modData) {
      let option = document.createElement("option");
      option.innerHTML = modData[element].displayname;
      document.getElementById("modlist").appendChild(option);
  }

  document.getElementById("mod-search").setAttribute("list", "modlist");
});

//redirection from searchbar
function search(element) {
  var opts = document.getElementById('modlist').childNodes;
  for (var i = 0; i < opts.length; i++) {
      if (opts[i].value === element.value) {
          window.location.href = "/stats.html?mod=" + opts[i].value;
          break;
      }
  }
}

// sidebar stuff
function closemenu() {
  document.getElementById("Sidebar").style.width = "0";
  document.body.style.marginLeft = "0";
}

function openmenu() {
  document.getElementById("Sidebar").style.width = "250px";
  document.body.style.marginLeft = "250px";
}

function getRandomMod() {
  let mods = document.getElementById("modlist").options;
  let rand = Math.floor(Math.random() * Math.floor(mods.length));
  window.location.href = "/stats.html?mod=" + mods[rand].innerHTML;
}
