package pppoe

import (
	"net/http"
	"github.com/prometheus/common/log"
	"github.com/cloudflare/cloudflare-go"
	"fmt"
	"encoding/json"
)

var cftoken string
var cfmail string

type HookData struct {
  TTY string `json:tty` // Calling interface name
  PPPName string `json:pppname` // ppp0 etc.
  ExternalIP string `json:externalIP`
  RemotePeerIP string `json:remotePeerIP`
  Speed int32 `json:speed` // usually 0
  ipparam string `json:ipparam` // arbitrary data
}

func Init(cfMail string, cfToken string) error {
	cfmail = cfMail
	cftoken = cfToken
	return nil
}

func Hook(w http.ResponseWriter, r *http.Request) {
	var data HookData

	if r.Body == nil {
		http.Error(w, "No body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Could not unmarshal json", 400)
		log.Errorf("Could not unmarshal json: %v", err)
		return
	}

	// Construct a new API object
	api, err := cloudflare.New(cftoken, cfmail)
	if err != nil {
		http.Error(w, "Could not instantiate cloudflare API", 500)
		log.Errorf("Could not instantiate cloudflare API: %v", err)
		return
	}

	// Fetch user details on the account
	u, err := api.UserDetails()
	if err != nil {
		http.Error(w, "Could not fetch user details", 500)
		log.Errorf("Could not fetch user details: %v", err)
		return
	}
	// Print user details
	fmt.Println(u)
	id, err := api.ZoneIDByName("t4cc0.re")
	if err != nil {
		http.Error(w, "Could not get zone", 500)
		log.Errorf("Could not get zone: %v", err)
		return
	}

	foo := cloudflare.DNSRecord{Name: "current.t4cc0.re"}
	recs, err := api.DNSRecords(id, foo)
	if err != nil {
		http.Error(w, "Could not get record", 500)
		log.Errorf("Could not get record: %v", err)
		return
	}

	if len(recs) < 1 {
		http.Error(w, "Could not get record id", 500)
		log.Errorf("Could not get record id: %v", err)
		return
	}
	rec := recs[:1][0]

	dnsrec := cloudflare.DNSRecord{Content: data.ExternalIP, Proxied: false, TTL: 120}

	err = api.UpdateDNSRecord(id, rec.ID, dnsrec)
	if err != nil {
		http.Error(w, "Could not update DNS", 500)
		log.Errorf("Could not update DNS: %v", err)
		return
	}

	w.WriteHeader(200)
}
