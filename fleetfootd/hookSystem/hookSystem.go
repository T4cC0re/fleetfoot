package hookSystem

import (
	"github.com/pkg/errors"
	"net/http"
	"io/ioutil"
	"github.com/prometheus/common/log"
	"regexp"
)

var hooks = map[string]Hook{}

type Hook func(payload interface{}) (interface{}, error)

var EExists = errors.New("Hook already registered")
var ENoHook = errors.New("Hook does not exist")
var EInvalidPayload = errors.New("Payload invalid")

func AddHook(name string, hook Hook) (error) {
	if _, ok := hooks[name]; ok {
		return EExists
	}
	hooks[name] = hook
	return nil
}

func Run(name string, payload interface{}) (interface{}, error) {
	if hook, ok := hooks[name]; ok {
		return hook(payload)
	}
	return nil, ENoHook
}

func Hooks() map[string]Hook {
	return hooks
}

func DataFromRequest(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return []byte{}, nil
	}

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	return data, err
}

func throw(w http.ResponseWriter, err error) {
	http.Error(w, "Error during execution", 500)
	log.Errorf("Error during execution: %v", err)
}

func RunWWW(w http.ResponseWriter, r *http.Request) {
	data, err := DataFromRequest(r)
	if err != nil {
		throw(w, err)
	}

	var re = regexp.MustCompile(`(?i)^/hook/([^?]+).*$`)
	rs := re.FindStringSubmatch(r.RequestURI)
	if len(rs) < 2 {
		throw(w, err)
		return
	}

	hook := rs[1]
	log.Infof("Running hook '%s'...", hook)
	ret, err := Run(rs[1], data)

	switch err {
	case ENoHook:
		log.Infof("hook '%s' not found", hook)
		w.WriteHeader(404)
	case EInvalidPayload:
		log.Infof("hook '%s' reported an invalid payload", hook)
		w.WriteHeader(400)
	case nil:
		log.Infof("hook '%s' executed OK", hook)
		w.WriteHeader(200)
		if rData, ok := ret.([]byte); ok {
			w.Write(rData)
		}
	default:
		throw(w, err)
	}
}
