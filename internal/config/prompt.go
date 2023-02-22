package config

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
)

// PromptWriteConfig prompts the user if they would like the new values to be written to configuration files
func PromptWriteConfig() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "The configuration changes shown above need to be written to the configuration files. Would you like wbi to write these changes and restart rstudio-server and rstudio-launcher?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the config change writing prompt")
	}
	return name, nil
}
