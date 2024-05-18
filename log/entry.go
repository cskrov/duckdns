package log

import "time"

type LogEntry struct {
	Level     string    `json:"level"`
	Timestamp Timestamp `json:"timestamp"`
	Message   string    `json:"message"`
	Updated   bool      `json:"updated"`
	IPv4      string    `json:"ipv4,omitempty"`
	IPv6      string    `json:"ipv6,omitempty"`
}

func createLogEntry(level string, message string) *LogEntry {
	return &LogEntry{
		Timestamp: Timestamp(time.Now().UTC().Unix()),
		Level:     level,
		Message:   message,
	}
}
