package connect

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// Prompt users if they wish to add a default Connect URL to Workbench
func PromptConnectChoice() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to provide a default Connect URL for Workbench?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Connect URL prompt")
	}
	return name, nil
}

// Prompt users for a default Connect URL
func PromptConnectURL() (string, error) {
	target := ""
	prompt := &survey.Input{
		Message: "Connect URL:",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a Connect URL: %w", err)
	}
	return target, nil
}
