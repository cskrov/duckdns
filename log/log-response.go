package log

import (
	"os"

	"github.com/cskrov/duckdns/response"
	"github.com/cskrov/duckdns/types"
)

func LogResponse(config *types.Config, res *response.OkResponse) {
	if !res.Updated && !config.VerboseLog {
		return
	}

	var message string
	if res.Updated {
		message = "Updated"
	} else if config.VerboseLog {
		message = "No change"
	}

	entry := createLogEntry("info", message)

	entry.Updated = res.Updated
	entry.IPv4 = res.IPv4
	entry.IPv6 = res.IPv6

	logLine(formatLog(entry), os.Stdout)
}
