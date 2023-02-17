package install

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Create a temporary file and download the R/Python installer to it.
func (installerInfo *InstallerInfo) DownloadLanguage(language string) (string, error) {

	url := installerInfo.URL
	name := installerInfo.Name

	languageTitleCase := strings.Title(language)

	fmt.Println("Downloading " + languageTitleCase + " installer from: " + url)

	// Create the file
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return tmpFile.Name(), err
	}
	defer tmpFile.Close()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, url, nil)
	if err != nil {
		return "", errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return "", errors.New("error downloading " + languageTitleCase + " installer")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("error retrieving " + languageTitleCase + " installer")
	}

	// Writer the body to file
	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
