package license

import (
	"fmt"

	"github.com/sol-eng/wbi/internal/system"
)

// Activate Workbench based on a license key
func ActivateLicenseKey(licenseKey string) error {
	cmdLicense := "rstudio-server license-manager activate " + licenseKey
	err := system.RunCommand(cmdLicense)
	if err != nil {
		return fmt.Errorf("issue activating Workbench license: %w", err)
	}

	// TODO add a real check that Workbench is activated
	fmt.Println("\nWorkbench has been successfully activated")
	return nil
}
