package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koron/go-ssdp"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend-display
var displayServerFrontend embed.FS

type DisplayServerData struct {
	Type string
	Meta map[string]interface{}
}

type DisplayMetaSong struct {
	Section string
	Text    string
}

type DisplayServer struct {
	ctx      context.Context
	server   *http.Server
	data     DisplayServerData
	dataChan chan DisplayServerData
}

func NewDisplayServer() *DisplayServer {
	return &DisplayServer{
		dataChan: make(chan DisplayServerData),
	}
}

func (d *DisplayServer) SetData(data DisplayServerData) {
	rt.LogDebugf(d.ctx, "DisplayServer.SetData: %s", data)
	d.data = data
	d.dataChan <- data
}

func (d *DisplayServer) startup(ctx context.Context) {
	d.ctx = ctx

	go d.ssdpInit()

	router := mux.NewRouter()

	var frontendFS = fs.FS(displayServerFrontend)
	staticContent, err := fs.Sub(frontendFS, "frontend-display/src")
	if err != nil {
		rt.LogErrorf(d.ctx, "Error: %e", err)
	}

	httpfs := http.FileServer(http.FS(staticContent))
	router.PathPrefix("/src/").Handler(http.StripPrefix("/src/", httpfs))

	router.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		marshaledData, _ := json.Marshal(d.data)
		w.Write(marshaledData)
	}).Methods("GET")

	router.HandleFunc("/data-events", func(w http.ResponseWriter, r *http.Request) {
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
		for data := range d.dataChan {

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
			sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))

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
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//marshaledMeta, _ := json.Marshal(d.data.Meta)

		// load index.html
		index, err := displayServerFrontend.ReadFile("frontend-display/index.html")
		if err != nil {
			rt.LogErrorf(d.ctx, "Error: %e", err)
		}

		_, err = displayServerFrontend.ReadFile("frontend-display/src/app.js")
		if err != nil {
			rt.LogErrorf(d.ctx, "Error: %e", err)
			w.Write([]byte(err.Error()))
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
