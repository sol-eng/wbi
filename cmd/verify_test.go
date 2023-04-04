package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestVerifyParamsValidate tests the verify command parameters
func TestVerifyParamsValidate(t *testing.T) {
	tests := map[string]struct {
		args        []string
		flags       verifyOpts
		expectError string
	}{
		// general arguement tests
		"no argument or key": {
			args:        []string{},
			flags:       verifyOpts{},
			expectError: "no arguments provided, please provide one argument",
		},
		"too many arguments": {
			args:        []string{"r", "python"},
			flags:       verifyOpts{},
			expectError: "too many arguments provided, please provide only one argument",
		},
		// packagemanager argument tests
		"packagemanager argument only fails": {
			args:        []string{"packagemanager"},
			flags:       verifyOpts{},
			expectError: "the url flag is required for packagemanager",
		},
		"packagemanager argument with URL flag succeeds": {
			args:        []string{"packagemanager"},
			flags:       verifyOpts{url: "https://packagemanager.posit.co/"},
			expectError: "",
		},
		"packagemanager argument with URL and repo flags fail": {
			args:        []string{"packagemanager"},
			flags:       verifyOpts{url: "https://packagemanager.posit.co/", repo: "cran"},
			expectError: "the language flag is required when the repo flag is provided",
		},
		"packagemanager argument with URL and language flags fail": {
			args:        []string{"packagemanager"},
			flags:       verifyOpts{url: "https://packagemanager.posit.co/", language: "r"},
			expectError: "the repo flag is required when the language flag is provided",
		},
		"packagemanager argument with URL, language and repo flag succeeds": {
			args:        []string{"packagemanager"},
			flags:       verifyOpts{url: "https://packagemanager.posit.co/", language: "r", repo: "cran"},
			expectError: "",
		},
		"packagemanager argument with cert-path flag fails": {
			args:        []string{"packagemanager"},
			flags:       verifyOpts{certPath: "cert.crt"},
			expectError: "the cert-path flag is only supported for ssl",
		},
		"packagemanager argument with key-path flag fails": {
			args:        []string{"packagemanager"},
			flags:       verifyOpts{keyPath: "cert.key"},
			expectError: "the key-path flag is only supported for ssl",
		},
		// connect-url argument tests
		"connect-url argument only fails": {
			args:        []string{"connect-url"},
			flags:       verifyOpts{},
			expectError: "the url flag is required for connect-url",
		},
		"connect-url argument with URL flag succeeds": {
			args:        []string{"connect-url"},
			flags:       verifyOpts{url: "https://colorado.posit.co/rsc"},
			expectError: "",
		},
		"connect-url argument with repo flag fails": {
			args:        []string{"connect-url"},
			flags:       verifyOpts{repo: "cran"},
			expectError: "the repo flag is only supported for packagemanager",
		},
		"connect-url argument with language flag fails": {
			args:        []string{"connect-url"},
			flags:       verifyOpts{language: "r"},
			expectError: "the language flag is only supported for packagemanager",
		},
		"connect-url argument with cert-path flag fails": {
			args:        []string{"connect-url"},
			flags:       verifyOpts{certPath: "cert.crt"},
			expectError: "the cert-path flag is only supported for ssl",
		},
		"connect-url argument with key-path flag fails": {
			args:        []string{"connect-url"},
			flags:       verifyOpts{keyPath: "cert.key"},
			expectError: "the key-path flag is only supported for ssl",
		},
		// workbench argument tests
		"workbench argument only succeeds": {
			args:        []string{"workbench"},
			flags:       verifyOpts{},
			expectError: "",
		},
		"workbench argument and a url flag fails": {
			args:        []string{"workbench"},
			flags:       verifyOpts{url: "https://colorado.posit.co/rsc"},
			expectError: "the url flag is only supported for packagemanager and connect",
		},
		"workbench argument and a language flag fails": {
			args:        []string{"workbench"},
			flags:       verifyOpts{language: "r"},
			expectError: "the language flag is only supported for packagemanager",
		},
		"workbench argument and a repo flag fails": {
			args:        []string{"workbench"},
			flags:       verifyOpts{repo: "cran"},
			expectError: "the repo flag is only supported for packagemanager",
		},
		"workbench argument and a cert-path flag fails": {
			args:        []string{"workbench"},
			flags:       verifyOpts{certPath: "cert.crt"},
			expectError: "the cert-path flag is only supported for ssl",
		},
		"workbench argument and a key-path flag fails": {
			args:        []string{"workbench"},
			flags:       verifyOpts{keyPath: "cert.key"},
			expectError: "the key-path flag is only supported for ssl",
		},
		// ssl argument tests
		"ssl argument only fails": {
			args:        []string{"ssl"},
			flags:       verifyOpts{},
			expectError: "the cert-path flag is required for ssl",
		},
		"ssl argument and a cert-path flag fails": {
			args:        []string{"ssl"},
			flags:       verifyOpts{certPath: "cert.crt"},
			expectError: "the key-path flag is required for ssl",
		},
		"ssl argument, a cert-path and a key-path flag succeeds": {
			args:        []string{"ssl"},
			flags:       verifyOpts{certPath: "cert.crt", keyPath: "cert.key"},
			expectError: "",
		},
		"ssl argument and a url flag fails": {
			args:        []string{"ssl"},
			flags:       verifyOpts{url: "https://colorado.posit.co/rsc"},
			expectError: "the url flag is only supported for packagemanager and connect",
		},
		"ssl argument and a language flag fails": {
			args:        []string{"ssl"},
			flags:       verifyOpts{language: "r"},
			expectError: "the language flag is only supported for packagemanager",
		},
		"ssl argument and a repo flag fails": {
			args:        []string{"ssl"},
			flags:       verifyOpts{repo: "cran"},
			expectError: "the repo flag is only supported for packagemanager",
		},
		// license argument tests
		"license argument only succeeds": {
			args:        []string{"license"},
			flags:       verifyOpts{},
			expectError: "",
		},
		"license argument and a url flag fails": {
			args:        []string{"license"},
			flags:       verifyOpts{url: "https://colorado.posit.co/rsc"},
			expectError: "the url flag is only supported for packagemanager and connect",
		},
		"license argument and a language flag fails": {
			args:        []string{"license"},
			flags:       verifyOpts{language: "r"},
			expectError: "the language flag is only supported for packagemanager",
		},
		"license argument and a repo flag fails": {
			args:        []string{"license"},
			flags:       verifyOpts{repo: "cran"},
			expectError: "the repo flag is only supported for packagemanager",
		},
		"license argument and a cert-path flag fails": {
			args:        []string{"license"},
			flags:       verifyOpts{certPath: "cert.crt"},
			expectError: "the cert-path flag is only supported for ssl",
		},
		"license argument and a key-path flag fails": {
			args:        []string{"license"},
			flags:       verifyOpts{keyPath: "cert.key"},
			expectError: "the key-path flag is only supported for ssl",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			verifyCmd := newVerifyCmd()
			// set the flags
			verifyCmd.opts = tc.flags
			// run validation
			err := verifyCmd.opts.Validate(tc.args)

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

// TestVerifyPackageManagerCommandIntegration tests the verify command with the packagemanager arg with just a URL flag in a Docker container.
func TestVerifyPackageManagerURLCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	installCommand := []string{"./wbi", "verify", "packagemanager", "--url=https://packagemanager.posit.co"}
	successMessage := []string{"Posit Package Manager URL has been successfull validated."}

	IntegrationContainerRunner(t, "Dockerfile.Ubuntu", installCommand, successMessage, false)
}

// TestVerifyPackageManagerURLRepoCommandIntegration tests the verify command with the packagemanager arg and URL, repo and language flags in a Docker container.
func TestVerifyPackageManagerURLRepoCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	installCommand := []string{"./wbi", "verify", "packagemanager", "--url=https://packagemanager.posit.co", "--repo cran", "--language r"}
	successMessage := []string{
		"Posit Package Manager URL has been successfull validated.",
		"Posit Package Manager Repository has been successfull validated.",
	}

	IntegrationContainerRunner(t, "Dockerfile.Ubuntu", installCommand, successMessage, false)
}

// TestVerifyConnectURLCommandIntegration tests the verify command with the connect-url arg with the URL flag in a Docker container.
func TestVerifyConnectURLCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	installCommand := []string{"./wbi", "verify", "connect-url", "--url=https://colorado.posit.co/rsc"}
	successMessage := []string{"Connect URL has been successfull validated."}

	IntegrationContainerRunner(t, "Dockerfile.Ubuntu", installCommand, successMessage, false)
}

// TestVerifyWorkbenchCommandIntegration tests the verify command with the workbench arg in a Docker container.
func TestVerifyWorkbenchCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	installCommand := []string{"./wbi", "verify", "workbench"}
	successMessage := []string{"Workbench installation detected:"}

	IntegrationContainerRunner(t, "Dockerfile.Workbench", installCommand, successMessage, false)
}

// TODO TestVerifySSLCommandIntegration tests the verify command with the ssl arg with the cert-path and key-path flags in a Docker container.

// TestVerifyLicenseCommandIntegration tests the verify command with the license arg in a Docker container.
func TestVerifyLicenseCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	installCommand := []string{"./wbi", "verify", "license"}
	successMessage := []string{"An active Workbench license was detected"}

	IntegrationContainerRunner(t, "Dockerfile.Workbench", installCommand, successMessage, false)
}
