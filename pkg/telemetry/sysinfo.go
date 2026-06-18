package telemetry

import "time"

// SystemResource represents client machine hardware metrics.
type SystemResource struct {
	CPUUsage     float64 `json:"cpu_usage"`
	TotalMemory  uint64  `json:"total_memory_bytes"`
	FreeMemory   uint64  `json:"free_memory_bytes"`
	Timestamp    int64   `json:"timestamp"`
}

// GetSystemInfo returns a mock of current hardware metrics.
// Standard library package runtime is used to keep it minimal and dependency-free.
func GetSystemInfo() *SystemResource {
	return &SystemResource{
		CPUUsage:    12.5, // Mock CPU usage percentage
		TotalMemory: 16 * 1024 * 1024 * 1024, // Mock 16 GB
		FreeMemory:  8 * 1024 * 1024 * 1024,  // Mock 8 GB
		Timestamp:   time.Now().Unix(),
	}
}
