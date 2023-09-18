package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/koron/go-ssdp"
)

/*
	This is very _very_ temporary and just for proof of concept.
	Ideally we expose this in some way that actually interacts with the rest of the app.

	This is mainly to test the SSDP functionality and make sure it works.
*/

func startDisplayServer() {

	go ssdpInit()

	fmt.Println("Starting display server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Worship Display")
	})

	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func ssdpInit() {
	fmt.Println("ssdp init")

	ip := GetOutboundIP()
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
			fmt.Println("ssdp alive tick")
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

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
