package workbench

import (
	"fmt"

	"github.com/sol-eng/wbi/internal/system"
)

// WriteRepoConfig writes the repo config to the Workbench config file
func WriteRepoConfig(url string, source string) error {
	if source == "cran" {
		filepath := "/etc/rstudio/repos.conf"
		// check to ensure the line doesn't already exist
		lineExists, err := system.CheckStringExists("CRAN=", filepath)
		if err != nil {
			return fmt.Errorf("failed to check if line exists: %w", err)
		}

		if !lineExists {
			writeLines := []string{
				"CRAN=" + url,
			}

			err := system.WriteStrings(writeLines, filepath, 0644, true)
			if err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
		} else {
			return fmt.Errorf("line already exists in repos.conf")
		}
	} else if source == "pypi" {
		filepath := "/etc/pip.conf"
		// check to ensure the line doesn't already exist
		lineExists, err := system.CheckStringExists("index-url="+url, filepath)
		if err != nil {
			return fmt.Errorf("failed to check if line exists: %w", err)
		}

		if !lineExists {
			writeLines := []string{
				"[global]",
				"index-url=" + url,
			}

			err := system.WriteStrings(writeLines, filepath, 0644, true)
			if err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
		} else {
			return fmt.Errorf("line already exists in pip.conf")
		}
	}
	return nil
}

func cleanServerURL(serverURL string) string {
	// remove trailing slash if present
	if serverURL[len(serverURL)-1] == '/' {
		serverURL = serverURL[:len(serverURL)-1]
	}
	// remove http:// or https:// if present
	if serverURL[:7] == "http://" {
		serverURL = serverURL[7:]
	} else if serverURL[:8] == "https://" {
		serverURL = serverURL[8:]
	}
	return serverURL
}

// WriteSSLConfig writes the SSL config to the Workbench config file
func WriteSSLConfig(certPath string, keyPath string, serverURL string) error {
	// clean the serverURL
	serverURLClean := cleanServerURL(serverURL)
	finalServerURL := "https://" + serverURLClean

	// check to ensure the lines don't already exist
	linesCheck := []string{
		"ssl-enabled=",
		"ssl-certificate=",
		"ssl-certificate-key=",
	}
	filepath := "/etc/rstudio/rserver.conf"

	var linesExist bool
	for _, value := range linesCheck {
		matched, err := system.CheckStringExists(value, filepath)
		if err == nil && matched {
			linesExist = true
		}
	}

	// remove launcher-sessions-callback-address and append the lines if they don't exist
	if !linesExist {
		// remove launcher-sessions-callback-address
		err := system.DeleteStrings([]string{"launcher-sessions-callback-address"}, filepath, 0644)
		if err != nil {
			return fmt.Errorf("failed to delete launcher-sessions-callback-address: %w", err)
		}

		// append new lines including the updated launcher-sessions-callback-address
		writeLines := []string{
			"",
			"launcher-sessions-callback-address=" + finalServerURL,
			"",
			"ssl-enabled=1",
			"ssl-certificate=" + certPath,
			"ssl-certificate-key=" + keyPath,
		}

		err = system.WriteStrings(writeLines, filepath, 0644, true)
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	} else {
		return fmt.Errorf("at least one line already exists in rserver.conf")
	}

	return nil
}

// WriteConnectURLConfig writes the Connect URL config to the Workbench config file
func WriteConnectURLConfig(url string) error {
	// check to ensure the line doesn't already exist
	filepath := "/etc/rstudio/rsession.conf"
	lineExists, err := system.CheckStringExists("default-rsconnect-server=", filepath)
	if err != nil {
		return fmt.Errorf("failed to check if line exists: %w", err)
	}

	if !lineExists {
		writeLines := []string{
			"default-rsconnect-server=" + url,
		}

		err := system.WriteStrings(writeLines, filepath, 0644, true)
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	} else {
		return fmt.Errorf("line already exists in rsession.conf")
	}
	return nil
}

// WriteJupyterConfig writes the Jupyter config to the Workbench config file
func WriteJupyterConfig(jupyterPath string) error {
	// TODO check to ensure line doesn't already exist and ideally take out the default commented out line to reduce confusion
	writeLines := []string{
		"jupyter-exe=" + jupyterPath,
	}
	filepath := "/etc/rstudio/jupyter.conf"

	err := system.WriteStrings(writeLines, filepath, 0644, true)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
