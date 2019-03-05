module gitlab.com/T4cC0re/fleetfoot

go 1.12

require (
	github.com/cloudflare/cloudflare-go v0.0.0-20190123000000-9837a599c0ba // manual: includes Spectrum API
	github.com/coreos/go-systemd v0.0.0-20190212144455-93d5ec2c7f76
	github.com/ghodss/yaml v1.0.0
	github.com/godbus/dbus v0.0.0-20181025153459-66d97aec3384 // indirect //manual: fixes issues with coreos/go-systemd
	github.com/okzk/sdnotify v0.0.0-20180710141335-d9becc38acbd
	github.com/pkg/errors v0.8.0
	github.com/prometheus/common v0.2.0
)
