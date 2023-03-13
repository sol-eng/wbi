package os

import (
	"fmt"
	"github.com/sol-eng/wbi/internal/config"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckFirewallStatus(osType config.OperatingSystem) (bool, error) {
	if osType == config.Redhat7 || osType == config.Redhat8 {
		stdout, stderr, _ := system.RunCommandAndCaptureOutput("rpm -q firewalld")
		fmt.Println("This is stdout: " + stdout)
		fmt.Println("This is stderr: " + stderr)

		if strings.Contains(stderr, "not installed") {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}
