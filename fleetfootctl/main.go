package main

import (
	"context"
	"net/http"
	"net"
	"flag"
	"os"
	"time"
)

var sockPath = flag.String("sockpath", "/var/run/fleetfootd.sock", "Socket to listen on")

type HookData struct {
	TTY          string `json:tty`     // Calling interface name
	PPPName      string `json:pppname` // ppp0 etc.
	ExternalIP   string `json:externalIP`
	RemotePeerIP string `json:remotePeerIP`
	Speed        int32  `json:speed`   // usually 0
	ipparam      string `json:ipparam` // arbitrary data
}

func main() {
	flag.Parse()

	http.DefaultClient = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", *sockPath)
			},
		},
	}

	if ran, exit := runHook(); ran == true {
		os.Exit(exit)
		return
	}
}
