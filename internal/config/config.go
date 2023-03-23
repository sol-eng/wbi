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
	Redhat9
)
