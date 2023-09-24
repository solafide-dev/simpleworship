package main

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/koron/go-ssdp"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend-display
var displayServerFrontend embed.FS

// DisplayServerData represents the data for a single slide
type DisplayServerData struct {
	Type string                 `json:"type"`
	Meta map[string]interface{} `json:"meta"`
}

// DisplayServer represents the display server
type DisplayServer struct {
	ctx      context.Context
	server   *http.Server
	data     map[int]DisplayServerData
	dataChan chan map[int]DisplayServerData
}

// NewDisplayServer creates a new DisplayServer
func NewDisplayServer() *DisplayServer {
	return &DisplayServer{
		data:     make(map[int]DisplayServerData),
		dataChan: make(chan map[int]DisplayServerData),
	}
}

// SetData sets the data for the display server
// TODO: Make this handle setting multiple slides
func (d *DisplayServer) SetData(data DisplayServerData) {
	rt.LogDebug(d.ctx, "DisplayServer.SetData called")
	d.data[0] = data
	d.dataChan <- map[int]DisplayServerData{0: data}
}

// startup is called at application startup and initializes the display server
func (d *DisplayServer) startup(ctx context.Context) {
	d.ctx = ctx

	go d.ssdpInit()

	// Define Router
	router := mux.NewRouter()

	// Setup static content servering for the display server frontend
	// these files are embedded in the binary
	var frontendFS = fs.FS(displayServerFrontend)
	staticContent, err := fs.Sub(frontendFS, "frontend-display/dist")
	if err != nil {
		rt.LogErrorf(d.ctx, "Error: %e", err)
	}

	httpfs := http.FileServer(http.FS(staticContent))
	router.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", httpfs))

	// Websocket client handling
	// Each new connection is added to the map,
	// and dataChan updates are sent to each client
	var clients = make(map[string]*websocket.Conn)

	// watch the dataBus for changes, and send the data to the client
	go func() {
		for dataMap := range d.dataChan {
			data := dataMap[0]
			rt.LogDebug(d.ctx, "DisplayServer.dataChan updated")

			marshaledData, _ := json.Marshal(data)

			for uuid, ws := range clients {
				rt.LogDebugf(d.ctx, "Sending data to client %s", uuid)
				err := ws.WriteMessage(websocket.TextMessage, marshaledData)
				if err != nil {
					rt.LogErrorf(d.ctx, "Error: %e", err)
					return
				}
			}
		}
	}()

	// Websocket Route, handles websocket connections
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			rt.LogErrorf(d.ctx, "Error: %e", err)
		}

		id := uuid.New().String()
		clients[id] = ws // add client to clients map

		rt.LogInfof(d.ctx, "New Websocket client connected: %s", id)

		// when a message is received, do something with it
		go func() {
			for {
				// TODO: Make this more robust, probably.
				_, msg, err := ws.ReadMessage()
				if err != nil {
					// websocket: close 1001 (going away): goodbye
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						rt.LogErrorf(d.ctx, "Error: %v", err)
					}
					rt.LogInfof(d.ctx, "Websocket client disconnected: %s", id)
					delete(clients, id) // yeet the client from the clients map
					return
				}
				rt.LogInfof(d.ctx, "ws message: %s", msg)
			}
		}()
	})

	// Data route, returns the current data as JSON
	router.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		marshaledData, _ := json.Marshal(d.data[0])
		w.Write(marshaledData)
	}).Methods("GET")

	// Root route, just servers index.html
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index, err := displayServerFrontend.ReadFile("frontend-display/index.html")
		if err != nil {
			rt.LogErrorf(d.ctx, "Error: %e", err)
		}

		w.Write(index)
	}).Methods("GET")

	// Define the server
	d.server = &http.Server{
		Addr:    ":7777",
		Handler: router,
	}

	// Start the server up
	go func() {
		if err := d.server.ListenAndServe(); err != nil {
			rt.LogErrorf(d.ctx, "Error: %e", err)
		}
	}()
}

// shutdown is called at application termination
func (d *DisplayServer) shutdown(ctx context.Context) {
	d.server.Shutdown(ctx)
}

// ssdpInit initializes the SSDP server
// TODO: Add this to the display server struct and handle shutdown properly
func (d *DisplayServer) ssdpInit() {
	rt.LogInfo(d.ctx, "ssdp init")

	ip := d.getOutboundIP()
	id := uuid.New().String()

	ad, err := ssdp.Advertise(
		"simpleworship:main",           // send as "ST"
		"uuid:"+id,                     // send as "USN"
		"http://"+ip.String()+":7777/", // send as "LOCATION"
		"simpleworship",                // send as "SERVER"
		1800)                           // send as "maxAge" in "CACHE-CONTROL"
	if err != nil {
		panic(err)
	}

	// to detect CTRL-C is pressed.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	var aliveTick <-chan time.Time

loop:
	for {
		select {
		case <-aliveTick:
			rt.LogInfo(d.ctx, "ssdp alive tick")
			ad.Alive()
		case <-quit:
			break loop
		}
	}

	// send/multicast "byebye" message.
	ad.Bye()
	// teminate Advertiser.
	ad.Close()
}

// getOutboundIP gets the outbound IP address (local IP address) to use with SSDP messages
func (d *DisplayServer) getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		rt.LogFatalf(d.ctx, "Fatal: %e", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
