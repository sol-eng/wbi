package system

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Runs a command in the terminal and streams the output
func RunCommand(command string, displayCommand bool, delay time.Duration) error {
	if displayCommand {
		fmt.Println("Running command: " + command)
	}

	// sleep for X seconds to allow the user to read the command
	time.Sleep(delay * time.Second)

	cmd := exec.Command("/bin/sh", "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("issue running command: %w", err)
	}

	return nil
}

// Runs a command in the terminal and return stdout/stderr as seperate strings
func RunCommandAndCaptureOutput(command string, displayCommand bool, delay time.Duration) (string, error) {
	if displayCommand {
		fmt.Println("Running command: " + command)
	}

	// sleep for X seconds to allow the user to read the command
	time.Sleep(delay * time.Second)

	cmd := exec.Command("/bin/sh", "-c", command)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("issue running command: %w", err)
	}

	return string(out), nil
}
