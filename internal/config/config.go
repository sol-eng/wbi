package config

type AuthType int

const (
	SAML AuthType = iota + 1
	OIDC
	LDAP
	PAM
	Other
)

func (d AuthType) String() string {
	return [...]string{"SAML", "OIDC", "LDAP", "PAM", "Other"}[d]
}

type OperatingSystem int

const (
	Unknown OperatingSystem = iota
	Ubuntu18
	Ubuntu20
	Ubuntu22
	Redhat7
	Redhat8
)

// WBConfig stores the entire workbench configuration
type WBConfig struct {
	SSLConfig            SSLConfig
	RConfig              RConfig
	PythonConfig         PythonConfig
	AuthConfig           AuthConfig
	PackageManagerConfig PackageManagerConfig
	ConnectConfig        ConnectConfig
	ConfigChange         bool
}

type AuthConfig struct {
	Using      bool
	AuthType   AuthType
	OIDCConfig OIDCConfig
	SAMLConfig SAMLConfig
}

type RConfig struct {
	Paths []string
}

type PythonConfig struct {
	Paths       []string
	JupyterPath string
}

// SSLConfig stores SSL config
type SSLConfig struct {
	Using    bool
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

type PackageManagerConfig struct {
	Using     bool
	RURL      string
	PythonURL string
}

type ConnectConfig struct {
	Using bool
	URL   string
}
