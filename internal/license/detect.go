package license

import (
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckLicenseActivation() (bool, error) {
	licenseActivateCommand := "rstudio-server license-manager status"
	licenseStatus, err := system.RunCommandAndCaptureOutput(licenseActivateCommand, false, 0, false)
	if err != nil {
		return false, fmt.Errorf("issue checking license activation with command '%s': %w", licenseActivateCommand, err)
	}

	if strings.Contains(licenseStatus, "Status: Activated") {
		system.PrintAndLogInfo("\nAn active Workbench license was detected")
		return true, nil
	}

	system.PrintAndLogInfo("\nNo active Workbench license was detected")
	return false, nil
}
