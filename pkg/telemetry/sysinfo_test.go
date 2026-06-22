package telemetry

import (
	"bytes"
	"testing"
)

func TestSystemResource_ToJSON(t *testing.T) {
	sr := &SystemResource{
		CPUUsage:    25.5,
		TotalMemory: 16000,
		FreeMemory:  8000,
		Timestamp:   123456789,
	}

	data, err := sr.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if !bytes.Contains(data, []byte("cpu_usage")) {
		t.Error("JSON does not contain cpu_usage")
	}
}

func TestSystemResource_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		sr    SystemResource
		valid bool
	}{
		{"valid metrics", SystemResource{50.0, 16000, 8000, 123}, true},
		{"cpu negative", SystemResource{-1.0, 16000, 8000, 123}, false},
		{"cpu overflow", SystemResource{101.0, 16000, 8000, 123}, false},
		{"zero total memory", SystemResource{50.0, 0, 0, 123}, false},
		{"free memory exceeds total", SystemResource{50.0, 16000, 17000, 123}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sr.IsValid() != tt.valid {
				t.Errorf("expected validity %v, got %v", tt.valid, tt.sr.IsValid())
			}
		})
	}
}

func TestGetSystemInfo(t *testing.T) {
	info := GetSystemInfo()
	if info == nil {
		t.Fatal("GetSystemInfo returned nil")
	}
	if !info.IsValid() {
		t.Errorf("GetSystemInfo returned invalid metrics: %+v", info)
	}
}
