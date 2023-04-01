package workbench

import (
	"fmt"
	"os"

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

			fmt.Println("\n=== Writing to the file " + filepath + ":")
			err := system.WriteStrings(writeLines, filepath, 0644)
			if err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
		} else {
			return fmt.Errorf("line already exists in repos.conf")
		}
	} else if source == "pypi" {
		// TODO instead of removing ensure that the correct URL is set
		// Remove pip.conf if it exists
		if _, err := os.Stat("/etc/pip.conf"); err == nil {
			err = os.Remove("/etc/pip.conf")
			if err != nil {
				return fmt.Errorf("failed to remove pip.conf: %w", err)
			}
		}

		writeLines := []string{
			"[global]",
			"index-url=" + url,
		}
		filepath := "/etc/pip.conf"

		fmt.Println("\n=== Writing to the file " + filepath + ":")
		err := system.WriteStrings(writeLines, filepath, 0644)
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
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

		fmt.Println("\n=== Writing to the file " + filepath + ":")
		err := system.WriteStrings(writeLines, filepath, 0644)
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

		fmt.Println("\n=== Writing to the file " + filepath + ":")
		err := system.WriteStrings(writeLines, filepath, 0644)
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	} else {
		return fmt.Errorf("line already exists in rsession.conf")
	}
	return nil
}

// WriteSAMLAuthConfig writes the SAML auth config to the Workbench config file
func WriteSAMLAuthConfig(idpUrl string) error {
	// check to ensure the lines don't already exist
	linesCheck := []string{
		"auth-saml=",
		"auth-saml-metadata-url=",
	}
	filepath := "/etc/rstudio/rserver.conf"

	var linesExist bool
	for _, value := range linesCheck {
		matched, err := system.CheckStringExists(value, filepath)
		if err == nil && matched {
			linesExist = true
		}
	}

	if !linesExist {
		writeLines := []string{
			"auth-saml=1",
			"auth-saml-metadata-url=" + idpUrl,
		}

		fmt.Println("\n=== Writing to the file " + filepath + ":")
		err := system.WriteStrings(writeLines, filepath, 0644)
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	} else {
		return fmt.Errorf("at least one line already exists in rserver.conf")
	}

	return nil
}

// WriteOIDCAuthConfig writes the OIDC auth config to the Workbench config file
func WriteOIDCAuthConfig(idpURL string, usernameClaim string, clientID string, clientSecret string) error {
	// default in the configuration https://docs.posit.co/ide/server-pro/authenticating_users/openid_connect_authentication.html#openid-claims
	if usernameClaim == "" {
		usernameClaim = "preferred_username"
	}

	// check to ensure the lines don't already exist
	linesCheck := []string{
		"auth-openid=",
		"auth-openid-issuer=",
		"auth-openid-username-claim=",
	}
	filepathRserver := "/etc/rstudio/rserver.conf"

	var linesExist bool
	for _, value := range linesCheck {
		matched, err := system.CheckStringExists(value, filepathRserver)
		if err == nil && matched {
			linesExist = true
		}
	}

	if !linesExist {
		// rserver.conf config options
		writeLinesRserver := []string{
			"auth-openid=1",
			"auth-openid-issuer=" + idpURL,
			"auth-openid-username-claim=" + usernameClaim,
		}

		fmt.Println("\n=== Writing to the file " + filepathRserver + ":")
		err := system.WriteStrings(writeLinesRserver, filepathRserver, 0644)
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	} else {
		return fmt.Errorf("at least one line already exists in rserver.conf")
	}

	if clientID != "" && clientSecret != "" {

		// check to ensure the lines don't already exist
		linesCheck := []string{
			"client-id=",
			"client-secret=",
		}
		filepathClientSecret := "/etc/rstudio/openid-client-secret"

		var linesExist bool
		for _, value := range linesCheck {
			matched, err := system.CheckStringExists(value, filepathClientSecret)
			if err == nil && matched {
				linesExist = true
			}
		}

		if !linesExist {
			// openid-client-secret config options
			writeLinesClientSecret := []string{
				"client-id=" + clientID,
				"client-secret=" + clientSecret,
			}

			fmt.Println("\n=== Writing to the file " + filepathClientSecret + ":")
			err := system.WriteStrings(writeLinesClientSecret, filepathClientSecret, 0644)
			if err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
		} else {
			return fmt.Errorf("at least one line already exists in openid-client-secret")
		}
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

	fmt.Println("\n=== Writing to the file " + filepath + ":")
	err := system.WriteStrings(writeLines, filepath, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
