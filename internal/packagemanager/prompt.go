package packagemanager

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// Prompt users if they wish to add a default Posit Package Manager URL to Workbench
func PromptPackageManagerChoice() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to provide a default Posit Package Manager URL for Workbench?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Posit Package Manager URL prompt")
	}
	return name, nil
}

// Prompt users for a default Posit Package Manager URL
func PromptPackageManagerURL() (string, error) {
	target := ""
	prompt := &survey.Input{
		Message: "Enter your Posit Package Manager base URL (for example, https://exampleaddress.com):",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a Posit Package Manager URL: %w", err)
	}
	return target, nil
}

// Prompt users for a Posit Package Manager repo name
func PromptPackageManagerRepo() (string, error) {
	target := ""
	prompt := &survey.Input{
		Message: "Enter the name of your CRAN repository on Posit Package Manager (for example, prod-cran):",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a Posit Package Manager repo: %w", err)
	}
	return target, nil
}
