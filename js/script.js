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
  if (event.key === 'Enter') {
    window.location.href = "https://modstats.repl.co/stats.html?mod=" + element.value;
  }
}
