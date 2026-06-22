package net

import (
	"net"
	"testing"
	"time"
)

func TestMeasureTCPConnectTime(t *testing.T) {
	// Start mock TCP listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start mock listener: %v", err)
	}
	defer l.Close()

	// Start accepting connections in background
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	addr := l.Addr().String()

	// Measure connect time
	latency, err := MeasureTCPConnectTime(addr, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("MeasureTCPConnectTime failed: %v", err)
	}

	if latency <= 0 {
		t.Errorf("Expected positive latency duration, got %v", latency)
	}

	// Verify reachability
	if !IsNodeReachable(addr, 500*time.Millisecond) {
		t.Error("Expected address to be reachable")
	}

	// Verify unreachable address
	unreachableAddr := "127.0.0.1:9999" // Assuming this port is unused
	if IsNodeReachable(unreachableAddr, 100*time.Millisecond) {
		t.Error("Expected unreachable address to return false")
	}
}
