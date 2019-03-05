package main

import (
	"fmt"
	"gitlab.com/T4cC0re/fleetfoot/shared"
	"os"
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
)

func runHook() (bool, int) {
	invoked := os.Args[0]
	regex := regexp.MustCompile(`(?m)(ip(?:v6)?-(?:up|down))\.d`)
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
		url = "http://_/hook/up_ipv4"
	case "ip-up6":
		url = "http://_/hook/up_ipv6"
	case "ip-down":
		url = "http://_/hook/down_ipv4"
	case "ip-down6":
		url = "http://_/hook/down_ipv6"
	default:
		fmt.Println("unknown hook. Not running")
		return false, 0
	}
	fmt.Println(matches[1])

	speed, err := strconv.ParseInt(os.Args[3], 10, 32)

	u := shared.HookData{
		TTY:          os.Args[2],
		PPPName:      os.Args[1],
		ExternalIP:   os.Args[4],
		RemotePeerIP: os.Args[5],
		Speed:        int32(speed),
		Ipparam:      "",
	}

	fmt.Println(u)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	res, err := http.Post(url, "application/json; charset=utf-8", b)
	fmt.Println(res, err)

	return true, 0
}
