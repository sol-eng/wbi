package workbench

import (
	"fmt"
	"os/exec"

	"github.com/sol-eng/wbi/internal/system"
)

// Checks if Workbench is installed
func VerifyWorkbench() bool {
	cmd := exec.Command("/bin/sh", "-c", "rstudio-server version")
	stdout, err := cmd.Output()

	if err != nil {
		return false
	} else {
		system.PrintAndLogInfo("\nWorkbench installation detected: " + string(stdout))
		return true
	}
}

// Runs verify-installation command
func VerifyInstallation() error {
	// stop rstudio-server
	err := StopRStudioServer()
	if err != nil {
		return fmt.Errorf("issue stopping rstudio-server: %w", err)
	}
	// run verify-installation
	err = system.RunCommand("rstudio-server verify-installation", true, 1)
	if err != nil {
		return fmt.Errorf("issue running verify-installation command 'rstudio-server verify-installation': %w", err)
	}
	// start rstudio-server
	err = StartRStudioServer()
	if err != nil {
		return fmt.Errorf("issue starting rstudio-server: %w", err)
	}
	return nil
}
