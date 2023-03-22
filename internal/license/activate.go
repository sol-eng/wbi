package license

import (
	"fmt"

	"github.com/sol-eng/wbi/internal/system"
)

// Activate Workbench based on a license key
func ActivateLicenseKey(licenseKey string) error {
	cmdLicense := "rstudio-server license-manager activate " + licenseKey
	err := system.RunCommand(cmdLicense, true, 1)
	if err != nil {
		return fmt.Errorf("issue activating Workbench license: %w", err)
	}

	// TODO add a real check that Workbench is activated
	fmt.Println("\nWorkbench has been successfully activated")
	return nil
}

// Check if Workbench is activated and if not perform the activation
func DetectAndActivateLicense(licenseKey string) error {
	licenseActivationStatus, err := CheckLicenseActivation()
	if err != nil {
		return fmt.Errorf("issue in checking for license activation: %w", err)
	}
	if !licenseActivationStatus {
		err := ActivateLicenseKey(licenseKey)
		if err != nil {
			return fmt.Errorf("issue activating license: %w", err)
		}
	} else {
		return fmt.Errorf("license is already activated")
	}
	return nil
}
