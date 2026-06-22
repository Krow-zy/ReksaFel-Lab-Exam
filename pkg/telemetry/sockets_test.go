package telemetry

import "testing"

func TestSocketScanner(t *testing.T) {
	// Forbidden external ports 443 (e.g. cloud AI proxies during exam)
	scanner := NewSocketScanner([]int{443, 8080})

	conns := []SocketConnection{
		{"127.0.0.1:8080", "127.0.0.1:51240", "ESTABLISHED", 4120},
		{"100.64.0.5:8081", "100.64.0.1:8080", "ESTABLISHED", 2980},
		{"192.168.1.15:443", "142.250.190.46:443", "ESTABLISHED", 1024},
	}

	// Verify detection of connection to forbidden port 443
	hasForbidden := scanner.HasForbiddenConnection(conns, 443)
	if !hasForbidden {
		t.Error("Expected scanner to find forbidden connection on port 443")
	}

	// Verify detection of connection to non-forbidden port 8081 (is in local list but not forbidden destination)
	hasForbidden = scanner.HasForbiddenConnection(conns, 8081)
	if hasForbidden {
		t.Error("Expected scanner to report false for port 8081 since it is not forbidden")
	}

	// Verify detection of connection to port 8080 (forbidden but local connection)
	hasForbidden = scanner.HasForbiddenConnection(conns, 8080)
	// In mock conns: "100.64.0.5:8081" connects to "100.64.0.1:8080" (remote port is 8080)
	if !hasForbidden {
		t.Error("Expected scanner to find forbidden connection on remote port 8080")
	}
}
