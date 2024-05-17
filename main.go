package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	Token      string   `json:"token"`
	Domains    []string `json:"domains"`
	Log        bool     `json:"log"`
	VerboseLog bool     `json:"verbose_log"`
	Interval   int      `json:"interval"`
}

const DUCKDNS_URL = "https://www.duckdns.org/update"

func main() {
	config := loadConfig()

	if config == nil {
		os.Exit(1)
	}

	isValid := validateConfig(config)

	if !isValid {
		os.Exit(1)
	}

	for {
		resposeString, err := updateDuckDNS(config)

		if err != nil {
			log("Failed to update DuckDNS: " + err.Error())
			time.Sleep(1 * time.Minute)
			continue
		}

		response, err := parseResponse(resposeString)

		if err != nil {
			log("Failed to parse DuckDNS response: " + err.Error())
			time.Sleep(1 * time.Minute)
			continue
		}

		if config.Log {
			line := formatLog(config, response)

			if line != "" {
				log(line)
			}
		}

		// Sleep for 5 minutes
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

func formatLog(config *Config, res *OkResponse) string {
	if !res.Updated && !config.VerboseLog {
		return ""
	}

	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
	message := timestamp + " - "

	if res.Updated {
		message += "Updated: "
	} else if config.VerboseLog {
		message += "No change: "
	}

	if res.IPv6 == "" {
		return message + res.IPv4
	}

	return message + res.IPv4 + " / " + res.IPv6
}

func log(line string) {
	println(line)

	date := time.Now().UTC().Format("01-2006")
	fileName := "log-" + date + ".log"
	file, err := os.OpenFile("/data/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open log file", err)
	}

	_, err = file.WriteString(line + "\n")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write to log file "+fileName, err)
	}
}

func loadConfig() *Config {
	file, err := os.Open("/data/config.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	// Parse configuration
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return &config
}

func validateConfig(config *Config) bool {
	valid := true

	if config.Token == "" {
		fmt.Fprintln(os.Stderr, "Token is required")
		valid = false
	}

	if len(config.Domains) == 0 {
		fmt.Fprintln(os.Stderr, "At least one subdomain is required")
		valid = false
	}

	if config.Interval == 0 {
		fmt.Fprintln(os.Stderr, "Interval is required")
		valid = false
	}

	return valid
}

func updateDuckDNS(config *Config) (string, error) {
	query := createQueryParams(config)
	url := DUCKDNS_URL + "?" + query + "&verbose=true"
	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	responseString := string(bodyBytes)

	if strings.HasPrefix(responseString, "KO") {
		return "", errors.New("DuckDNS responded with KO")
	}

	return responseString, nil
}

func createQueryParams(config *Config) string {
	return "domains=" + strings.Join(config.Domains, ",") + "&token=" + config.Token
}

type OkResponse struct {
	IPv4    string
	IPv6    string
	Updated bool
}

func parseResponse(res string) (*OkResponse, error) {
	parts := strings.Split(res, "\n")

	status := parts[0]

	if status != "OK" {
		return nil, errors.New("Unexpected response status from DuckDNS: " + status)
	}

	if len(parts) == 1 {
		return nil, errors.New("Invalid response from DuckDNS\n\n" + res)
	}

	hasIPv6 := strings.Contains(parts[2], ":")

	if hasIPv6 {
		return &OkResponse{
			IPv4:    parts[1],
			IPv6:    parts[2],
			Updated: parts[3] == "UPDATED",
		}, nil
	}

	return &OkResponse{
		IPv4:    parts[1],
		Updated: parts[2] == "UPDATED",
	}, nil
}
