package pppoe

import (
	"github.com/prometheus/common/log"
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	"fmt"
	"gitlab.com/T4cC0re/fleetfoot/daemon/hookSystem"
	"gitlab.com/T4cC0re/fleetfoot/shared"
	"encoding/json"
)

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

	hookSystem.AddHook("up_ipv4", HookUp)
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

func HookUp(data interface{}) (interface{}, error) {
	if err := initialize(); err != nil {
		return nil, err
	}

	var raw []byte
	var ok bool
	if raw, ok = data.([]byte); !ok {
		return nil, hookSystem.EInvalidPayload
	}

	var hookData shared.HookData
	err := json.Unmarshal(raw, &hookData)

	if err != nil {
		return nil, err
	}

	zoneID, err := FetchZoneID("t4cc0.re")
	if err != nil {
		return nil, err
	}

	apps, err := api.SpectrumApplications(zoneID)
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		if app.DNS.Name == "ssh.t4cc0.re" {
			log.Infoln(app.OriginDirect)
			direct := fmt.Sprintf("tcp://%s:1337", hookData.ExternalIP)
			AppId := app.ID
			app.OriginDirect = []string{direct}
			app.ID = ""
			log.Infoln(app.OriginDirect)
			_, err := api.UpdateSpectrumApplication(zoneID, AppId, app)
			if err != nil {
				log.Errorln(err)
			}
		}
	}

	records := []string{"current.t4cc0.re"}
	for _, record := range records {
		recID, err := FetchDNSRecordID(zoneID, record)
		if err != nil {
			return nil, err
		}

		err = UpdateDNSRecord(zoneID, recID, hookData.ExternalIP, false, 120)
		if err != nil {
			return nil, err
		}
	}

	return true, nil
}
