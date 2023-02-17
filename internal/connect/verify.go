package connect

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func cleanConnectURL(connectURL string) string {
	// remove trailing slash if present
	if connectURL[len(connectURL)-1] == '/' {
		connectURL = connectURL[:len(connectURL)-1]
	}
	// if the url does not contain :// then add it, either http if port 3939 or https in other cases (most installs)
	if !strings.Contains(connectURL, "://") {
		if strings.Contains(connectURL, ":3939") {
			connectURL = "http://" + connectURL
		} else {
			connectURL = "https://" + connectURL
		}
	}
	return connectURL
}

// Information is only needed to check if the URL is valid
type ProhibitedUsernames struct {
	ProhibitedUsernames []string `json:"prohibited_usernames"`
}

// VerifyConnectURL checks if the Connect URL is valid
func VerifyConnectURL(connectURL string) (ProhibitedUsernames, error) {

	fullTestURL := cleanConnectURL(connectURL) + "/__api__/server_settings"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, fullTestURL, nil)
	if err != nil {
		return ProhibitedUsernames{}, errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return ProhibitedUsernames{}, errors.New("error retrieving JSON data")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return ProhibitedUsernames{}, errors.New("error in HTTP status code")
	}
	var prohibitedUsernames ProhibitedUsernames
	err = json.NewDecoder(res.Body).Decode(&prohibitedUsernames)
	if err != nil {
		return ProhibitedUsernames{}, errors.New("error unmarshalling JSON data")
	}

	fmt.Println("\nConnect URL has been successfull validated.\n")
	return prohibitedUsernames, nil
}
