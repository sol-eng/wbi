package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestActivateParamsValidate tests the activate command parameters
func TestActivateParamsValidate(t *testing.T) {

	tests := map[string]struct {
		args        []string
		flags       activateOpts
		expectError string
	}{
		"no argument": {
			args:        []string{},
			flags:       activateOpts{},
			expectError: "no arguments provided, please provide one argument",
		},
		"too many arguments": {
			args:        []string{"license", "secondarg"},
			flags:       activateOpts{},
			expectError: "too many arguments provided, please provide only one argument",
		},
		"no key flag": {
			args:        []string{"license"},
			flags:       activateOpts{},
			expectError: "the key flag is required for license",
		},
		"correct command but incorrect key flag": {
			args:        []string{"license"},
			flags:       activateOpts{key: ""},
			expectError: "the key flag is required for license",
		},
		"correct command and key flag": {
			args:        []string{"license"},
			flags:       activateOpts{key: "1234"},
			expectError: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			activateCmd := newActivateCmd()
			// set the flags
			activateCmd.opts = tc.flags
			// run validation
			err := activateCmd.opts.Validate(tc.args)

			if err != nil {
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

// TestActivateCommandIntegration tests the activate command in a Docker container. $RSW_LICENSE must be set to a valid license key.
func TestActivateCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	command := []string{"./wbi", "activate", "license", fmt.Sprintf("--key=%s", os.Getenv("RSW_LICENSE"))}
	successMessage := []string{"Workbench has been successfully activated"}

	IntegrationContainerRunner(t, "Dockerfile.Workbench", command, successMessage, false)
}
