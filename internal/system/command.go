package system

import (
	"bytes"
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

// Runs a command in the terminal and return stdout/stderr as seperate strings
func RunCommandAndCaptureOutput(command string) (string, string, error) {
	fmt.Println("Running command for output: " + command)

	testcommand := "/bin/sh -c " + command
	cmd := exec.Command(testcommand)

	var errb bytes.Buffer
	//cmd.Stdout = &outb
	cmd.Stderr = &errb
	output, err := cmd.Output()
	fmt.Println(string(output)) // when success

	if err != nil {
		return "", "", fmt.Errorf("issue running command: %w", err)
	}

	return string(output), errb.String(), nil
}
