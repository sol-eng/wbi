package workbench

import (
	"fmt"
	"os/exec"
)

// Checks if Workbench is installed
func VerifyWorkbench() bool {
	cmd := exec.Command("/bin/sh", "-c", "rstudio-server version")
	stdout, err := cmd.Output()

	if err != nil {
		return false
	} else {
		fmt.Println("\nWorkbench installation detected: ", string(stdout))
		return true
	}
}
