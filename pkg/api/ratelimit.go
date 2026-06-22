package api

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimiter manages client request frequency limits.
type RateLimiter struct {
	mu           sync.Mutex
	clients      map[string][]time.Time
	limit        int
	window       time.Duration
}

// NewRateLimiter creates a RateLimiter configured for a window duration and limit.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

// Limit returns a middleware that limits request rates.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)

		// Prune old request times
		var active []time.Time
		for _, t := range rl.clients[ip] {
			if t.After(cutoff) {
				active = append(active, t)
			}
		}

		if len(active) >= rl.limit {
			rl.mu.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		rl.clients[ip] = append(active, now)
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
