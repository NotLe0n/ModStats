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

function search(element) {
    // has enter been pressed?
    if (event.keyCode === 13) {
        window.location.href = `/stats?mod=${encodeURIComponent(element.value)}`;
    }
}