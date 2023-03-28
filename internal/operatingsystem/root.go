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
		return fmt.Errorf("issue running the user identification command 'id -u': %w", err)
	}
	if strings.TrimSpace(idOutput) != "0" {
		return errors.New("wbi must be as root, the command 'id -u' did not return 0. Please run wbi with sudo and try again")
	}
	return nil
}
