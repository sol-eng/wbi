package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInstallParamsValidate tests the install command parameters
func TestInstallParamsValidate(t *testing.T) {
	//TODO: add valid R and Python versions from the availble versions function in languages package instead of hardcode
	tests := map[string]struct {
		args        []string
		flags       installOpts
		expectError string
	}{
		// general arguement tests
		"no argument or key": {
			args:        []string{},
			flags:       installOpts{},
			expectError: "no arguments provided, please provide one argument",
		},
		"too many arguments": {
			args:        []string{"workbench", "r"},
			flags:       installOpts{},
			expectError: "too many arguments provided, please provide only one argument",
		},
		// r argument tests
		"r argument only succeeds": {
			args:        []string{"r"},
			flags:       installOpts{},
			expectError: "",
		},
		"r argument with a symlink flag succeeds": {
			args:        []string{"r"},
			flags:       installOpts{symlink: true},
			expectError: "",
		},
		"r argument with a valid r version and a symlink flag succeeds": {
			args:        []string{"r"},
			flags:       installOpts{symlink: true, versions: []string{"3.6.1"}},
			expectError: "",
		},
		"r argument with an invalid r version and a symlink flag fails": {
			args:        []string{"r"},
			flags:       installOpts{symlink: true, versions: []string{"1.6.1"}},
			expectError: "version 1.6.1 is not a valid R version",
		},
		"r argument with single valid version flag succeeds": {
			args:        []string{"r"},
			flags:       installOpts{versions: []string{"3.6.1"}},
			expectError: "",
		},
		"r argument with single invalid version flag fails": {
			args:        []string{"r"},
			flags:       installOpts{versions: []string{"1.6.1"}},
			expectError: "version 1.6.1 is not a valid R version",
		},
		"r argument with multiple valid version flags succeeds": {
			args:        []string{"r"},
			flags:       installOpts{versions: []string{"3.6.1", "4.2.2"}},
			expectError: "",
		},
		"r argument with multiple invalid version flags fails": {
			args:        []string{"r"},
			flags:       installOpts{versions: []string{"1.6.1", "2.2.2"}},
			expectError: "version 1.6.1 is not a valid R version",
		},
		"r argument with one valid and one invalid version flags fails": {
			args:        []string{"r"},
			flags:       installOpts{versions: []string{"3.6.1", "2.2.2"}},
			expectError: "version 2.2.2 is not a valid R version",
		},
		"r argument with path flag fails": {
			args:        []string{"r"},
			flags:       installOpts{path: "/opt/python/3.6.1/bin/python"},
			expectError: "the path flag is only supported for jupyter",
		},
		"r argument with addToPATH flag fails": {
			args:        []string{"r"},
			flags:       installOpts{addToPATH: true},
			expectError: "the add-to-path flag is only supported for python",
		},
		// python argument tests
		"python argument only succeeds": {
			args:        []string{"python"},
			flags:       installOpts{},
			expectError: "",
		},
		"python argument with a symlink flag fails": {
			args:        []string{"python"},
			flags:       installOpts{symlink: true},
			expectError: "the symlink flag is only supported for r",
		},
		"python argument with a addToPATH flag succeeds": {
			args:        []string{"python"},
			flags:       installOpts{addToPATH: true},
			expectError: "",
		},
		"python argument with a valid python version and a addToPATH flag succeeds": {
			args:        []string{"python"},
			flags:       installOpts{addToPATH: true, versions: []string{"3.11.0"}},
			expectError: "",
		},
		"python argument with an invalid python version and a addToPATH flag fails": {
			args:        []string{"python"},
			flags:       installOpts{addToPATH: true, versions: []string{"2.11.1"}},
			expectError: "version 2.11.1 is not a valid Python version",
		},
		"python argument with single valid version flag succeeds": {
			args:        []string{"python"},
			flags:       installOpts{versions: []string{"3.11.0"}},
			expectError: "",
		},
		"python argument with single invalid version flag fails": {
			args:        []string{"python"},
			flags:       installOpts{versions: []string{"1.6.1"}},
			expectError: "version 1.6.1 is not a valid Python version",
		},
		"python argument with multiple valid version flags succeeds": {
			args:        []string{"python"},
			flags:       installOpts{versions: []string{"3.11.0", "3.10.0"}},
			expectError: "",
		},
		"python argument with multiple invalid version flags fails": {
			args:        []string{"python"},
			flags:       installOpts{versions: []string{"3.15.1", "2.2.2"}},
			expectError: "version 3.15.1 is not a valid Python version",
		},
		"python argument with one valid and one invalid version flags fails": {
			args:        []string{"python"},
			flags:       installOpts{versions: []string{"3.11.0", "3.3.2"}},
			expectError: "version 3.3.2 is not a valid Python version",
		},
		"python argument with path flag fails": {
			args:        []string{"python"},
			flags:       installOpts{path: "/opt/python/3.6.1/bin/python"},
			expectError: "the path flag is only supported for jupyter",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			installCmd := newInstallCmd()
			// set the flags
			installCmd.opts = tc.flags
			// run validation
			err := installCmd.opts.Validate(tc.args)

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
