package packagemanager

import (
	"errors"
	"strings"

	"github.com/dpastoor/wbi/internal/config"
)

func BuildPackagemanagerFullURL(url string, repo string, osType config.OperatingSystem) (string, error) {
	var osName string
	switch osType {
	case config.Ubuntu18:
		osName = "bionic"
	case config.Ubuntu20:
		osName = "focal"
	case config.Ubuntu22:
		osName = "jammy"
	case config.Redhat7:
		osName = "centos7"
	case config.Redhat8:
		osName = "centos8"
	default:
		return "", errors.New("operating system not supported")
	}

	fullURL := url + "/" + strings.ToLower(repo) + "/__linux__/" + osName + "/" + "latest"

	return fullURL, nil
}
