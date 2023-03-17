package os

import (
	"fmt"
	"github.com/sol-eng/wbi/internal/config"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckLinuxSecurityStatus(osType config.OperatingSystem) (bool, error) {
	if osType == config.Redhat7 || osType == config.Redhat8 {
		fmt.Println("Checking to see if SELinux is active on this server")
		enforceStatus, err := system.RunCommandAndCaptureOutput("getenforce")
		if err != nil {
			return false, fmt.Errorf("issue running getenforce command: %w", err)
		}
		fmt.Println("SELinux Status: " + enforceStatus)
		if strings.Contains(enforceStatus, "Enforcing") {
			return true, nil
		}
	}
	return false, nil
}
