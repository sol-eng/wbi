package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSetupParamsValidate tests the setup command parameters
func TestSetupParamsValidate(t *testing.T) {
	tests := map[string]struct {
		args        []string
		flags       setupOpts
		expectError string
	}{
		// general arguement tests
		"no argument or flag succeeds": {
			args:        []string{},
			flags:       setupOpts{},
			expectError: "",
		},
		"too many arguments": {
			args:        []string{"python"},
			flags:       setupOpts{},
			expectError: "no arguments are supported for this command",
		},
		"no argument with a valid step succeeds": {
			args:        []string{},
			flags:       setupOpts{step: "auth"},
			expectError: "",
		},
		"no argument with an invalid step fails": {
			args:        []string{},
			flags:       setupOpts{step: "configure"},
			expectError: "invalid step: configure",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			setupCmd := newSetupCmd()
			// set the flags
			setupCmd.opts = tc.flags
			// run validation
			err := setupCmd.opts.Validate(tc.args)

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

// TODO add integration tests for the setup prompts
