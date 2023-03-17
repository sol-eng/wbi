package workbench

import (
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func StatusRStudioServerAndLauncher() error {
	err := StatusRStudioServer()
	if err != nil {
		return fmt.Errorf("issue running status for rstudio-server: %w", err)
	}
	err = StatusRStudioLauncher()
	if err != nil {
		return fmt.Errorf("issue running status for rstudio-launcher: %w", err)
	}
	return nil
}

func StatusRStudioServer() error {
	stdout, _, err := system.RunCommandAndCaptureOutput("rstudio-server status")
	if err != nil {
		return fmt.Errorf("issue running status for rstudio-server: %w", err)
	}
	if strings.Contains(stdout, "active (running)") {
		fmt.Println("\nrstudio-server status showing active (running)!\n")
		return nil
	}

	return nil
}

func StatusRStudioLauncher() error {
	stdout, _, err := system.RunCommandAndCaptureOutput("rstudio-launcher status")
	if err != nil {
		return fmt.Errorf("issue running status for rstudio-launcher: %w", err)
	}
	if strings.Contains(stdout, "active (running)") {
		fmt.Println("\nrstudio-launcher status showing active (running)!\n")
		return nil
	}
	return nil
}
