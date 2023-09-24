console.log("Loaded Simple Worship Display")


socket = new WebSocket("ws://"+location.host+"/ws");
console.log("Attempting Connection to Websocket...");

socket.onopen = () => {
    console.log("Successfully Connected");
    socket.send("Hi From the Client!")
};

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    socket.send("Client Closed!")
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};

socket.onmessage = (msg) => {
    console.log("data-events message", msg)
    const parsedData = JSON.parse(msg.data);
    updateDisplay(parsedData)
};


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