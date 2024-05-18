package configuration

import (
	"encoding/json"
	"os"

	"github.com/cskrov/duckdns/log"
	"github.com/cskrov/duckdns/types"
)

func LoadConfig() *types.Config {
	file, err := os.Open("/data/config.json")
	if err != nil {
		log.LogError("Failed to open config file", err)
		return nil
	}

	// Parse configuration
	config := &types.Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		log.LogError("Failed to parse config file", err)
		return nil
	}

	return config
}

func ValidateConfig(config *types.Config) bool {
	valid := true

	if config.Token == "" {
		log.LogError("Token is required", nil)
		valid = false
	}

	if len(config.Domains) == 0 {
		log.LogError("At least one subdomain is required", nil)
		valid = false
	}

	if config.Interval == 0 {
		log.LogError("Interval is required", nil)
		valid = false
	}

	return valid
}
