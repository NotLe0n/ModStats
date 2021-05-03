document.addEventListener("DOMContentLoaded", async function() {
  // get all mod names from the database
  var response = await fetch('/api/getModlist', { method: "GET" });
  let modData = await response.json();
  console.log(modData)
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
      window.location.href = "/stats?mod=" + opts[i].value;
      break;
    }
  }
}