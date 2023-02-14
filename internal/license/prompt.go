package license

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

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
