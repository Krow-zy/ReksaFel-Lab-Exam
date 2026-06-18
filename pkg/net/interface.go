package net

import (
	"errors"
	"net"
	"strings"
)

// HasInterfacePrefix checks if there is an active interface matching the prefix.
// Useful to verify if Tailscale or ReksaFel VPN interfaces are running locally.
func HasInterfacePrefix(prefix string) (bool, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}

	for _, iface := range interfaces {
		if strings.HasPrefix(strings.ToLower(iface.Name), strings.ToLower(prefix)) {
			return true, nil
		}
	}
	return false, nil
}

// GetLocalIPv4 returns the first non-loopback local IPv4 address.
func GetLocalIPv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("no active IPv4 address found")
}
