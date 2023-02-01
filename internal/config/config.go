package config

// WBConfig stores the entire workbench configuration
type WBConfig struct {
	SSLConfig SSLConfig
	AuthType  string
	// TODO: should probably nest this or otherwise
	LDAPConfig LDAPConfig
	OIDCConfig OIDCConfig
	SAMLConfig SAMLConfig
}

// SSLConfig stores SSL config
type SSLConfig struct {
	KeyPath  string
	CertPath string
	UseSSL   bool
}

// TODO: what actually do we need for ldap
type LDAPConfig struct {
}

// OIDCConfig stores OIDC config
type OIDCConfig struct {
}

// SAMLConfig stores SAML config
type SAMLConfig struct {
}
