package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/koron/go-ssdp"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

type DisplayServerData struct {
	Type string
	Meta map[string]interface{}
}

type DisplayMetaSong struct {
	Section string
	Text    string
}

type DisplayServer struct {
	ctx    context.Context
	server *http.Server
	data   DisplayServerData
}

func NewDisplayServer() *DisplayServer {
	return &DisplayServer{}
}

func (d *DisplayServer) SetData(data DisplayServerData) {
	rt.LogDebugf(d.ctx, "DisplayServer.SetData: %s", data)
	d.data = data
}

func (d *DisplayServer) startup(ctx context.Context) {
	d.ctx = ctx

	go d.ssdpInit()
	d.server = &http.Server{Addr: ":7777"}

	d.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			This function handles all the "display side display logic"
			Its very basic and just a proof of concept at the moment.

			TODO:
			- Possibly figure out a way for this to be written in react
			- Figure out how we use events to update the display automatically when data changes.
			// https://thedevelopercafe.com/articles/server-sent-events-in-go-595ae2740c7a might be a possible place to start

		*/

		marshaledMeta, _ := json.Marshal(d.data.Meta)

		if d.data.Type == "song" {
			meta := &DisplayMetaSong{}
			err := json.Unmarshal(marshaledMeta, meta)
			if err != nil {
				rt.LogErrorf(d.ctx, "Error: %e", err)
			}

			text := meta.Text
			text = strings.ReplaceAll(text, "\n", "<br>")

			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf("<html><head><title>%s</title></head><body><h1>%s</h1>%s</body></html>", meta.Section, meta.Section, text)))
			return
		}

		// write data to the display
		data, _ := json.MarshalIndent(d.data, "", "  ")
		tempData := fmt.Sprintf("Display: %s", data)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(tempData))
	})

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
