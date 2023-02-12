package system

import (
	"bufio"
	"fmt"
	"os/exec"
)

// Runs a command in the terminal and streams the output
func RunCommand(command string) error {
	fmt.Println("Running command: " + command)
	cmd := exec.Command("/bin/sh", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("issue running command: %w", err)
	}
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return fmt.Errorf("issue running command: %w", scanner.Err())
	}
	cmd.Wait()
	return nil
}
