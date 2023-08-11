package license

import (
	"fmt"

	"github.com/sol-eng/wbi/internal/prompt"
)

// Prompt users if they wish to activate Workbench with a license key
func PromptLicenseChoice() (bool, error) {
	confirmText := "Would you like to activate Workbench with a license key?"

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in Workbench license activation confirm prompt: %w", err)
	}

	return result, nil
}

// Prompt users for a Workbench license key
func PromptLicense() (string, error) {
	promptText := "Workbench license key:"

	result, err := prompt.PromptText(promptText)
	if err != nil {
		return "", fmt.Errorf("issue occured in license key text prompt: %w", err)
	}

	return result, nil
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
