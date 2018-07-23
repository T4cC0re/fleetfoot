package main

import (
	"./pppoe"
	"flag"
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"github.com/okzk/sdnotify"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var sockPath = flag.String("sockpath", "/var/run/fleetfootd.sock", "Socket to listen on")
var debugPort = flag.Int("debugport", 0, "If set listens to HTTP requests on given port. DO NOT USE IN PROD!")
var cloudflareToken = flag.String("cloudflaretoken", "", "Token used for cloudflare API invokation")
var cloudflareMail = flag.String("cloudflaremail", "", "Email used for cloudflare API invokation")

var hookedInterfaces = map[string]string{}
var hookedInterfaces6 = map[string]string{}
var dbusConn *dbus.Conn

func _404(w http.ResponseWriter, r *http.Request) {
}

func reload(notify bool) {
	if notify {
		sdnotify.Reloading()
	}

	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPAddr:
				ip = v.IP
			case *net.IPNet:
				ip = v.IP
			}
			// process IP address
			log.Println(i, ip)
			if _ip, ok := hookedInterfaces[i.Name]; ok && ip != nil && ip.To4() != nil {
				log.Println("Updating", i.Name, _ip, ip.String())
				hookedInterfaces[i.Name] = string(ip.String())
			}
			if _ip, ok := hookedInterfaces6[i.Name]; ok && ip != nil && ip.To4() == nil {
				log.Println("Updating6", i.Name, _ip, ip.String())
				hookedInterfaces6[i.Name] = string(ip.String())
			}
		}
	}

	if notify {
		sdnotify.Ready()
	}
}
func main() {
	flag.Parse()
	os.Remove(*sockPath)
	unixListener, err := net.Listen("unix", *sockPath)
	if err != nil {
		log.Fatal("Listen (UNIX socket): ", err)
	}
	defer unixListener.Close()

	log.Println("Opened listening socket at", *sockPath)

	pppoe.Init(*cloudflareMail, *cloudflareToken)
	dbusConn, err = dbus.New()
	info, err := dbusConn.GetUnitProperties("not-a-docker.service")
	log.Println(info["LoadState"], err)

	demo()

	http.HandleFunc("/", _404)
	http.HandleFunc("/pppd/hook/up", pppoe.HookUp)
	http.HandleFunc("/pppd/hook/down", pppoe.HookDown)
	http.HandleFunc("/pppd/hook/up6", pppoe.HookUp6)
	http.HandleFunc("/pppd/hook/down6", pppoe.HookDown6)

	reload(false)

	if *debugPort != 0 {
		log.Println("Opening debug HTTP port", *debugPort)
		go http.ListenAndServe(fmt.Sprintf(":%d", *debugPort), nil)
	}
	go http.Serve(unixListener, nil)
	// Talks to sd_notify, sets up reloading and watchdog
	enterSystemdLifecycle()
}

func enterSystemdLifecycle() {
	log.Println("entering systemd/sd_notify compatible lifecycle...")
	sdnotify.Ready()
	go func() {
		tick := time.Tick(30 * time.Second)
		for {
			<-tick
			log.Println("watchdog reporting")
			sdnotify.Watchdog()
		}
	}()
	go func() {
		tick := time.Tick(5 * time.Minute)
		for {
			<-tick
			reload(true)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Println("enterred systemd/sd_notify compatible lifecycle...")
	for sig := range sigCh {
		if sig == syscall.SIGHUP {
			reload(true)
		} else {
			break
		}
	}

	sdnotify.Stopping()
	log.Println("exiting...")
}
