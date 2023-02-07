package config

type AuthType int

const (
	SAML AuthType = iota
	OIDC
	LDAP
	PAM
	Other
)

func (d AuthType) String() string {
	return [...]string{"SAML", "OIDC", "LDAP", "PAM", "Other"}[d]
}

// WBConfig stores the entire workbench configuration
type WBConfig struct {
	SSLConfig    SSLConfig
	AuthType     AuthType
	RConfig      RConfig
	PythonConfig PythonConfig
	OIDCConfig   OIDCConfig
	SAMLConfig   SAMLConfig
}

type RConfig struct {
	Paths []string
}

type PythonConfig struct {
	Paths []string
}

// SSLConfig stores SSL config
type SSLConfig struct {
	KeyPath  string
	CertPath string
	UseSSL   bool
}

// OIDCConfig stores OIDC config
type OIDCConfig struct {
	AuthOpenID              int
	ClientID                string
	ClientSecret            string
	AuthOpenIDIssuer        string
	AuthOpenIDUsernameClaim string
}

// SAMLConfig stores SAML config
type SAMLConfig struct {
	AuthSAML                    int
	AuthSamlSpAttributeUsername string
	AuthSamlMetadataURL         string
}
