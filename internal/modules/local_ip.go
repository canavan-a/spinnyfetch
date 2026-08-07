package modules

import (
	"fmt"
	"net"
)

type localIPModule struct{}

func (localIPModule) Key() string { return "Local IP" }

func (localIPModule) Fetch() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.To4() == nil {
				continue
			}
			ones, _ := ipnet.Mask.Size()
			return fmt.Sprintf("%s: %s/%d", iface.Name, ipnet.IP.String(), ones), nil
		}
	}
	return "", fmt.Errorf("no active non-loopback IPv4 interface found")
}
