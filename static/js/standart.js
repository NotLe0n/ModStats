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