package system

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Runs a command in the terminal and streams the output
func RunCommand(command string) error {
	fmt.Println("Running command: " + command)
	// sleep for 3 seconds to allow the user to read the command
	time.Sleep(3 * time.Second)

	cmd := exec.Command("/bin/sh", "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("issue running command: %w", err)
	}

	return nil
}
