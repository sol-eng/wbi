package workbench

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
)

// Prompt users if they would like to install Workbench
func WorkbenchInstallPrompt() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Workbench is required to be installed to continue. Would you like to install Workbench?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Workbench install prompt")
	}
	return name, nil
}

func PromptInstallVerify() (bool, error) {
	name := false
	prompt := &survey.Confirm{
		Message: "Would you like to verify the installation of Workbench?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with verify Workbench install prompt")
	}
	return name, nil
}
