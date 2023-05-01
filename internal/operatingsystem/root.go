package operatingsystem

import (
	"errors"
	"fmt"
	"strings"

	cmdlog "github.com/sol-eng/wbi/internal/logging"
	"github.com/sol-eng/wbi/internal/system"
)

func CheckIfRunningAsRoot() error {
	idOutput, err := system.RunCommandAndCaptureOutput("id -u", false, 0, false)
	if err != nil {
		return fmt.Errorf("issue running the user identification command 'id -u': %w", err)
	}

	// save to command log
	cmdlog.Info("if [ \"$(id -u)\" -ne 0 ]; then\n")
	cmdlog.Info("  echo \"This script must be run as root\" 1>&2\n")
	cmdlog.Info("  exit 1\n")
	cmdlog.Info("fi\n")

	if strings.TrimSpace(idOutput) != "0" {
		return errors.New("wbi must be as root, the command 'id -u' did not return 0. Please run wbi with sudo and try again")
	}
	return nil
}
