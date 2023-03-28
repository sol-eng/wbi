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
	status, err := system.RunCommandAndCaptureOutput("rstudio-server status | cat", true, 1)
	if err != nil {
		return fmt.Errorf("issue running status for rstudio-server: %w", err)
	}
	if strings.Contains(status, "active (running)") {
		fmt.Println(status)
		fmt.Println("\nrstudio-server status is active (running)!\n")
	} else {
		fmt.Println(status)
		fmt.Println("\nrstudio-server status is not active!\n")
	}

	return nil
}

func StatusRStudioLauncher() error {
	status, err := system.RunCommandAndCaptureOutput("rstudio-launcher status | cat", true, 1)
	if err != nil {
		return fmt.Errorf("issue running status for rstudio-launcher with the command 'rstudio-launcher status | cat': %w", err)
	}
	if strings.Contains(status, "active (running)") {
		fmt.Println(status)
		fmt.Println("\nrstudio-launcher status is active (running)!\n")
	} else {
		fmt.Println(status)
		fmt.Println("\nrstudio-launcher status is not active!\n")
	}
	return nil
}
