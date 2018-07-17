package pppoe

import (
	"net/http"
	"github.com/prometheus/common/log"
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	"fmt"
	"encoding/json"
)

type HookData struct {
  TTY string `json:tty` // Calling interface name
  PPPName string `json:pppname` // ppp0 etc.
  ExternalIP string `json:externalIP`
  RemotePeerIP string `json:remotePeerIP`
  Speed int32 `json:speed` // usually 0
  ipparam string `json:ipparam` // arbitrary data
}

var zoneIDs map[string]string
var DNSRecords map[string]string
var user cloudflare.User
var api *cloudflare.API

func Init(cfMail string, cfToken string) error {
	zoneIDs = map[string]string{}
	DNSRecords = map[string]string{}
	var err error
	api, err = cloudflare.New(cfToken, cfMail)
	if err != nil {
		return err
	}
	user, err = api.UserDetails()
	if err != nil {
		return err
	}
	fmt.Println(user)
	return nil
}

func FetchZoneID(zoneName string) (string, error) {
	if val, ok := zoneIDs[zoneName]; ok {
		return val, nil
	}
	id, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return "", err
	}
	zoneIDs[zoneName] = id
	return id, nil
}
func FetchDNSRecordID(zoneID string, recordName string) (string, error) {
	if val, ok := DNSRecords[recordName]; ok {
		return val, nil
	}
	tmp := cloudflare.DNSRecord{Name: recordName}
	recs, err := api.DNSRecords(zoneID, tmp)
	if err != nil {
		return "", err
	}

	if len(recs) < 1 {
		return "", errors.New("no records found")
	}
	if len(recs) > 1 {
		return "", errors.New("not a unique record")
	}

	rec := recs[:1][0].ID
	DNSRecords[recordName] = rec
	return rec, nil
}

func UpdateDNSRecord(zoneID string, recordID string, content string, proxied bool, TTL int) error {
	dnsrec := cloudflare.DNSRecord{Content: content, Proxied: proxied, TTL: TTL}
	err := api.UpdateDNSRecord(zoneID, recordID, dnsrec)
	if err != nil {
		return err
	}

	return nil
}

func throw(w http.ResponseWriter, err error) {
	http.Error(w, "Error during execution", 500)
	log.Errorf("Error during execution: %v", err)
}

func HookUp (w http.ResponseWriter, r *http.Request) {
	var data HookData

	if r.Body == nil {
		http.Error(w, "No body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		throw(w, err)
		return
	}

	zoneID, err := FetchZoneID("t4cc0.re")
	if err != nil {
		throw(w, err)
		return
	}

	records := []string{"current.t4cc0.re"}
	for _, record := range records {
		recID, err := FetchDNSRecordID(zoneID, record)
		if err != nil {
			throw(w, err)
			return
		}

		err = UpdateDNSRecord(zoneID, recID, data.ExternalIP, false, 120)
		if err != nil {
			throw(w, err)
			return
		}
	}

	w.WriteHeader(200)
}
func HookDown (w http.ResponseWriter, r *http.Request) {}
func HookUp6 (w http.ResponseWriter, r *http.Request) {}
func HookDown6 (w http.ResponseWriter, r *http.Request) {}
