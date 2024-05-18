package request

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/cskrov/duckdns/types"
)

const DUCKDNS_URL = "https://www.duckdns.org/update"

func UpdateDuckDNS(config *types.Config) (string, error) {
	query := CreateQueryParams(config)
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

func CreateQueryParams(config *types.Config) string {
	return "domains=" + strings.Join(config.Domains, ",") + "&token=" + config.Token
}
