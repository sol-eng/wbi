package languages

import (
	"errors"


	"github.com/sol-eng/wbi/internal/config"

)

// InstallerInfo contains the information needed to download and install R and Python
type InstallerInfo struct {
	Name    string
	URL     string
	Version string
}

func PopulateInstallerInfo(language string, version string, osType config.OperatingSystem) (InstallerInfo, error) {
	switch osType {
	case config.Ubuntu18:
		return InstallerInfo{
			Name:    language + "-" + version + "_1_amd64.deb",
			URL:     "https://cdn.rstudio.com/" + language + "/ubuntu-1804/pkgs/" + language + "-" + version + "_1_amd64.deb",
			Version: version,
		}, nil
	case config.Ubuntu20:
		return InstallerInfo{
			Name:    language + "-" + version + "_1_amd64.deb",
			URL:     "https://cdn.rstudio.com/" + language + "/ubuntu-2004/pkgs/" + language + "-" + version + "_1_amd64.deb",
			Version: version,
		}, nil
	case config.Ubuntu22:
		return InstallerInfo{
			Name:    language + "-" + version + "_1_amd64.deb",
			URL:     "https://cdn.rstudio.com/" + language + "/ubuntu-2204/pkgs/" + language + "-" + version + "_1_amd64.deb",
			Version: version,
		}, nil
	case config.Redhat7: //TODO: Upper the R language variable for RHEL
		if language == "r" {
			language = strings.ToUpper(language)
		}

		return InstallerInfo{
			Name:    language + "-" + version + "-1-1.x86_64.rpm",
			URL:     "https://cdn.rstudio.com/" + language + "/centos-7/pkgs/" + language + "-" + version + "-1-1.x86_64.rpm",
			Version: version,
		}, nil
	case config.Redhat8:
		return InstallerInfo{
			Name:    language + "-" + version + "-1-1.x86_64.rpm",
			URL:     "https://cdn.rstudio.com/" + language + "/centos-8/pkgs/" + language + "-" + version + "-1-1.x86_64.rpm",
			Version: version,
		}, nil
	default:
		return InstallerInfo{}, errors.New("operating system not supported")
	}
}
