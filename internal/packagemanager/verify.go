package packagemanager

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sol-eng/wbi/internal/system"
)

type RepoInformation []struct {
	Name string `json:"name"`
}

func cleanPackageManagerURL(packageManagerURL string) string {
	// remove trailing slash if present
	if packageManagerURL[len(packageManagerURL)-1] == '/' {
		packageManagerURL = packageManagerURL[:len(packageManagerURL)-1]
	}
	return packageManagerURL
}

func VerifyPackageManagerURL(packageManagerURL string) (string, error) {
	cleanPackageManagerURL := cleanPackageManagerURL(packageManagerURL)
	fullTestURL := cleanPackageManagerURL + "/__ping__"

	client := &http.Client{
		Timeout: 5 * time.Second,
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

	system.PrintAndLogInfo("\nPosit Package Manager URL has been successfull validated.")

	return cleanPackageManagerURL, nil
}

func VerifyPackageManagerRepo(packageManagerURL string, packageManagerRepo string, language string) error {
	repoSearchURL := packageManagerURL + "/__api__/repos?type=" + language

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, repoSearchURL, nil)
	if err != nil {
		return errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return errors.New("error retrieving JSON data")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("error in HTTP status code")
	}

	var repoInformation RepoInformation
	err = json.NewDecoder(res.Body).Decode(&repoInformation)
	if err != nil {
		return errors.New("error unmarshalling JSON data")
	}

	// verify the repo name exists in the list of repos
	matchFound := false
	for _, repo := range repoInformation {
		if repo.Name == packageManagerRepo {
			matchFound = true
			break
		}
	}
	if !matchFound {
		return errors.New("error finding the " + packageManagerRepo + " repository in Posit Package Manager")
	}

	system.PrintAndLogInfo("\nPosit Package Manager Repository has been successfull validated.")
	return nil

}
