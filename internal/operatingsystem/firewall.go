package operatingsystem

import (
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/config"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckFirewallStatus(osType config.OperatingSystem) (bool, error) {

	if osType == config.Redhat7 || osType == config.Redhat8 || osType == config.Redhat9 {
		rpmOutput, err := system.RunCommandAndCaptureOutput("rpm -q firewalld || true")
		if err != nil {
			return false, fmt.Errorf("issue in rpmOutput check: %w", err)
		}

		if strings.Contains(rpmOutput, "not installed") {
			return false, nil
		}

		firewallActive, err := system.RunCommandAndCaptureOutput("systemctl is-active firewalld || true")
		if err != nil {
			return false, fmt.Errorf("issue in firewallActive: %w", err)
		}

		if !strings.Contains(firewallActive, "inactive") {
			return true, nil
		}

		firewallEnabled, err := system.RunCommandAndCaptureOutput("systemctl is-enabled firewalld || true")
		if err != nil {
			return false, fmt.Errorf("issue in firewallEnabled check: %w", err)
		}

		if strings.Contains(firewallEnabled, "enabled") {
			return true, nil
		}
	}
	return false, nil
}
