package connect

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sol-eng/wbi/internal/system"
)

func cleanConnectURL(connectURL string) string {
	// remove trailing slash if present
	if connectURL[len(connectURL)-1] == '/' {
		connectURL = connectURL[:len(connectURL)-1]
	}
	return connectURL
}

// VerifyConnectURL checks if the Connect URL is valid
func VerifyConnectURL(connectURL string) (string, error) {

	cleanConnectURL := cleanConnectURL(connectURL)
	fullTestURL := cleanConnectURL + "/__ping__"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, fullTestURL, nil)
	if err != nil {
		return "", errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return "", errors.New("error retrieving JSON data")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("error in HTTP status code")
	}

	system.PrintAndLogInfo("\nConnect URL has been successfull validated.")
	return cleanConnectURL, nil
}
