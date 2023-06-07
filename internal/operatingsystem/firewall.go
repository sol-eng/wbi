package operatingsystem

import (
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/system"
)

func CheckFirewallStatus(osType config.OperatingSystem) (bool, error) {
	if osType == config.Redhat7 || osType == config.Redhat8 || osType == config.Redhat9 {
		firewallCheckCommand := "rpm -q firewalld || true"
		rpmOutput, err := system.RunCommandAndCaptureOutput(firewallCheckCommand, false, 0, false)
		if err != nil {
			return false, fmt.Errorf("issue in rpmOutput check with command '%s': %w", firewallCheckCommand, err)
		}

		if strings.Contains(rpmOutput, "not installed") {
			return false, nil
		}

		firewallIsActiveCommand := "systemctl is-active firewalld || true"
		firewallActive, err := system.RunCommandAndCaptureOutput(firewallIsActiveCommand, false, 0, false)
		if err != nil {
			return false, fmt.Errorf("issue in firewallActive with the command '%s': %w", firewallIsActiveCommand, err)
		}

		if !strings.Contains(firewallActive, "inactive") {
			return true, nil
		}

		firewallEnabledCommand := "systemctl is-enabled firewalld || true"
		firewallEnabled, err := system.RunCommandAndCaptureOutput(firewallEnabledCommand, false, 0, false)
		if err != nil {
			return false, fmt.Errorf("issue in firewallEnabled check with the command '%s': %w", firewallEnabledCommand, err)
		}

		if strings.Contains(firewallEnabled, "enabled") {
			return true, nil
		}
	}
	return false, nil
}
