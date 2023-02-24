package os

import (
	"errors"
	"os"
	"runtime"
	"strings"

	"github.com/sol-eng/wbi/internal/config"
)

// Detect which operating system WBI is running on
func DetectOS() (config.OperatingSystem, error) {
	osType := runtime.GOOS
	if osType == "linux" {
		// Check RHEL
		if _, err := os.Stat("/etc/redhat-release"); err == nil {
			releaseVersionRHEL, err := os.ReadFile("/etc/redhat-release")
			if err != nil {
				return config.Unknown, err
			}
			if strings.Contains(string(releaseVersionRHEL), "release 7") {
				return config.Redhat7, nil
			} else if strings.Contains(string(releaseVersionRHEL), "release 8") {
				return config.Redhat8, nil
			} else {
				return config.Unknown, errors.New("unsupported operating system")
			}
		} else if _, err := os.Stat("/etc/issue"); err == nil {
			releaseVersionUbuntu, err := os.ReadFile("/etc/issue")
			if err != nil {
				return config.Unknown, err
			}
			if strings.Contains(string(releaseVersionUbuntu), "Ubuntu 22") {
				return config.Ubuntu22, nil
			} else if strings.Contains(string(releaseVersionUbuntu), "Ubuntu 20") {
				return config.Ubuntu20, nil
			} else if strings.Contains(string(releaseVersionUbuntu), "Ubuntu 18") {
				return config.Ubuntu18, nil
			} else {
				return config.Unknown, errors.New("unsupported operating system")
			}
		} else {
			return config.Unknown, errors.New("unsupported operating system")
		}
	} else {
		return config.Unknown, errors.New("unsupported operating system")
	}
}
