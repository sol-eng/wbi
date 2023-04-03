package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestScanParamsValidate tests the scan command parameters
func TestScanParamsValidate(t *testing.T) {
	tests := map[string]struct {
		args        []string
		flags       scanOpts
		expectError string
	}{
		// general arguement tests
		"no argument or key": {
			args:        []string{},
			flags:       scanOpts{},
			expectError: "no arguments provided, please provide one argument",
		},
		"too many arguments": {
			args:        []string{"r", "python"},
			flags:       scanOpts{},
			expectError: "too many arguments provided, please provide only one argument",
		},
		// r argument tests
		"r argument only succeeds": {
			args:        []string{"r"},
			flags:       scanOpts{},
			expectError: "",
		},
		// python argument tests
		"python argument only succeeds": {
			args:        []string{"python"},
			flags:       scanOpts{},
			expectError: "",
		},
		// non r or python argument test
		"non r or python argument only fails": {
			args:        []string{"workbench"},
			flags:       scanOpts{},
			expectError: "invalid language provided, please provide one of the following: r, python",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			scanCmd := newScanCmd()
			// set the flags
			scanCmd.opts = tc.flags
			// run validation
			err := scanCmd.opts.Validate(tc.args)

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

// TestScanRCommandIntegration tests the scan command with the r arg in a Docker container.
func TestScanRCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	installCommand := []string{"./wbi", "scan", "r"}
	successMessage := []string{"/opt/R/4.2.2/bin/R"}

	IntegrationContainerRunner(t, "Dockerfile.Ubuntu-R-Python", installCommand, successMessage, false)
}

// TestScanPythonCommandIntegration tests the scan command with the python arg in a Docker container.
func TestScanPythonCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	installCommand := []string{"./wbi", "scan", "python"}
	successMessage := []string{
		"/opt/python/3.11.2/bin/python",
		"/usr/bin/python3",
	}

	IntegrationContainerRunner(t, "Dockerfile.Ubuntu-R-Python", installCommand, successMessage, false)
}
