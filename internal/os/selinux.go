package os

import (
	"fmt"
	"github.com/sol-eng/wbi/internal/config"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckSELinuxStatus(osType config.OperatingSystem) (bool, error) {
	if osType == config.Redhat7 || osType == config.Redhat8 {
		stdout, _, err := system.RunCommandAndCaptureOutput("getenforce")
		if err != nil {
			return false, fmt.Errorf("issue running getenforce command: %w", err)
		}
		if strings.Contains(stdout, "Enforcing") || strings.Contains(stdout, "Permissive") {
			return true, nil
		}
	}
	return false, nil
}
