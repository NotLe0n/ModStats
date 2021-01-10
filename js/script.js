document.addEventListener("DOMContentLoaded", async function() {
  // get all mod names from the database
  var response = await fetch('/idk', { method: "POST" });
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
      window.location.href = "https://modstats.repl.co/stats.html?mod=" + opts[i].value;
      break;
    }
  }
}