package workbench

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/prompt"
)

// Prompt users if they would like to install Workbench
func WorkbenchInstallPrompt() (bool, error) {
	confirmText := "Workbench is required to be installed to continue. Would you like to install Workbench?"

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in Workbench install confirm prompt: %w", err)
	}

	return result, nil
}

func PromptInstallVerify() (bool, error) {
	confirmText := "Would you like to verify the installation of Workbench?"

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in verify install confirm prompt: %w", err)
	}

	return result, nil
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
