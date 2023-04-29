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
		"ssl argument with cert-path and key-path flags succeeds": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key"},
			expectError: "",
		},
		"ssl argument with url flag fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key", url: "https://www.packagemanager.rstudio.com"},
			expectError: "the url flag is only valid for repo and connect-url",
		},
		"ssl argument with source flag fails": {
			args:        []string{"ssl"},
			flags:       configOpts{certPath: "cert.crt", keyPath: "cert.key", source: "cran"},
			expectError: "the source flag is only valid for repo",
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
