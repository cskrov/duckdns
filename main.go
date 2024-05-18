package main

import (
	"os"
	"time"

	"github.com/cskrov/duckdns/configuration"
	"github.com/cskrov/duckdns/log"
	"github.com/cskrov/duckdns/request"
	"github.com/cskrov/duckdns/response"
)

func main() {
	config := configuration.LoadConfig()

	if config == nil {
		os.Exit(1)
	}

	isValid := configuration.ValidateConfig(config)

	if !isValid {
		os.Exit(1)
	}

	for {
		resposeString, err := request.UpdateDuckDNS(config)

		if err != nil {
			log.LogError("Failed to update DuckDNS", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		res, err := response.ParseResponse(resposeString)

		if err != nil {
			log.LogError("Failed to parse DuckDNS response", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		if config.Log {
			log.LogResponse(config, res)
		}

		// Sleep for 5 minutes
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}
