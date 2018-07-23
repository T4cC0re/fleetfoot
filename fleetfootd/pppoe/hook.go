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
	TTY          string `json:tty`     // Calling interface name
	PPPName      string `json:pppname` // ppp0 etc.
	ExternalIP   string `json:externalIP`
	RemotePeerIP string `json:remotePeerIP`
	Speed        int32  `json:speed`   // usually 0
	ipparam      string `json:ipparam` // arbitrary data
}

const UP_IPV4 = 1
const UP_IPV6 = 2
const DOWN_IPV4 = 3
const DOWN_IPV6 = 4

var cfToken string
var cfMail string
var zoneIDs map[string]string
var DNSRecords map[string]string
var user cloudflare.User
var api *cloudflare.API
var initialized bool

func initialize() error {
	if initialized {
		return nil
	}
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

func Init(CloudflareMail string, CloudflareToken string) {
	cfToken = CloudflareToken
	cfMail = CloudflareMail
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

var up_ipv4 *HookData

func TryExec() {
	if up_ipv4 != nil {
		data := &up_ipv4
		zoneID, err := FetchZoneID("t4cc0.re")
		if err != nil {
			return
		}

		records := []string{"current.t4cc0.re"}
		for _, record := range records {
			recID, err := FetchDNSRecordID(zoneID, record)
			if err != nil {
				return
			}

			err = UpdateDNSRecord(zoneID, recID, (*data).ExternalIP, false, 120)
			if err != nil {
				return
			}
		}
	}
}

func Schedule(kind int, data *HookData) {
	if kind == UP_IPV4 {
		up_ipv4 = data
	}
}

func DataFromRequest(r *http.Request) (*HookData, error) {
	var data HookData

	if r.Body == nil {
		return nil, errors.New("no body")
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func HookUp(w http.ResponseWriter, r *http.Request) {
	if err := initialize(); err != nil {
		http.Error(w, "not ready", 503)
		return
	}
	data, err := DataFromRequest(r)
	if err != nil {
		throw(w, err)
	}

	Schedule(UP_IPV4, data)
	go w.WriteHeader(206)
	TryExec()
}

func HookDown(w http.ResponseWriter, r *http.Request) {
	if err := initialize(); err != nil {
		http.Error(w, "not ready", 503)
		return
	}
	w.WriteHeader(200)
}
func HookUp6(w http.ResponseWriter, r *http.Request) {
	if err := initialize(); err != nil {
		http.Error(w, "not ready", 503)
		return
	}
	w.WriteHeader(200)
}
func HookDown6(w http.ResponseWriter, r *http.Request) {
	if err := initialize(); err != nil {
		http.Error(w, "not ready", 503)
		return
	}
	w.WriteHeader(200)
}
