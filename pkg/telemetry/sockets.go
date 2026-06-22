package telemetry

// SocketConnection represents an active TCP socket connection.
type SocketConnection struct {
	LocalAddr  string `json:"local_address"`
	RemoteAddr string `json:"remote_address"`
	State      string `json:"state"`
	PID        int    `json:"pid"`
}

// SocketScanner lists active TCP sockets on the host.
type SocketScanner struct {
	ForbiddenPorts []int
}

// NewSocketScanner configures a scanner to inspect active sockets.
func NewSocketScanner(forbiddenPorts []int) *SocketScanner {
	return &SocketScanner{ForbiddenPorts: forbiddenPorts}
}

// ScanMockConnections returns mock connections to simulate socket scanning.
func (s *SocketScanner) ScanMockConnections() []SocketConnection {
	return []SocketConnection{
		{"127.0.0.1:8080", "127.0.0.1:51240", "ESTABLISHED", 4120},
		{"100.64.0.5:8081", "100.64.0.1:8080", "ESTABLISHED", 2980},
		{"192.168.1.15:443", "142.250.190.46:443", "ESTABLISHED", 1024}, // External web connection
	}
}

// HasForbiddenConnection checks if any connections link to forbidden destination ports.
func (s *SocketScanner) HasForbiddenConnection(conns []SocketConnection, targetPort int) bool {
	for _, p := range s.ForbiddenPorts {
		if p == targetPort {
			for _, conn := range conns {
				// Simple check if remote port matches forbidden port
				if s.isPortMatch(conn.RemoteAddr, targetPort) {
					return true
				}
			}
		}
	}
	return false
}

func (s *SocketScanner) isPortMatch(addr string, port int) bool {
	// Address is usually ip:port. We check if address suffix matches :port
	suffix := ":"
	switch port {
	case 80: suffix = ":80"
	case 443: suffix = ":443"
	case 8080: suffix = ":8080"
	case 8081: suffix = ":8081"
	}
	return len(addr) >= len(suffix) && addr[len(addr)-len(suffix):] == suffix
}
