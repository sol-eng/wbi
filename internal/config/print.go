package config

import "fmt"

// Prints the WBConfig configuration struct information to the console
func (WBConfig *WBConfig) ConfigStructToText() {
	anyPrint := false
	if WBConfig.PythonConfig.JupyterPath != "" {
		anyPrint = true
		WBConfig.PythonConfig.PythonStructToText()
	}
	if WBConfig.SSLConfig.UseSSL {
		anyPrint = true
		WBConfig.SSLConfig.SSLStructToText()
	}
	if WBConfig.AuthConfig.Using {
		anyPrint = true
		WBConfig.AuthConfig.AuthStructToText()
	}
	if WBConfig.PackageManagerConfig.Using {
		anyPrint = true
		WBConfig.PackageManagerStringToText()
	}
	if WBConfig.ConnectConfig.Using {
		anyPrint = true
		WBConfig.ConnectStringToText()
	}
	if anyPrint {
		fmt.Println("\n=== Please restart Workbench after making these changes")
	} else {
		fmt.Println("\n=== No configuration changes are needed")
	}
}

// Prints the PythonConfig configuration struct information to the console
func (PythonConfig *PythonConfig) PythonStructToText() {
	fmt.Println("\n=== Add to config file: /etc/rstudio/jupyter.conf:")
	fmt.Println("jupyter-exe=" + PythonConfig.JupyterPath)
}

// Prints the SSLConfig configuration struct information to the console
func (SSLConfig *SSLConfig) SSLStructToText() {
	fmt.Println("\n=== Add to config file: /etc/rstudio/rserver.conf:")
	fmt.Println("ssl-enabled=1")
	fmt.Println("ssl-certificate=" + SSLConfig.CertPath)
	fmt.Println("ssl-certificate-key=" + SSLConfig.KeyPath)
}

// Prints the AuthConfig configuration struct information to the console
func (AuthConfig *AuthConfig) AuthStructToText() {
	switch AuthConfig.AuthType {
	case SAML:
		AuthConfig.SAMLConfig.AuthSAMLStructToText()
	case OIDC:
		AuthConfig.OIDCConfig.AuthOIDCStructToText()
	default:

	}
}

// Prints the SAMLConfig configuration struct information to the console
func (SAMLConfig *SAMLConfig) AuthSAMLStructToText() {
	fmt.Println("\n=== Add to config file: /etc/rstudio/rserver.conf:")
	fmt.Println("auth-saml=1")
	fmt.Println("auth-saml-sp-attribute-username=" + SAMLConfig.AuthSamlSpAttributeUsername)
	fmt.Println("auth-saml-metadata-url=" + SAMLConfig.AuthSamlMetadataURL)
}

// Prints the OIDCConfig configuration struct information to the console
func (OIDCConfig *OIDCConfig) AuthOIDCStructToText() {
	fmt.Println("\n=== Add to config file: /etc/rstudio/rserver.conf:")
	fmt.Println("auth-openid=1")
	fmt.Println("auth-openid-issuer=" + OIDCConfig.AuthOpenIDIssuer)
	fmt.Println("auth-openid-username-claim=" + OIDCConfig.AuthOpenIDUsernameClaim)

	fmt.Println("\n=== Add to config file: /etc/rstudio/openid-client-secret:")
	fmt.Println("client-id=" + OIDCConfig.ClientID)
	fmt.Println("client-secret=" + OIDCConfig.ClientSecret)
}

// Prints the Package Manager URL configuration string information to the console
func (WBConfig *WBConfig) PackageManagerStringToText() {
	fmt.Println("\n=== Add to config file: /etc/rstudio/repos.conf:")
	fmt.Println("CRAN=" + WBConfig.PackageManagerConfig.URL)
}

// Prints the ConnectURL configuration string information to the console
func (WBConfig *WBConfig) ConnectStringToText() {
	fmt.Println("\n=== Add to config file: /etc/rstudio/rsession.conf:")
	fmt.Println("default-rsconnect-server=" + WBConfig.ConnectConfig.URL)
}
