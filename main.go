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
			logError("Failed to update DuckDNS", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		response, err := parseResponse(resposeString)

		if err != nil {
			logError("Failed to parse DuckDNS response", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		if config.Log {
			logResponse(config, response)
		}

		// Sleep for 5 minutes
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

type Timestamp int64

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Unix(int64(t), 0).UTC().Format("2006-01-02T15:04:05Z") + "\""), nil
}

type LogEntry struct {
	Level     string    `json:"level"`
	Timestamp Timestamp `json:"timestamp"`
	Message   string    `json:"message"`
	Updated   bool      `json:"updated"`
	IPv4      string    `json:"ipv4,omitempty"`
	IPv6      string    `json:"ipv6,omitempty"`
}

func logResponse(config *Config, res *OkResponse) {
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

	log(formatLog(entry), os.Stdout)
}

func logError(message string, err error) {
	if err != nil {
		message += ": " + err.Error()
	}

	log(formatLog(createLogEntry("error", message)), os.Stderr)
}

func createLogEntry(level string, message string) *LogEntry {
	return &LogEntry{
		Timestamp: Timestamp(time.Now().UTC().Unix()),
		Level:     level,
		Message:   message,
	}
}

func formatLog(entry *LogEntry) string {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to format log entry", err)
		return ""
	}

	return string(jsonBytes)
}

func log(line string, output *os.File) {
	if line == "" {
		return
	}

	fmt.Fprintln(output, line)

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
		logError("Failed to open config file", err)
		return nil
	}

	// Parse configuration
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		logError("Failed to parse config file", err)
		return nil
	}

	return &config
}

func validateConfig(config *Config) bool {
	valid := true

	if config.Token == "" {
		logError("Token is required", nil)
		valid = false
	}

	if len(config.Domains) == 0 {
		logError("At least one subdomain is required", nil)
		valid = false
	}

	if config.Interval == 0 {
		logError("Interval is required", nil)
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
