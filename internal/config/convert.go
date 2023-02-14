package config

import "fmt"

func (WBConfig *WBConfig) ConfigStructToText() {
	WBConfig.PythonConfig.PythonStructToText()
	WBConfig.SSLConfig.SSLStructToText()
	WBConfig.AuthConfig.AuthStructToText()
}

func (PythonConfig *PythonConfig) PythonStructToText() {
	fmt.Println("\n=== Adding to config file: /etc/rstudio/jupyter.conf:")
	fmt.Println("jupyter-exe=" + PythonConfig.JupyterPath)
}

func (SSLConfig *SSLConfig) SSLStructToText() {
	fmt.Println("\n=== Adding to config file: /etc/rstudio/rserver.conf:")
	fmt.Println("ssl-enabled=1")
	fmt.Println("ssl-certificate=" + SSLConfig.CertPath)
	fmt.Println("ssl-certificate-key=" + SSLConfig.KeyPath)
}

func (AuthConfig *AuthConfig) AuthStructToText() {
	switch AuthConfig.AuthType {
	case SAML:
		AuthConfig.SAMLConfig.AuthSAMLStructToText()
	case OIDC:
		AuthConfig.OIDCConfig.AuthOIDCStructToText()
	default:

	}
}

func (SAMLConfig *SAMLConfig) AuthSAMLStructToText() {
	fmt.Println("\n=== Adding to config file: /etc/rstudio/rserver.conf:")
	fmt.Println("auth-saml=1")
	fmt.Println("auth-saml-sp-attribute-username=" + SAMLConfig.AuthSamlSpAttributeUsername)
	fmt.Println("auth-saml-metadata-url=" + SAMLConfig.AuthSamlMetadataURL)
}

func (OIDCConfig *OIDCConfig) AuthOIDCStructToText() {
	fmt.Println("\n=== Adding to config file: /etc/rstudio/rserver.conf:")
	fmt.Println("auth-openid=1")
	fmt.Println("auth-openid-issuer=" + OIDCConfig.AuthOpenIDIssuer)
	fmt.Println("auth-openid-username-claim=" + OIDCConfig.AuthOpenIDUsernameClaim)

	fmt.Println("\n=== Adding to config file: /etc/rstudio/openid-client-secret:")
	fmt.Println("client-id=" + OIDCConfig.ClientID)
	fmt.Println("client-secret=" + OIDCConfig.ClientSecret)
}
