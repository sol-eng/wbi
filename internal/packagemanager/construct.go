package packagemanager

import (
	"errors"

	"github.com/dpastoor/wbi/internal/config"
)

func BuildPackagemanagerFullURL(url string, repo string, osType config.OperatingSystem) (string, error) {

	osName, err := ConvertOSTypeToOSName(osType)
	if err != nil {
		return "", errors.New("there was an issue converting the operating system type to an os name")
	}

	fullURL := url + "/" + repo + "/__linux__/" + osName + "/" + "latest"

	return fullURL, nil
}

func BuildPublicPackageManagerFullURL(osType config.OperatingSystem) (string, error) {

	osName, err := ConvertOSTypeToOSName(osType)
	if err != nil {
		return "", errors.New("there was an issue converting the operating system type to an os name")
	}

	fullURL := "https://packagemanager.rstudio.com/cran/__linux__/" + osName + "/" + "latest"

	return fullURL, nil
}

func ConvertOSTypeToOSName(osType config.OperatingSystem) (string, error) {
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

	return osName, nil
}
