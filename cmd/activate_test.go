package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	type errorTestCases struct {
		description   string
		opts          activateOpts
		arguments     []string
		expectedError string
	}

	for _, scenario := range []errorTestCases{
		{
			description:   "no arguments",
			opts:          activateOpts{key: "1234"},
			arguments:     []string{},
			expectedError: "no arguments provided, please provide one argument",
		},
		{
			description: "too many arguments, more than 1",
			opts:        activateOpts{key: "1234"},
			arguments: []string{
				"license",
				"somethingelse",
			},
			expectedError: "too many arguments provided, please provide only one argument",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			err := scenario.opts.Validate(scenario.arguments)
			require.Error(t, err)
			assert.Equal(t, scenario.expectedError, err.Error())
		})
	}
}
