package os

import (
	"fmt"
	"github.com/sol-eng/wbi/internal/config"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckFirewallStatus(osType config.OperatingSystem) (bool, error) {
	if osType == config.Redhat7 || osType == config.Redhat8 {
		stdout, _, err := system.RunCommandAndCaptureOutput("rpm -q firewalld")
		if err != nil {
			return false, fmt.Errorf("issue running firewalld installation command: %w", err)
		}
		if strings.Contains(stdout, "not installed") {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}
