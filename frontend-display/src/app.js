console.log("Loaded Simple Worship Display")

// Event Source data updater!
const es = new EventSource("/data-events");
es.onerror = (err) => { console.log("data-events error", err) }
es.onmessage = (msg) => { console.log("data-events message", msg) }
es.onopen = (...args) => { console.log("Connected to Display Updater", args) }
es.addEventListener("data-update", (event) => {
    const parsedData = JSON.parse(event.data);
    updateDisplay(parsedData)
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
    document.querySelector("#app").innerHTML = JSON.stringify(data)
}