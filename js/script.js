//redirection from searchbar
function search(element) {
  if (event.key === 'Enter') {
    window.location.href = "https://modstats.repl.co/stats.html?mod=" + element.value;
  }
}