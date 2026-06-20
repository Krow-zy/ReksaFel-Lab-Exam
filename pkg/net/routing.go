package net

import (
	"net"
)

// IsIPInSubnet checks if the given IP address string belongs to the target CIDR block.
// Standard library net package is used to remain lightweight and dependency-free.
func IsIPInSubnet(ipStr, cidrStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, nil
	}

	_, subnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false, err
	}

	return subnet.Contains(ip), nil
}
