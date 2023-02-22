package license

import (
	"fmt"
	"strings"

	"github.com/dpastoor/wbi/internal/system"
)

func CheckLicenseActivation() (bool, error) {

	stdout, _, err := system.RunCommandAndCaptureOutput("rstudio-server license-manager status")
	if err != nil {
		return false, fmt.Errorf("issue checking license activation: %w", err)
	}

	if strings.Contains(stdout, "Status: Activated") {
		fmt.Println("\nAn active Workbench license was detected\n")
		return true, nil
	}

	fmt.Println("\nNo active Workbench license was detected\n")
	return false, nil
}
