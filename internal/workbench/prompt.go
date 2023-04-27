package workbench

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sol-eng/wbi/internal/config"
)

// Prompt users if they would like to install Workbench
func WorkbenchInstallPrompt() (bool, error) {
	name := true
	messageText := "Workbench is required to be installed to continue. Would you like to install Workbench?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Workbench install prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

func PromptInstallVerify() (bool, error) {
	name := false
	messageText := "Would you like to verify the installation of Workbench?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with verify Workbench install prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

func CheckPromptDownloadAndInstallWorkbench(osType config.OperatingSystem) error {
	workbenchInstalled := VerifyWorkbench()
	if !workbenchInstalled {
		installWorkbenchChoice, err := WorkbenchInstallPrompt()
		if err != nil {
			return fmt.Errorf("issue selecting Workbench installation: %w", err)
		}
		if installWorkbenchChoice {
			err := DownloadAndInstallWorkbench(osType)
			if err != nil {
				return fmt.Errorf("issue installing Workbench: %w", err)
			}
		} else {
			log.Fatal("Workbench installation is required to continue")
		}
	}
	return nil
}
