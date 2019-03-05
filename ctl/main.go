package main

import (
	"context"
	"net/http"
	"net"
	"flag"
	"os"
	"time"
)

var sockPath = flag.String("sockpath", "/var/run/fleetfootd.sock", "Socket to connect to")

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
