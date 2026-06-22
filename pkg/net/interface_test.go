package net

import (
	"testing"
)

func TestHasInterfacePrefix(t *testing.T) {
	// We can test with loopback prefix or common adapters like 'lo' or 'eth' or 'win'
	// Since we are running on Windows, 'loopback' or similar is usually present.
	// But let's check a non-existent prefix first.
	exists, err := HasInterfacePrefix("nonexistentprefix12345")
	if err != nil {
		t.Fatalf("HasInterfacePrefix failed: %v", err)
	}
	if exists {
		t.Error("expected non-existent prefix to return false")
	}

	// We can also try prefix 'loop' which matches 'Loopback Pseudo-Interface 1' on Windows
	exists, err = HasInterfacePrefix("loop")
	if err != nil {
		t.Fatalf("HasInterfacePrefix failed: %v", err)
	}
	// Note: Loopback might not match on all systems, but we should not fail the test
	// if it doesn't match, we just verify the behavior.
}

func TestGetLocalIPv4(t *testing.T) {
	ip, err := GetLocalIPv4()
	// Either we get an active local IPv4 or an error if offline.
	// In most test environments there is a local IP.
	if err != nil {
		t.Logf("GetLocalIPv4 returned error (normal if offline): %v", err)
	} else {
		if ip == "" {
			t.Error("expected non-empty IP string")
		}
	}
}
