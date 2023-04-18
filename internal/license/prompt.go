package license

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
)

// Prompt users if they wish to activate Workbench with a license key
func PromptLicenseChoice() (bool, error) {
	name := true
	messageText := "Would you like to activate Workbench with a license key?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Workbench activation prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

// Prompt users for a Workbench license key
func PromptLicense() (string, error) {
	target := ""
	messageText := "Workbench license key:"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a license key: %w", err)
	}
	log.Info(messageText)
	log.Info(target)
	return target, nil
}

func CheckPromptAndActivateLicense() error {
	licenseActivationStatus, err := CheckLicenseActivation()
	if err != nil {
		return fmt.Errorf("issue in checking for license activation: %w", err)
	}

	if !licenseActivationStatus {
		licenseChoice, err := PromptLicenseChoice()
		if err != nil {
			return fmt.Errorf("issue in prompt for license activate choice: %w", err)
		}

		if licenseChoice {
			licenseKey, err := PromptLicense()
			if err != nil {
				return fmt.Errorf("issue entering license key: %w", err)
			}
			ActivateErr := ActivateLicenseKey(licenseKey)
			if ActivateErr != nil {
				return fmt.Errorf("issue activating license key: %w", ActivateErr)
			}
		}
	}
	return nil
}
