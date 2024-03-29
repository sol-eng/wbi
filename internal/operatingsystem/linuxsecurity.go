package operatingsystem

import (
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/config"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckLinuxSecurityStatus(osType config.OperatingSystem) (bool, error) {
	if osType == config.Redhat7 || osType == config.Redhat8 || osType == config.Redhat9 {
		system.PrintAndLogInfo("Checking to see if SELinux is active on this server")
		enforceStatus, err := system.RunCommandAndCaptureOutput("getenforce", false, 0, false)
		if err != nil {
			return false, fmt.Errorf("issue running the getenforce command: %w", err)
		}
		system.PrintAndLogInfo("SELinux Status: " + enforceStatus)
		if strings.Contains(enforceStatus, "Enforcing") {
			return true, nil
		}
	}
	return false, nil
}
