package telemetry

import (
	"encoding/json"
	"testing"
)

func TestSystemResourceValidation(t *testing.T) {
	tests := []struct {
		name  string
		sys   SystemResource
		valid bool
	}{
		{
			name: "Valid Resource Data",
			sys: SystemResource{
				CPUUsage:    45.2,
				TotalMemory: 8000,
				FreeMemory:  4000,
			},
			valid: true,
		},
		{
			name: "Invalid CPU Usage High",
			sys: SystemResource{
				CPUUsage:    101.0,
				TotalMemory: 8000,
				FreeMemory:  4000,
			},
			valid: false,
		},
		{
			name: "Invalid CPU Usage Negative",
			sys: SystemResource{
				CPUUsage:    -1.5,
				TotalMemory: 8000,
				FreeMemory:  4000,
			},
			valid: false,
		},
		{
			name: "Free Memory Exceeds Total",
			sys: SystemResource{
				CPUUsage:    50.0,
				TotalMemory: 8000,
				FreeMemory:  9000,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sys.IsValid(); got != tt.valid {
				t.Errorf("SystemResource.IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestSystemResourceSerialization(t *testing.T) {
	sys := &SystemResource{
		CPUUsage:    25.0,
		TotalMemory: 1000,
		FreeMemory:  500,
		Timestamp:   123456789,
	}

	data, err := sys.ToJSON()
	if err != nil {
		t.Fatalf("failed to serialize: %v", err)
	}

	var parsed SystemResource
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to deserialize: %v", err)
	}

	if parsed.CPUUsage != sys.CPUUsage || parsed.TotalMemory != sys.TotalMemory || parsed.FreeMemory != sys.FreeMemory {
		t.Errorf("deserialized data does not match original struct")
	}
}
