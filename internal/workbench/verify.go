package workbench

import (
	"fmt"
	"os/exec"
)

func VerifyWorkbench() bool {
	cmd := exec.Command("/bin/sh", "-c", "sudo rstudio-server version")
	stdout, err := cmd.Output()

	if err != nil {
		return false
	} else {
		fmt.Println("Workbench installation detected: ", string(stdout))
		return true
	}
}
