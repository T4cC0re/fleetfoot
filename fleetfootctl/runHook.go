package main

import (
	"fmt"
	"os"
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
)

type HookData struct {
	TTY          string `json:tty`     // Calling interface name
	PPPName      string `json:pppname` // ppp0 etc.
	ExternalIP   string `json:externalIP`
	RemotePeerIP string `json:remotePeerIP`
	Speed        int32  `json:speed`   // usually 0
	ipparam      string `json:ipparam` // arbitrary data
}

func runHook() (bool, int) {
	invoked := os.Args[0]
	regex := regexp.MustCompile(`(?m)(ip(?:v6)?-(?:up|down))\.d/`)
	matches := regex.FindStringSubmatch(invoked)
	if matches == nil || len(matches) < 2 {
		return false, 0
	}

	if len(os.Args) < 6 {
		fmt.Errorf("invalid argument count")
		return false, 0
	}
	fmt.Print("Hook detected! Running ")

	var url string
	switch matches[1] {
	case "ip-up":
		url = "http://_/pppd/hook/up"
	case "ip-up6":
		url = "http://_/pppd/hook/up6"
	case "ip-down":
		url = "http://_/pppd/hook/down"
	case "ip-down6":
		url = "http://_/pppd/hook/down6"
	default:
		fmt.Println("unknown hook. Not running")
		return false, 0
	}
	fmt.Println(matches[1])

	speed, err := strconv.ParseInt(os.Args[3], 10, 32)

	u := HookData{
		TTY:          os.Args[2],
		PPPName:      os.Args[1],
		ExternalIP:   os.Args[4],
		RemotePeerIP: os.Args[5],
		Speed:        int32(speed),
		ipparam:      "",
	}

	fmt.Println(u)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	res, err := http.Post(url, "application/json; charset=utf-8", b)
	fmt.Println(res, err)

	return true, 0
}
