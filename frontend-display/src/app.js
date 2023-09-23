console.log("Loaded Simple Worship Display")

// Event Source data updater!
const es = new EventSource("/data-events");
es.onerror = (err) => { console.log("data-events error", err) }
es.onmessage = (msg) => { console.log("data-events message", msg) }
es.onopen = (...args) => { console.log("Connected to Display Updater", args) }
es.addEventListener("data-update", (event) => {
    const parsedData = JSON.parse(event.data);
    updateDisplay(parsedData.data)
});
es.addEventListener("data-app-start", (event) => {
    // on a fresh app start, reload the display view
    window.location.reload()
});

// Load initial data from server on first load
fetch("/data").then((response) => {
    response.json().then((data) => {
        updateDisplay(data)
    })
})


// Handle the data
function updateDisplay(data) {
    console.log("updateDisplay", data)

    if (data.type == "song") {
        document.querySelector("#app").innerHTML = `
            <h1>${data.meta.section}</h1>
            <p>${data.meta.text}</p>
        `
        return
    }

    document.querySelector("#app").innerHTML = JSON.stringify(data)
}