package shared

type HookData struct {
	TTY          string `json:tty`     // Calling interface name
	PPPName      string `json:pppname` // ppp0 etc.
	ExternalIP   string `json:externalIP`
	RemotePeerIP string `json:remotePeerIP`
	Speed        int32  `json:speed,omitempty`   // usually 0
	Ipparam      string `json:ipparam,omitempty` // arbitrary data, optional
}
