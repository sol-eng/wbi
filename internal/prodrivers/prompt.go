package prodrivers

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
)

// Prompt users if they would like to install Posit Pro Drivers
func ProDriversInstallPrompt() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to install Post Pro Drivers?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Pro Drivers install prompt")
	}
	return name, nil
}
