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
	status, err := system.RunCommandAndCaptureOutput("rstudio-server status | cat")
	if err != nil {
		return fmt.Errorf("issue running status for rstudio-server: %w", err)
	}
	if strings.Contains(status, "active (running)") {
		fmt.Println("\nrstudio-server status is active (running)! See details below:\n")
		fmt.Println(status)
	} else {
		fmt.Println("\nrstudio-server status is inactive (dead)! See details below:\n")
		fmt.Println(status)
	}

	return nil
}

func StatusRStudioLauncher() error {
	status, err := system.RunCommandAndCaptureOutput("rstudio-launcher status | cat")
	if err != nil {
		return fmt.Errorf("issue running status for rstudio-launcher: %w", err)
	}
	if strings.Contains(status, "active (running)") {
		fmt.Println("\nrstudio-launcher status is active (running)! See details below:\n")
		fmt.Println(status)
	} else {
		fmt.Println("\nrstudio-launcher status is inactive (dead)! See details below:\n")
		fmt.Println(status)
	}
	return nil
}
