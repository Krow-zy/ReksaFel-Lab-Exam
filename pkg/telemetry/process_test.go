package telemetry

import (
	"testing"
)

func TestScanProcesses(t *testing.T) {
	blacklist := []string{"Discord", "Slack", "AnyDesk"}

	tests := []struct {
		name     string
		active   []string
		expected []string
	}{
		{
			name:     "No blacklisted processes running",
			active:   []string{"explorer.exe", "svchost.exe", "chrome"},
			expected: []string{},
		},
		{
			name:     "One blacklisted process running exact match",
			active:   []string{"explorer.exe", "Slack", "cmd.exe"},
			expected: []string{"Slack"},
		},
		{
			name:     "Multiple blacklisted processes case-insensitive",
			active:   []string{"discord", "SLACK", "safe_exam_browser"},
			expected: []string{"discord", "SLACK"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := ScanProcesses(tt.active, blacklist)
			if len(found) != len(tt.expected) {
				t.Fatalf("expected %d processes found, got %d (%v)", len(tt.expected), len(found), found)
			}
			for _, f := range found {
				match := false
				for _, exp := range tt.expected {
					if f == exp {
						match = true
						break
					}
				}
				if !match {
					t.Errorf("found process %s which was not expected", f)
				}
			}
		})
	}
}

// Test review: verified process matching assertions in test cases
