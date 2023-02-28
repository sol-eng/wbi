package config

import (
	"bufio"
	"fmt"
	"os"
)

func WriteConfig(WBConfig WBConfig) error {
	if WBConfig.PythonConfig.JupyterPath != "" {
		err := WBConfig.PythonConfig.PythonStructToConfigWrite()
		if err != nil {
			return fmt.Errorf("failed to write python config: %w", err)
		}
	}
	if WBConfig.SSLConfig.UseSSL {
		err := WBConfig.SSLConfig.SSLStructToConfigWrite()
		if err != nil {
			return fmt.Errorf("failed to write SSL config: %w", err)
		}
	}
	if WBConfig.AuthConfig.Using {
		err := WBConfig.AuthConfig.AuthStructToConfigWrite()
		if err != nil {
			return fmt.Errorf("failed to write Authentication config: %w", err)
		}
	}
	if WBConfig.PackageManagerConfig.Using {
		err := WBConfig.PackageManagerStringToConfigWrite()
		if err != nil {
			return fmt.Errorf("failed to write Package Manager config: %w", err)
		}
	}
	if WBConfig.ConnectConfig.Using {
		err := WBConfig.ConnectStringToConfigWrite()
		if err != nil {
			return fmt.Errorf("failed to write Connect config: %w", err)
		}
	}
	return nil
}

func WriteStrings(lines []string, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range lines {
		_, err := datawriter.WriteString(data + "\n")
		if err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}
	}

	datawriter.Flush()
	file.Close()

	return nil
}

// Prints the PythonConfig configuration struct information to the console
func (PythonConfig *PythonConfig) PythonStructToConfigWrite() error {
	writeLines := []string{
		"jupyter-exe=" + PythonConfig.JupyterPath,
	}
	filepath := "/etc/rstudio/jupyter.conf"

	fmt.Println("\n=== Writing to the file " + filepath + ":")
	err := WriteStrings(writeLines, filepath)
	if err != nil {
		fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

// Prints the SSLConfig configuration struct information to the console
func (SSLConfig *SSLConfig) SSLStructToConfigWrite() error {
	writeLines := []string{
		"ssl-enabled=1",
		"ssl-certificate=" + SSLConfig.CertPath,
		"ssl-certificate-key=" + SSLConfig.KeyPath,
	}
	filepath := "/etc/rstudio/rserver.conf"

	fmt.Println("\n=== Writing to the file " + filepath + ":")
	err := WriteStrings(writeLines, filepath)
	if err != nil {
		fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

// Prints the AuthConfig configuration struct information to the console
func (AuthConfig *AuthConfig) AuthStructToConfigWrite() error {
	switch AuthConfig.AuthType {
	case SAML:
		err := AuthConfig.SAMLConfig.AuthSAMLStructToConfigWrite()
		if err != nil {
			return fmt.Errorf("failed to write SAML config: %w", err)
		}
	case OIDC:
		err := AuthConfig.OIDCConfig.AuthOIDCStructToConfigWrite()
		if err != nil {
			return fmt.Errorf("failed to write OIDC config: %w", err)
		}
	default:
		return nil
	}
	return nil
}

// Prints the SAMLConfig configuration struct information to the console
func (SAMLConfig *SAMLConfig) AuthSAMLStructToConfigWrite() error {
	writeLines := []string{
		"auth-saml=1",
		"auth-saml-sp-attribute-username=" + SAMLConfig.AuthSamlSpAttributeUsername,
		"auth-saml-metadata-url=" + SAMLConfig.AuthSamlMetadataURL,
	}
	filepath := "/etc/rstudio/rserver.conf"

	fmt.Println("\n=== Writing to the file " + filepath + ":")
	err := WriteStrings(writeLines, filepath)
	if err != nil {
		fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

// Prints the OIDCConfig configuration struct information to the console
func (OIDCConfig *OIDCConfig) AuthOIDCStructToConfigWrite() error {
	// rserver.conf config options
	writeLinesRserver := []string{
		"auth-openid=1",
		"auth-openid-issuer=" + OIDCConfig.AuthOpenIDIssuer,
		"auth-openid-username-claim=" + OIDCConfig.AuthOpenIDUsernameClaim,
	}
	filepathRserver := "/etc/rstudio/rserver.conf"

	fmt.Println("\n=== Writing to the file " + filepathRserver + ":")
	err := WriteStrings(writeLinesRserver, filepathRserver)
	if err != nil {
		fmt.Errorf("failed to write config: %w", err)
	}

	// openid-client-secret config options
	writeLinesClientSecret := []string{
		"client-id=" + OIDCConfig.ClientID,
		"client-secret=" + OIDCConfig.ClientSecret,
	}
	filepathClientSecret := "/etc/rstudio/openid-client-secret"

	fmt.Println("\n=== Writing to the file " + filepathClientSecret + ":")
	err = WriteStrings(writeLinesClientSecret, filepathClientSecret)
	if err != nil {
		fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Prints the Package Manager URL configuration string information to the console
func (WBConfig *WBConfig) PackageManagerStringToConfigWrite() error {
	if WBConfig.PackageManagerConfig.RURL != "" {
		writeLines := []string{
			"CRAN=" + WBConfig.PackageManagerConfig.RURL,
		}
		filepath := "/etc/rstudio/repos.conf"

		fmt.Println("\n=== Writing to the file " + filepath + ":")
		err := WriteStrings(writeLines, filepath)
		if err != nil {
			fmt.Errorf("failed to write config: %w", err)
		}
	}

	if WBConfig.PackageManagerConfig.PythonURL != "" {
		// Remove pip.conf if it exists
		if _, err := os.Stat("/etc/pip.conf"); err == nil {
			err = os.Remove("/etc/pip.conf")
			if err != nil {
				return fmt.Errorf("failed to remove pip.conf: %w", err)
			}
		}

		writeLines := []string{
			"[global]",
			"index-url=" + WBConfig.PackageManagerConfig.PythonURL,
		}
		filepath := "/etc/pip.conf"

		fmt.Println("\n=== Writing to the file " + filepath + ":")
		err := WriteStrings(writeLines, filepath)
		if err != nil {
			fmt.Errorf("failed to write config: %w", err)
		}
	}
	return nil
}

// Prints the ConnectURL configuration string information to the console
func (WBConfig *WBConfig) ConnectStringToConfigWrite() error {
	writeLines := []string{
		"default-rsconnect-server=" + WBConfig.ConnectConfig.URL,
	}
	filepath := "/etc/rstudio/rsession.conf"

	fmt.Println("\n=== Writing to the file " + filepath + ":")
	err := WriteStrings(writeLines, filepath)
	if err != nil {
		fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}
