package license

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// Prompt users if they wish to activate Workbench with a license key
func PromptLicenseChoice() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to activate Workbench with a license key?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Workbench activation prompt")
	}
	return name, nil
}

// Prompt users for a Workbench license key
func PromptLicense() (string, error) {
	target := ""
	prompt := &survey.Input{
		Message: "Workbench license key:",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a license key: %w", err)
	}
	return target, nil
}
