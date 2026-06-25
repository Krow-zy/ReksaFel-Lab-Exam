package telemetry

import (
	"encoding/json"
	"time"
)

// ClientSession represents the state of a single monitored client in the network grid.
type ClientSession struct {
	ClientID   string    `json:"client_id"`
	Hostname   string    `json:"hostname"`
	IPAddress  string    `json:"ip_address"`
	MacAddress string    `json:"mac_address"`
	OS         string    `json:"os"`
	Uptime     int64     `json:"uptime_seconds"`
	JoinedAt   time.Time `json:"joined_at"`
	LastSeen   time.Time `json:"last_seen"`
	IsActive   bool      `json:"is_active"`
}

// ProcessSnapshot captures the system resources and running processes for anomaly detection.
type ProcessSnapshot struct {
	Timestamp  time.Time `json:"timestamp"`
	CPUUsage   float64   `json:"cpu_usage_percent"`
	MemoryUsed uint64    `json:"memory_used_bytes"`
	MemoryFree uint64    `json:"memory_free_bytes"`
	Processes  []string  `json:"processes"`
}

// AlertLog represents a security event logged by the client agent.
type AlertLog struct {
	EventID   string    `json:"event_id"`
	ClientID  string    `json:"client_id"`
	Timestamp time.Time `json:"timestamp"`
	Severity  string    `json:"severity"` // INFO, WARNING, CRITICAL
	Message   string    `json:"message"`
	Trigger   string    `json:"trigger_type"` // e.g., "forbidden_process", "screen_change"
}

// ToJSON serializes the ClientSession into a JSON byte slice.
func (cs *ClientSession) ToJSON() ([]byte, error) {
	return json.Marshal(cs)
}

// FromJSON deserializes a JSON byte slice into a ClientSession.
func FromJSON(data []byte) (*ClientSession, error) {
	var cs ClientSession
	if err := json.Unmarshal(data, &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}

// Model review: confirmed serialization tags for client payload structs
