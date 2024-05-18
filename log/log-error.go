package log

import "os"

func LogError(message string, err error) {
	if err != nil {
		message += ": " + err.Error()
	}

	logLine(formatLog(createLogEntry("error", message)), os.Stderr)
}
