package prodrivers

import (
	"fmt"

	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/prompt"
	"github.com/sol-eng/wbi/internal/system"
)

// Prompt users if they would like to install Posit Pro Drivers
func ProDriversInstallPrompt() (bool, error) {
	confirmText := "Would you like to install Post Pro Drivers?"

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in Pro Drivers install confirm prompt: %w", err)
	}

	return result, nil
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
