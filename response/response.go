package response

import (
	"errors"
	"strings"
)

type OkResponse struct {
	IPv4    string
	IPv6    string
	Updated bool
}

func ParseResponse(res string) (*OkResponse, error) {
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
