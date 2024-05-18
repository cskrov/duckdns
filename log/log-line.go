package log

import (
	"fmt"
	"os"
	"time"
)

func logLine(line string, output *os.File) {
	if line == "" {
		return
	}

	fmt.Fprintln(output, line)

	date := time.Now().UTC().Format("01-2006")
	fileName := "duckdns-log-" + date + ".log"
	stat, err := os.Stat("/data/logs")
	if os.IsNotExist(err) {
		err = os.Mkdir("/data/logs", 0770)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create log directory", err)
		}
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to check if log directory exists", err)
	} else if !stat.IsDir() {
		fmt.Fprintln(os.Stderr, "/data/logs exists but is not a directory")
		return
	}

	file, err := os.OpenFile("/data/logs/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open log file", err)
		return
	}

	_, err = file.WriteString(line + "\n")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write to log file "+fileName, err)
	}
}
