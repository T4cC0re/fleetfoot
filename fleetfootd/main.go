package main

import (
	"net/http"
	"fmt"
	"os"
	"net"
	"log"
	"flag"
	"./pppoe"
)

var sockPath = flag.String("sockpath","/var/run/fleetfootd.sock", "Socket to listen on")
var debugPort = flag.Int("debugport",0, "If set listens to HTTP requests on given port. DO NOT USE IN PROD!")
var cloudflareToken = flag.String("cloudflaretoken", "", "Token used for cloudflare API invokation")
var cloudflareMail = flag.String("cloudflaremail", "", "Email used for cloudflare API invokation")

func _404 (w http.ResponseWriter, r * http.Request) {
}

func main() {
	flag.Parse()
	os.Remove(*sockPath)
	unixListener, err := net.Listen("unix", *sockPath)
	if err != nil {
		log.Fatal("Listen (UNIX socket): ", err)
	}
	defer unixListener.Close()

	log.Println("Opened listening socket at", *sockPath)

	pppoe.Init(*cloudflareMail, *cloudflareToken)

	http.HandleFunc("/", _404)
	http.HandleFunc("/pppd/hook/up", pppoe.HookUp)
	http.HandleFunc("/pppd/hook/down", pppoe.HookDown)
	http.HandleFunc("/pppd/hook/up6", pppoe.HookUp6)
	http.HandleFunc("/pppd/hook/down6", pppoe.HookDown6)

	if *debugPort != 0 {
		go http.Serve(unixListener, nil)
		log.Println("Opening debug HTTP port", *debugPort)
		http.ListenAndServe(fmt.Sprintf(":%d", *debugPort), nil)
	}
	http.Serve(unixListener, nil)
}
