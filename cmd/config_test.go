package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfigParamsValidate tests the config command parameters
func TestConfigParamsValidate(t *testing.T) {
	tests := map[string]struct {
		args        []string
		flags       configOpts
		expectError string
	}{
		// general arguement tests
		"no argument or key": {
			args:        []string{},
			flags:       configOpts{},
			expectError: "no arguments provided, please provide one argument",
		},
		"too many arguments": {
			args:        []string{"r", "python"},
			flags:       configOpts{},
			expectError: "too many arguments provided, please provide only one argument",
		},
		// ssl argument tests
		"ssl argument only fails": {
			args:        []string{"ssl"},
			flags:       configOpts{},
			expectError: "the cert-path flag is required for ssl",
		},
		"ssl argument with cert-path flag fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt"},
			expectError: "the key-path flag is required for ssl",
		},
		"ssl argument with cert-path and key-path flags fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key"},
			expectError: "the url flag is required for ssl",
		},
		"ssl argument with cert-path, key-path and url flags succeeds": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key", url: "myserverurl.com"},
			expectError: "",
		},
		"ssl argument with source flag fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key", url: "myserverurl.com", source: "cran"},
			expectError: "the source flag is only valid for repo",
		},
		"ssl argument with auth-type flag fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key", url: "myserverurl.com", authType: "saml"},
			expectError: "the auth-type flag is only valid for auth",
		},
		"ssl argument with idp-url flag fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key", url: "myserverurl.com", idpURL: "https://www.example.com"},
			expectError: "the idp-url flag is only valid with auth as an argument and a auth-type flag of saml or oidc",
		},
		"ssl argument with username-claim flag fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key", url: "myserverurl.com", usernameClaim: "user"},
			expectError: "the username-claim flag is only valid with auth as an argument and a auth-type flag of oidc",
		},
		"ssl argument with client-id flag fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key", url: "myserverurl.com", clientID: "user"},
			expectError: "the client-id flag is only valid with auth as an argument and a auth-type flag of oidc",
		},
		"ssl argument with client-secret flag fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key", url: "myserverurl.com", clientSecret: "user"},
			expectError: "the client-secret flag is only valid with auth as an argument and a auth-type flag of oidc",
		},
		// auth argument tests
		"auth argument only fails": {
			args:        []string{"auth"},
			flags:       configOpts{},
			expectError: "the auth-type flag is required for auth",
		},
		"auth argument with the auth-type flag fails": {
			args:        []string{"auth"},
			flags:       configOpts{authType: "saml"},
			expectError: "the idp-url flag is required for argument auth and auth-type flag of saml or odic",
		},
		"auth argument with the auth-type saml and idp-url flags set succeeds": {
			args:        []string{"auth"},
			flags:       configOpts{authType: "saml", idpURL: "https://www.example.com"},
			expectError: "",
		},
		"auth argument with the auth-type oidc and idp-url flags set fails": {
			args:        []string{"auth"},
			flags:       configOpts{authType: "oidc", idpURL: "https://www.example.com"},
			expectError: "the client-id flag is required for argument auth and auth-type flag of odic",
		},
		"auth argument with the auth-type oidc and idp-url and client-id flags set fails": {
			args:        []string{"auth"},
			flags:       configOpts{authType: "oidc", idpURL: "https://www.example.com", clientID: "dadawdaw"},
			expectError: "the client-secret flag is required for argument auth and auth-type flag of odic",
		},
		"auth argument with the auth-type oidc and idp-url, client-id, and client-secret flags set succeeds": {
			args:        []string{"auth"},
			flags:       configOpts{authType: "oidc", idpURL: "https://www.example.com", clientID: "adwada", clientSecret: "adawdaw"},
			expectError: "",
		},
		"auth argument with url flag fails": {
			args:        []string{"auth"},
			flags:       configOpts{authType: "saml", idpURL: "https://www.example.com", url: "https://www.packagemanager.rstudio.com"},
			expectError: "the url flag is only valid for repo, connect-url and url",
		},
		"auth argument with source flag fails": {
			args:        []string{"auth"},
			flags:       configOpts{authType: "saml", idpURL: "https://www.example.com", source: "cran"},
			expectError: "the source flag is only valid for repo",
		},
		"auth argument with cert-path flag fails": {
			args:        []string{"auth"},
			flags:       configOpts{authType: "saml", idpURL: "https://www.example.com", certPath: "cert.crt"},
			expectError: "the cert-path flag is only valid for ssl",
		},
		"auth argument with key-path flag fails": {
			args:        []string{"auth"},
			flags:       configOpts{authType: "saml", idpURL: "https://www.example.com", keyPath: "cert.key"},
			expectError: "the key-path flag is only valid for ssl",
		},
		// repo argument tests
		"repo argument only fails": {
			args:        []string{"repo"},
			flags:       configOpts{},
			expectError: "the url flag is required for repo",
		},
		"repo argument with a URL flag fails": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co"},
			expectError: "the source flag is required for repo",
		},
		"repo argument with a URL and source of cran flags succeeds": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co", source: "cran"},
			expectError: "",
		},
		"repo argument with a URL and source of pypi flags succeeds": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co", source: "pypi"},
			expectError: "",
		},
		"repo argument with auth-type flag fails": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co", source: "cran", authType: "saml"},
			expectError: "the auth-type flag is only valid for auth",
		},
		"repo argument with idp-url flag fails": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co", source: "cran", idpURL: "https://www.example.com"},
			expectError: "the idp-url flag is only valid with auth as an argument and a auth-type flag of saml or oidc",
		},
		"repo argument with username-claim flag fails": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co", source: "cran", usernameClaim: "user"},
			expectError: "the username-claim flag is only valid with auth as an argument and a auth-type flag of oidc",
		},
		"repo argument with client-id flag fails": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co", source: "cran", clientID: "user"},
			expectError: "the client-id flag is only valid with auth as an argument and a auth-type flag of oidc",
		},
		"repo argument with client-secret flag fails": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co", source: "cran", clientSecret: "user"},
			expectError: "the client-secret flag is only valid with auth as an argument and a auth-type flag of oidc",
		},
		"repo argument with cert-path flag fails": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co", source: "cran", certPath: "cert.crt"},
			expectError: "the cert-path flag is only valid for ssl",
		},
		"repo argument with key-path flag fails": {
			args:        []string{"repo"},
			flags:       configOpts{url: "https://packagemanager.posit.co", source: "cran", keyPath: "cert.key"},
			expectError: "the key-path flag is only valid for ssl",
		},
		// connect-url argument tests
		"connect-url argument only fails": {
			args:        []string{"connect-url"},
			flags:       configOpts{},
			expectError: "the url flag is required for connect-url",
		},
		"connect-url argument with a url flag succeeds": {
			args:        []string{"connect-url"},
			flags:       configOpts{url: "https://colorado.posit.co/rsc"},
			expectError: "",
		},
		"connect-url argument with a source flag succeeds": {
			args:        []string{"connect-url"},
			flags:       configOpts{url: "https://colorado.posit.co/rsc", source: "cran"},
			expectError: "the source flag is only valid for repo",
		},
		"connect-url argument with auth-type flag fails": {
			args:        []string{"connect-url"},
			flags:       configOpts{url: "https://packagemanager.posit.co", authType: "saml"},
			expectError: "the auth-type flag is only valid for auth",
		},
		"connect-url argument with idp-url flag fails": {
			args:        []string{"connect-url"},
			flags:       configOpts{url: "https://packagemanager.posit.co", idpURL: "https://www.example.com"},
			expectError: "the idp-url flag is only valid with auth as an argument and a auth-type flag of saml or oidc",
		},
		"connect-url argument with username-claim flag fails": {
			args:        []string{"connect-url"},
			flags:       configOpts{url: "https://packagemanager.posit.co", usernameClaim: "user"},
			expectError: "the username-claim flag is only valid with auth as an argument and a auth-type flag of oidc",
		},
		"connect-url argument with client-id flag fails": {
			args:        []string{"connect-url"},
			flags:       configOpts{url: "https://packagemanager.posit.co", clientID: "user"},
			expectError: "the client-id flag is only valid with auth as an argument and a auth-type flag of oidc",
		},
		"connect-url argument with client-secret flag fails": {
			args:        []string{"connect-url"},
			flags:       configOpts{url: "https://packagemanager.posit.co", clientSecret: "user"},
			expectError: "the client-secret flag is only valid with auth as an argument and a auth-type flag of oidc",
		},
		"connect-url argument with cert-path flag fails": {
			args:        []string{"connect-url"},
			flags:       configOpts{url: "https://packagemanager.posit.co", certPath: "cert.crt"},
			expectError: "the cert-path flag is only valid for ssl",
		},
		"connect-url argument with key-path flag fails": {
			args:        []string{"connect-url"},
			flags:       configOpts{url: "https://packagemanager.posit.co", keyPath: "cert.key"},
			expectError: "the key-path flag is only valid for ssl",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			configCmd := newConfigCmd()
			// set the flags
			configCmd.opts = tc.flags
			// run validation
			err := configCmd.opts.Validate(tc.args)

			if err != nil && tc.expectError != "" {
				// if we expect an error, check that it contains the expected error
				assert.Containsf(t, err.Error(), tc.expectError, "expected error containing %q, got %s", tc.expectError, err)
			} else if err != nil && tc.expectError == "" {
				// if we expect no error but get one then fail
				t.Fatalf("expected no error, but got %s", err)
			} else if err == nil && tc.expectError != "" {
				// if we expect an error but don't get one then fail
				t.Fatalf("expected error containing %q, but the command ran without error", tc.expectError)
			}
			// otherwise we expect the command to succeed so pass the test
		})
	}

}

// TODO TestConfigSSLCommandIntegration tests the config command with the ssl arg in a Docker container.

// TestConfigAuthSAMLCommandIntegration tests the config command with the auth arg, auth-type flag set to saml and idp-url flag set in a Docker container.
func TestConfigAuthSAMLCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// TODO actually verify the contents of rserver.conf
	t.Parallel()

	installCommand := []string{"./wbi", "config", "auth", "--auth-type=saml", "--idp-url=https://www.example.com"}
	successMessage := []string{"=== Writing to the file /etc/rstudio/rserver.conf:"}

	IntegrationContainerRunner(t, "Dockerfile.Workbench", installCommand, successMessage, false)
}

// TestConfigAuthOIDCCommandIntegration tests the config command with the auth arg, auth-type flag set to oidc, idp-url flag, client-id flag and client-secret flag set in a Docker container.
func TestConfigAuthOIDCCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// TODO actually verify the contents of rserver.conf
	t.Parallel()

	installCommand := []string{"./wbi", "config", "auth", "--auth-type=oidc", "--idp-url=https://www.example.com", "--client-id=awdawdawd", "--client-secret=adwawdawdfgawa"}
	successMessage := []string{"=== Writing to the file /etc/rstudio/rserver.conf:"}

	IntegrationContainerRunner(t, "Dockerfile.Workbench", installCommand, successMessage, false)
}

// TestConfigRepoCRANCommandIntegration tests the config command with the repo arg, url flag set to Public Package Manager and source flag set to cran in a Docker container.
func TestConfigRepoCRANCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// TODO actually verify the contents of repos.conf
	t.Parallel()

	installCommand := []string{"./wbi", "config", "repo", "--url=https://packagemanager.posit.co/cran/__linux__/jammy/latest", "--source=cran"}
	successMessage := []string{"=== Writing to the file /etc/rstudio/repos.conf:"}

	IntegrationContainerRunner(t, "Dockerfile.Workbench", installCommand, successMessage, false)
}

// TestConfigRepoPyPICommandIntegration tests the config command with the repo arg, url flag set to Public Package Manager and source flag set to pypi in a Docker container.
func TestConfigRepoPyPICommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// TODO actually verify the contents of pip.conf
	t.Parallel()

	installCommand := []string{"./wbi", "config", "repo", "--url=https://packagemanager.posit.co/pypi/latest/simple", "--source=pypi"}
	successMessage := []string{"=== Writing to the file /etc/pip.conf:"}

	IntegrationContainerRunner(t, "Dockerfile.Workbench", installCommand, successMessage, false)
}

// TestConfigConnectCommandIntegration tests the config command with the connect-url arg and the url flag set to Colorado Connect test server in a Docker container.
func TestConfigConnectCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// TODO actually verify the contents of resession.conf
	t.Parallel()

	installCommand := []string{"./wbi", "config", "connect-url", "--url=https://colorado.posit.co/rsc"}
	successMessage := []string{"=== Writing to the file /etc/rstudio/rsession.conf:"}

	IntegrationContainerRunner(t, "Dockerfile.Workbench", installCommand, successMessage, false)
}
