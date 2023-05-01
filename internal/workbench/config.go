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

			system.PrintAndLogInfo("\n=== Writing to the file " + filepath + ":")
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

			system.PrintAndLogInfo("\n=== Writing to the file " + filepath + ":")
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

// WriteSSLConfig writes the SSL config to the Workbench config file
func WriteSSLConfig(certPath string, keyPath string) error {
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

	// append the lines if they don't exist
	if !linesExist {
		writeLines := []string{
			"ssl-enabled=1",
			"ssl-certificate=" + certPath,
			"ssl-certificate-key=" + keyPath,
		}

		system.PrintAndLogInfo("\n=== Writing to the file " + filepath + ":")
		err := system.WriteStrings(writeLines, filepath, 0644, true)
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

		system.PrintAndLogInfo("\n=== Writing to the file " + filepath + ":")
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

	system.PrintAndLogInfo("\n=== Writing to the file " + filepath + ":")
	err := system.WriteStrings(writeLines, filepath, 0644, true)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
