package os

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

func DetectOS() (string, error) {
	osType := runtime.GOOS
	if osType == "linux" {
		// Check RHEL
		if _, err := os.Stat("/etc/redhat-release"); err == nil {
			releaseVersionRHEL, err := os.ReadFile("/etc/redhat-release")
			if err != nil {
				return "", err
			}
			if strings.Contains(string(releaseVersionRHEL), "release 7") {
				osType = "rhel7"
			}
			if strings.Contains(string(releaseVersionRHEL), "release 8") {
				osType = "rhel8"
			}
			if strings.Contains(string(releaseVersionRHEL), "release 9") {
				osType = "rhel9"
			}
			return osType, nil
		} else if _, err := os.Stat("/etc/issue"); err == nil {
			releaseVersionUbuntu, err := os.ReadFile("/etc/issue")
			if err != nil {
				return "", err
			}
			if strings.Contains(string(releaseVersionUbuntu), "Ubuntu 22") {
				osType = "ubuntu22"
			}
			if strings.Contains(string(releaseVersionUbuntu), "Ubuntu 20") {
				osType = "ubuntu20"
			}
			if strings.Contains(string(releaseVersionUbuntu), "Ubuntu 18") {
				osType = "ubuntu18"
			}
			return osType, nil
		} else {
			return "", errors.New("unsupported operating system")
		}
	} else {
		return "", errors.New("unsupported operating system")
	}
}
