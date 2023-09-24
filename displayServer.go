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

type DisplayServerData struct {
	Type string                 `json:"type"`
	Meta map[string]interface{} `json:"meta"`
}

type DisplayServer struct {
	ctx      context.Context
	server   *http.Server
	data     map[int]DisplayServerData
	dataChan chan map[int]DisplayServerData
}

func NewDisplayServer() *DisplayServer {
	return &DisplayServer{
		data:     make(map[int]DisplayServerData),
		dataChan: make(chan map[int]DisplayServerData),
	}
}

func (d *DisplayServer) SetData(data DisplayServerData) {
	rt.LogDebug(d.ctx, "DisplayServer.SetData called")
	d.data[0] = data
	d.dataChan <- map[int]DisplayServerData{0: data}
}

func (d *DisplayServer) startup(ctx context.Context) {
	d.ctx = ctx

	go d.ssdpInit()

	router := mux.NewRouter()

	var frontendFS = fs.FS(displayServerFrontend)
	staticContent, err := fs.Sub(frontendFS, "frontend-display/dist")
	if err != nil {
		rt.LogErrorf(d.ctx, "Error: %e", err)
	}

	httpfs := http.FileServer(http.FS(staticContent))
	router.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", httpfs))

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

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

	router.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		marshaledData, _ := json.Marshal(d.data[0])
		w.Write(marshaledData)
	}).Methods("GET")

	/*router.HandleFunc("/data-events", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			rt.LogError(d.ctx, "SSE not supported")
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// listens to change to the data channel, and writes the data to the response writer
		for dataMap := range d.dataChan {
			data := dataMap[0]
			rt.LogDebugf(d.ctx, "DisplayServer.dataChan: %s", data)

			m := map[string]any{"data": data}

			buff := bytes.NewBuffer([]byte{})
			encoder := json.NewEncoder(buff)
			err := encoder.Encode(m)
			if err != nil {
				fmt.Println(err)
				break
			}

			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("event: %s\n", "data-update"))
			sb.WriteString(fmt.Sprintf("data: %v\n", buff.String()))
			event := sb.String()

			if err != nil {
				fmt.Println(err)
				break
			}

			_, err = fmt.Fprint(w, event)
			if err != nil {
				fmt.Println(err)
				break
			}

			flusher.Flush()
		}
	})*/

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//marshaledMeta, _ := json.Marshal(d.data.Meta)

		// load index.html
		index, err := displayServerFrontend.ReadFile("frontend-display/index.html")
		if err != nil {
			rt.LogErrorf(d.ctx, "Error: %e", err)
		}

		w.Write(index)

	}).Methods("GET")

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

func (d *DisplayServer) shutdown(ctx context.Context) {
	d.server.Shutdown(ctx)
}

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

func (d *DisplayServer) getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		rt.LogFatalf(d.ctx, "Fatal: %e", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
