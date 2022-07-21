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