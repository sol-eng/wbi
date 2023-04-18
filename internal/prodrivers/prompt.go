package prodrivers

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/system"
)

// Prompt users if they would like to install Posit Pro Drivers
func ProDriversInstallPrompt() (bool, error) {
	name := true
	messageText := "Would you like to install Post Pro Drivers?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Pro Drivers install prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

func CheckPromptDownloadAndInstallProDrivers(osType config.OperatingSystem) error {
	proDriversExistingStatus, err := CheckExistingProDrivers()
	if err != nil {
		return fmt.Errorf("issue in checking for prior pro driver installation: %w", err)
	}
	if !proDriversExistingStatus {
		installProDriversChoice, err := ProDriversInstallPrompt()
		if err != nil {
			return fmt.Errorf("issue selecting Pro Drivers installation: %w", err)
		}
		if installProDriversChoice {
			err := DownloadAndInstallProDrivers(osType)
			if err != nil {
				return fmt.Errorf("issue installing Pro Drivers: %w", err)
			}
		}
	} else {
		system.PrintAndLogInfo("Pro Drivers are already installed")
	}
	return nil
}
