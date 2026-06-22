package net

import (
	"fmt"
	"net"
	"time"
)

// MeasureTCPConnectTime measures the latency to establish a TCP connection.
// This is used for heartbeat ping latency monitoring in the mesh network.
func MeasureTCPConnectTime(address string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return 0, fmt.Errorf("ping fail: %w", err)
	}
	defer conn.Close()
	return time.Since(start), nil
}

// IsNodeReachable checks if a remote network node answers to TCP connection request.
func IsNodeReachable(address string, timeout time.Duration) bool {
	_, err := MeasureTCPConnectTime(address, timeout)
	return err == nil
}
