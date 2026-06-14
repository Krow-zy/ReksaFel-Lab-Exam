package net

import (
	"context"
	"net"
	"sync"
	"time"
)

// ProbeResult holds the connection outcome of a single node check.
type ProbeResult struct {
	IPAddress string        `json:"ip_address"`
	Port      string        `json:"port"`
	IsAlive   bool          `json:"is_alive"`
	Latency   time.Duration `json:"latency"`
	Error     error         `json:"error,omitempty"`
}

// ProbeNode performs a TCP dial check to test connectivity and measure latency.
func ProbeNode(ctx context.Context, ip string, port string, timeout time.Duration) ProbeResult {
	start := time.Now()
	dialer := net.Dialer{Timeout: timeout}

	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(ip, port))
	duration := time.Since(start)

	if err != nil {
		return ProbeResult{
			IPAddress: ip,
			Port:      port,
			IsAlive:   false,
			Latency:   0,
			Error:     err,
		}
	}
	defer conn.Close()

	return ProbeResult{
		IPAddress: ip,
		Port:      port,
		IsAlive:   true,
		Latency:   duration,
	}
}

// ConcurrentProbe runs network checks on multiple target IPs in parallel.
func ConcurrentProbe(ctx context.Context, ips []string, port string, timeout time.Duration) []ProbeResult {
	results := make([]ProbeResult, 0, len(ips))
	resultChan := make(chan ProbeResult, len(ips))

	var wg sync.WaitGroup

	for _, ip := range ips {
		wg.Add(1)
		go func(targetIP string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				resultChan <- ProbeResult{
					IPAddress: targetIP,
					Port:      port,
					IsAlive:   false,
					Error:     ctx.Err(),
				}
			default:
				res := ProbeNode(ctx, targetIP, port, timeout)
				resultChan <- res
			}
		}(ip)
	}

	// Wait in a separate goroutine to prevent blocking the channel read
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect all results
	for res := range resultChan {
		results = append(results, res)
	}

	return results
}
