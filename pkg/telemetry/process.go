package telemetry

import "time"

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
