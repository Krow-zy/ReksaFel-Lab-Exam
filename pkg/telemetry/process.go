package telemetry

import (
	"strings"
	"time"
)

// ForbiddenProcess represents a blacklisted process running on a client machine.
type ForbiddenProcess struct {
	Name      string    `json:"name"`
	PID       int       `json:"pid"`
	StartTime time.Time `json:"start_time"`
}

// DefaultBlacklist returns a list of default forbidden process names.
// These typically represent tools that are prohibited during proctored exams.
func DefaultBlacklist() []string {
	return []string{
		"discord",
		"slack",
		"chrome",
		"firefox",
		"msedge",
		"teamviewer",
		"anydesk",
		"wireshark",
	}
}

// ScanProcesses matches a list of active process names against a blacklist.
// Case-insensitive comparisons are performed using the strings standard package.
func ScanProcesses(active []string, blacklist []string) []string {
	var found []string
	for _, act := range active {
		for _, black := range blacklist {
			if strings.EqualFold(act, black) {
				found = append(found, act)
				break
			}
		}
	}
	return found
}

// Process review: verified process blacklist filtering logic
