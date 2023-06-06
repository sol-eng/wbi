package system

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
	cmdlog "github.com/sol-eng/wbi/internal/logging"
)

// Runs a command in the terminal and streams the output
func RunCommand(command string, displayCommand bool, delay time.Duration, save bool) error {
	if displayCommand {
		PrintAndLogInfo("Running command: " + command)
	}

	// sleep for X seconds to allow the user to read the command
	time.Sleep(delay * time.Second)

	var errBuf, outBuf bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", command)

	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("issue running the command '%s': %w", command, err)
	}
	if save {
		cmdlog.Info(command)
	}
	if len(outBuf.String()) > 0 {
		log.Info(outBuf.String())
	}
	if len(errBuf.String()) > 0 {
		log.Error(errBuf.String())
	}

	return nil
}

// Runs a command in the terminal and return stdout/stderr as seperate strings
func RunCommandAndCaptureOutput(command string, displayCommand bool, delay time.Duration, save bool) (string, error) {
	if displayCommand {
		PrintAndLogInfo("Running command: " + command)
	}

	// sleep for X seconds to allow the user to read the command
	time.Sleep(delay * time.Second)

	cmd := exec.Command("/bin/sh", "-c", command)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("issue running the command '%s': %w", command, err)
	}
	if save {
		cmdlog.Info(command)
	}
	log.Info(string(out))

	return string(out), nil
}
