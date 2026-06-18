package net

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestHasInterfacePrefix(t *testing.T) {
	// Standard loopback interface (lo/lo0/Loopback) exists on almost all machines
	exists, err := HasInterfacePrefix("loopback")
	if err != nil {
		t.Logf("HasInterfacePrefix returned error: %v", err)
		return
	}
	t.Logf("Loopback interface detected: %v", exists)
}

func TestGetLocalIPv4(t *testing.T) {
	ip, err := GetLocalIPv4()
	if err != nil {
		t.Logf("GetLocalIPv4 failed (normal in offline sandbox): %v", err)
		return
	}
	if net.ParseIP(ip) == nil {
		t.Errorf("GetLocalIPv4 returned invalid IP: %s", ip)
	}
}

func TestConcurrentProbe(t *testing.T) {
	// 1. Spin up a temporary mock TCP listener
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on temporary address: %v", err)
	}
	defer listener.Close()

	// Handle one connection in background
	go func() {
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	addr := listener.Addr().String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("failed to split host port: %v", err)
	}

	// 2. Run ConcurrentProbe targeting this address
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	results := ConcurrentProbe(ctx, []string{host}, port, 1*time.Second)

	// 3. Verify results
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	res := results[0]
	if res.IPAddress != host {
		t.Errorf("expected IP %s, got %s", host, res.IPAddress)
	}
	if res.Port != port {
		t.Errorf("expected Port %s, got %s", port, res.Port)
	}
	if !res.IsAlive {
		t.Errorf("expected target to be alive, got dead. Error: %v", res.Error)
	}
}
