document.addEventListener("DOMContentLoaded", async function () {
  // get all mod names from the database
  var response = await fetch('/api/getModlist', {
    method: "GET"
  });
  let modData = await response.json();
  modData.forEach(el => {
    let option = document.createElement("option");
    option.innerHTML = el.DisplayName;
    document.getElementById("modlist").appendChild(option);
  })
  document.getElementById("mod-search").setAttribute("list", "modlist");
});

function search(element) {
  if (event.keyCode === 13) {
    window.location.href = `/stats?mod=${element.value}`;
  }
}