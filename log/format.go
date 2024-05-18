package log

import (
	"encoding/json"
	"fmt"
	"os"
)

func formatLog(entry *LogEntry) string {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to format log entry", err)
		return ""
	}

	return string(jsonBytes)
}
