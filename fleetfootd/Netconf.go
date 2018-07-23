package main

import (
	"errors"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
)

type DHCPSettings struct {
	Enabled    bool   `json:"enabled"`
	RangeStart string `json:"range_start"`
	RangeEnd   string `json:"range_end"`
}

type Whitelist struct {
	URLs       *[]string `json:"urls,omitempty"`
	Ranges     *[]string `json:"ranges,omitempty"`
	AllowEmpty bool      `json:"allow_empty,omitempty"` // Can be used to allow empty whitelists, useful for URLs that might fail
}

type Interface struct {
	Name            string            `json:"name"`                       // Required
	IntelDriverBug  bool              `json:"intel_driver_bug,omitempty"` // Set to true on i219 devices
	Bridge          string            `json:"bridge,omitempty"`           // Omit if VLAN
	NetworkSettings map[string]string `json:"network_settings,omitempty"` // Can be omitted
	DHCP            *DHCPSettings     `json:"dhcp,omitempty"`             // Can be omitted
	MAC             string            `json:"mac,omitempty"`              // Omit if VLAN
	VLANS           *[]int16          `json:"vlans,omitempty"`            // Can be omitted
}

type PPPoEConfig struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	DefaultRoute bool   `json:"default_route,omitempty"`
}

type DHCPLease struct {
	Name string `json:"name"`
	MAC  string `json:"mac"`
	IP   string `json:"ip"`
}

type PortConf struct {
	Name         string `json:"name,omitempty"`
	TargetIP     string `json:"target_ip,omitempty"`
	TargetPort   uint16 `json:"target_port,omitempty"`
	TargetName   string `json:"target_name,omitempty"`
	Protocol     string `json:"protocol,omitempty"`
	ExternalPort uint16 `json:"external_port,omitempty"`
	Whitelist    string `json:"whitelist,omitempty"`
}

type VPNCConfig struct {
	TunnelInterface string    `json:"tunnel_interface,omitempty"`
	Remote          string    `json:"remote,omitempty"`
	PreSharedKey    string    `json:"pre_shared_key,omitempty"`
	Username        string    `json:"username,omitempty"`
	Password        string    `json:"password,omitempty"`
	DefaultRoute    bool      `json:"default_route,omitempty"`
	Routes          *[]string `json:"routes,omitempty"`
}

type Netconf struct {
	valid bool // Set internally after validation

	EnableMSSClamping bool                    `json:"enable_mss_clamping,omitempty"`
	Interfaces        map[string]*Interface   `json:"interfaces,omitempty"`
	PortForwarding    []*PortConf             `json:"port_forwarding,omitempty"`
	PPPoE             map[string]*PPPoEConfig `json:"pppoe,omitempty"`
	StaticDHCPLeases  map[string]*DHCPLease            `json:"static_dhcp_leases,omitempty"`
	Version           int                     `json:"version"`
	VLANS             map[int16]*Interface    `json:"vlans,omitempty"`
	VPNC              map[string]*VPNCConfig  `json:"vpnc,omitempty"`
	Whitelists        map[string]*Whitelist   `json:"whitelists,omitempty"`
}

func demo() {
	file, _ := os.Open("./netconf.yml")
	data, _ := ioutil.ReadAll(file)

	var nc Netconf
	yaml.Unmarshal(data, &nc)

	if !nc.IsValid() {
		log.Println("Errors have occurred during configuration parsing")
		os.Exit(1)
	}

	os.Exit(0)

	str, _ := yaml.Marshal(nc)
	log.Println(string(str))
	//
	//for k, v := range nc.Interfaces {
	//	log.Println(k, v)
	//	if v.DHCP != nil {
	//		log.Println(k, "DHCP:", *v.DHCP)
	//	}
	//	if v.VLANS != nil {
	//		log.Println(k, "VLANS:", *v.VLANS)
	//	}
	//	if v.NetworkSettings != nil {
	//		log.Println(k, "Network Settings:", v.NetworkSettings)
	//	}
	//}
	//
	//for k, v := range nc.VLANS {
	//	log.Println(k, v)
	//	if v.DHCP != nil {
	//		log.Println(k, "DHCP:", *v.DHCP)
	//	}
	//	if v.NetworkSettings != nil {
	//		log.Println(k, "Network Settings:", v.NetworkSettings)
	//	}
	//}
	//log.Println("MSS Clamping", nc.EnableMSSClamping)
	//
	//uplink0 := Interface{
	//	Name:           "uplink0",
	//	MAC:            "70:85:c2:02:c5:ef",
	//	IntelDriverBug: true,
	//	VLANS: &[]int16{
	//		8,
	//		10,
	//		601,
	//		666,
	//		3000,
	//		4000,
	//	},
	//	NetworkSettings: map[string]string{
	//		"Address": "10.0.0.1/28",
	//		"Gateway": "10.0.0.1",
	//	},
	//	DHCP: &DHCPSettings{
	//		Enabled:    true,
	//		RangeStart: "10.0.0.2",
	//		RangeEnd:   "10.0.0.14",
	//	},
	//}
	//
	//fiber0 := Interface{
	//	Name:   "fiber0",
	//	MAC:    "00:02:c9:55:11:76",
	//	Bridge: "servers",
	//}
	//
	//net := Netconf{
	//	Interfaces: map[string]*Interface{
	//		uplink0.Name: &uplink0,
	//		fiber0.Name:  &fiber0,
	//	},
	//	PPPoE:             nil,
	//	VLANS:             nil,
	//	StaticDHCPLeases:  nil,
	//	EnableMSSClamping: true,
	//}
	////log.Println(net)
	//str, _ = yaml.Marshal(net)
	////log.Println(string(str), err)
	//var nc2 Netconf
	//yaml.Unmarshal(str, &nc2)

	//log.Println(nc2)

}

var E_UNVALIDATED = errors.New("a method was called on an unvalidated object")
var E_INVALID_WHITELIST = errors.New("whitelist definition incomplete")
var E_INVALID_PORT_FORWARDING_NIL = errors.New("port forwarding is nil")
var E_INVALID_PORT_FORWARDING_NAME = errors.New("port forwarding name is illegal or empty")
var E_INVALID_PORT_FORWARDING_IP = errors.New("port forwarding target_ip is illegal or empty")
var E_INVALID_PORT_FORWARDING_PORT = errors.New("port forwarding target_port is illegal, empty or duplicate")
var E_INVALID_PORT_FORWARDING_PROTOCOL = errors.New("port forwarding protocol is not supported")
var E_INVALID_DHCP_MAC = errors.New("DHCP lease mac is illegal or empty")
var E_INVALID_PORT_FORWARDING_WHITELIST = errors.New("port forwarding whitelist is invalid")

func (netconf *Netconf) IsValid() bool {
	return netconf.Validate() == nil
}

func (netconf *Netconf) Validate() []*error {
	if netconf.valid {
		return nil
	}

	var validationErrors []*error

	// Matches all a-Z, 0-9 incl. - and _. It may not begin with 1 or 2 - (-, --)
	nameRegexp := regexp.MustCompile("^[^-]{1,2}[a-zA-Z0-9-_]+?$")
	macRegexp := regexp.MustCompile("^(?:[a-fA-F0-9]{2}:){5}[a-fA-F0-9]{2}$")
	println(macRegexp)
	allocatedTCPPorts := map[uint16]bool{}
	allocatedUDPPorts := map[uint16]bool{}

	for name, whitelist := range netconf.Whitelists {
		if whitelist == nil {
			delete(netconf.Whitelists, name)
			log.Printf("WARN: whitelist '%s' is nil. dropped it", name)
			continue
		}
		if whitelist.URLs == nil {
			whitelist.URLs = &[]string{}
		}
		if whitelist.Ranges == nil {
			whitelist.Ranges = &[]string{}
		}
		//TODO IP validation
		if len(*whitelist.URLs) == 0 && len(*whitelist.Ranges) == 0 && !whitelist.AllowEmpty {
			delete(netconf.Whitelists, name)
			log.Printf("ERR: whitelist '%s' empty. dropped it", name)
			validationErrors = append(validationErrors, &E_INVALID_WHITELIST)
			continue
		}
	}

	for name, lease := range netconf.StaticDHCPLeases {
		//TODO Messages
		if lease == nil {
			log.Printf("ERR: port forwarding '%s' is nil", name)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_NIL)
			continue
		}
		lease.Name = name
		if lease.Name == "" || !nameRegexp.MatchString(lease.Name) {
			log.Printf("ERR: port forwarding '%s' has an invalid name", lease.Name)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_NAME)
			continue
		}
		if lease.MAC == "" || !macRegexp.MatchString(lease.MAC) {
			log.Printf("ERR: port forwarding '%s' has an invalid mac '%s'", lease.Name, lease.MAC)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_NAME)
			continue
		}
		if lease.IP == "" || net.ParseIP(lease.IP) == nil {
			log.Printf("ERR: port forwarding '%s' invalid target_ip '%s'", lease.Name, lease.IP)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_IP)
		}

		//TODO Match IP against defined networks
	}

	for index, forwarding := range netconf.PortForwarding {
		if forwarding == nil {
			log.Printf("ERR: port forwarding at index %d is nil", index)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_NIL)
			continue
		}
		if forwarding.Name == "" || !nameRegexp.MatchString(forwarding.Name) {
			log.Printf("ERR: port forwarding '%s' (index %d) has an invalid name", forwarding.Name, index)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_NAME)
		}
		if forwarding.Protocol == "" {
			forwarding.Protocol = "tcp"
		}
		if forwarding.Whitelist != "" {
			if val, ok := netconf.Whitelists[forwarding.Whitelist]; !ok || val == nil {
				log.Printf("ERR: port forwarding '%s' (index %d) contains non-existing whitelist '%s'", forwarding.Name, index, forwarding.Whitelist)
				validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_WHITELIST)
			}
		}
		if forwarding.TargetPort == 0 {
			log.Printf("ERR: port forwarding '%s' (index %d) invalid target_port %d", forwarding.Name, index, forwarding.TargetPort)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_PORT)
		}
		if forwarding.ExternalPort == 0 {
			forwarding.ExternalPort = forwarding.TargetPort
		}

		switch forwarding.Protocol {
		case "tcp":
			if used, ok := allocatedTCPPorts[forwarding.ExternalPort]; ok && used {
				log.Printf("ERR: port forwarding '%s' (index %d) uses uccupied tcp port %d", forwarding.Name, index, forwarding.ExternalPort)
				validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_PORT)
			} else {
				allocatedTCPPorts[forwarding.ExternalPort] = true
			}
		case "udp":
			if used, ok := allocatedUDPPorts[forwarding.ExternalPort]; ok && used {
				log.Printf("ERR: port forwarding '%s' (index %d) uses uccupied udp port %d", forwarding.Name, index, forwarding.ExternalPort)
				validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_PORT)
			} else {
				allocatedUDPPorts[forwarding.ExternalPort] = true
			}
		default:
			log.Printf("ERR: port forwarding '%s' (index %d) contains unsupported protocol '%s'", forwarding.Name, index, forwarding.Protocol)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_PROTOCOL)
		}


		/// TODO Mash with the names and ips


		if forwarding.TargetIP == "" || net.ParseIP(forwarding.TargetIP) == nil {
			log.Printf("ERR: port forwarding '%s' (index %d) invalid target_ip '%s'", forwarding.Name, index, forwarding.TargetIP)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_IP)
			continue
		}
		_, ok := netconf.StaticDHCPLeases[forwarding.TargetName]
		if forwarding.TargetName != "" && !ok  {
			log.Printf("ERR: port forwarding '%s' (index %d) invalid target_name '%s'", forwarding.Name, index, forwarding.TargetIP)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_IP)
		}
		if forwarding.TargetIP != "" &&  forwarding.TargetName != "" {
			log.Printf("ERR: port forwarding '%s' (index %d) target_ip and target_host are set", forwarding.Name, index)
			validationErrors = append(validationErrors, &E_INVALID_PORT_FORWARDING_IP)
		}

	}

	return validationErrors
}
