package operatingsystem

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckIfRunningAsRoot() error {
	idOutput, err := system.RunCommandAndCaptureOutput("id -u", false, 0)
	if err != nil {
		return fmt.Errorf("issue running user identification command: %w", err)
	}
	if strings.TrimSpace(idOutput) != "0" {
		return errors.New("wbi must be as root. Please run wbi with sudo and try again")
	}
	return nil
}
