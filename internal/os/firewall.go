package os

import (
	"github.com/sol-eng/wbi/internal/config"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckFirewallStatus(osType config.OperatingSystem) (bool, error) {
	if osType == config.Redhat7 || osType == config.Redhat8 {
		rpmOutput, _, _ := system.RunCommandAndCaptureOutput("rpm -q firewalld")

		if strings.Contains(rpmOutput, "not installed") {
			return false, nil
		}

		firewallActive, _, _ := system.RunCommandAndCaptureOutput("systemctl is-active firewalld")

		if strings.Contains(firewallActive, "active") {
			return true, nil
		}

		firewallEnabled, _, _ := system.RunCommandAndCaptureOutput("systemctl is-active firewalld")

		if strings.Contains(firewallEnabled, "enabled") {
			return true, nil
		}
	}
	return false, nil
}
