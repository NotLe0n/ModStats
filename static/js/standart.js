async function setupModList() {
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

    return modList;
};

async function search(element) {
    // has enter been pressed?
    if (event.keyCode === 13) {
        if (element.value.startsWith("author:")){
            window.location.href = `/stats?author=${encodeURIComponent(element.value.substr(7).trim())}`;
        }
        else {
            window.location.href = `/stats?mod=${encodeURIComponent(element.value)}`;
        }
    }
}