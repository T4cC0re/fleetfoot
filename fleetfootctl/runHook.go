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

func runHook() (bool, int) {
	fmt.Println(os.Args[0])

	invoked := os.Args[0]

	fmt.Println(len(os.Args))

	if match, _ := regexp.MatchString(`(?m)ip-up`, invoked); !match {
		return false, 0
	}

	if len(os.Args) < 6 {
		fmt.Errorf("invalid argument count")
		return false, 0
	}


	fmt.Println("Hook detected! Running...")

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
	res, err := http.Post("http://_/pppd/hook", "application/json; charset=utf-8", b)
	fmt.Println(res, err)

	return true, 0
}
