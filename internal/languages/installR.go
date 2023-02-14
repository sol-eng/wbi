package languages

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dpastoor/wbi/internal/system"
)

var availableRVersions = []string{
	"4.2.2", "4.2.1", "4.2.0", "4.1.3", "4.1.2", "4.1.1", "4.1.0", "4.0.5", "4.0.4", "4.0.3", "4.0.2", "4.0.1", "4.0.0", "3.6.3", "3.6.2", "3.6.1", "3.6.0", "3.5.3", "3.5.2", "3.5.1", "3.5.0", "3.4.4", "3.4.3", "3.4.2", "3.4.1", "3.4.0", "3.3.3", "3.3.2", "3.3.1", "3.3.0",
}

const Ubuntu22 = "ubuntu22"
const Ubuntu20 = "ubuntu20"
const Ubuntu18 = "ubuntu18"

const Redhat8 = "redhat8"
const Redhat7 = "redhat7"

// InstallerInfo contains the information needed to download and install R
type InstallerInfo struct {
	Name    string
	URL     string
	Version string
}

// Prompt users if they would like to install R versions
func RInstallPrompt() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to install version(s) of R?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the R install prompt")
	}
	return name, nil
}

func RetrieveValidRVersions() ([]string, error) {
	// TODO make this dynamic based on https://cran.r-project.org/src/base/R-4/ and https://cran.r-project.org/src/base/R-3/
	return availableRVersions, nil
}

// Prompt asking users which R version(s) they would like to install
func RSelectVersionsPrompt(availableRVersions []string) ([]string, error) {
	var qs = []*survey.Question{
		{
			Name: "rversions",
			Prompt: &survey.MultiSelect{
				Message: "Which version(s) of R would you like to install?",
				Options: availableRVersions,
				Default: availableRVersions[0],
			},
		},
	}
	rVersionsAnswers := struct {
		RVersions []string `survey:"rversions"`
	}{}
	err := survey.Ask(qs, &rVersionsAnswers)
	if err != nil {
		return []string{}, errors.New("there was an issue with the R versions selection prompt")
	}
	if len(rVersionsAnswers.RVersions) == 0 {
		return []string{}, errors.New("At least one R version must be selected")
	}
	return rVersionsAnswers.RVersions, nil
}

// Downloads the R installer, and installs R
func DownloadAndInstallR(rVersion string, osType string) error {
	// Create InstallerInfo with the proper information
	installerInfo, err := PopulateInstallerInfo(rVersion, osType)
	if err != nil {
		return fmt.Errorf("PopulateInstallerInfo: %w", err)
	}
	// Download installer
	filepath, err := installerInfo.DownloadR()
	if err != nil {
		return fmt.Errorf("DownloadR: %w", err)
	}
	// Install R
	err = InstallR(filepath, osType, rVersion)
	if err != nil {
		return fmt.Errorf("InstallR: %w", err)
	}
	return nil
}

// Create a temporary file and download the R installer to it.
func (installerInfo *InstallerInfo) DownloadR() (string, error) {

	url := installerInfo.URL
	name := installerInfo.Name

	fmt.Println("Downloading R installer from: " + url)

	// Create the file
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return tmpFile.Name(), err
	}
	defer tmpFile.Close()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, url, nil)
	if err != nil {
		return "", errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return "", errors.New("error downloading R installer")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("error retrieving R installer")
	}

	// Writer the body to file
	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// Installs R in a certain way based on the operating system
func InstallR(filepath string, osType string, rVersion string) error {
	// Update apt and install gdebi-core if an Ubuntu system
	if osType == Ubuntu22 || osType == Ubuntu20 || osType == Ubuntu18 {
		AptErr := UpgradeApt()
		if AptErr != nil {
			return fmt.Errorf("UpgradeApt: %w", AptErr)
		}

		GdebiCoreErr := InstallGdebiCore()
		if GdebiCoreErr != nil {
			return fmt.Errorf("InstallGdebiCore: %w", GdebiCoreErr)
		}
	} else if osType == Redhat8 {
		// Enable the Extra Packages for Enterprise Linux (EPEL) repository
		EnableEPELErr := EnableEPELRepo(osType)
		if EnableEPELErr != nil {
			return fmt.Errorf("EnableEPELRepo: %w", EnableEPELErr)
		}
		// Enable the CodeReady Linux Builder repository
		EnableCodeReadyErr := EnableCodeReadyRepo()
		if EnableCodeReadyErr != nil {
			return fmt.Errorf("EnableCodeReadyRepo: %w", EnableCodeReadyErr)
		}
	} else if osType == Redhat7 {
		// Enable the Extra Packages for Enterprise Linux (EPEL) repository
		EnableEPELErr := EnableEPELRepo(osType)
		if EnableEPELErr != nil {
			return fmt.Errorf("EnableEPELRepo: %w", EnableEPELErr)
		}
		//Enable the Optional repository
		EnableOptionalRepoErr := EnableOptionalRepo()
		if EnableOptionalRepoErr != nil {
			return fmt.Errorf("EnableOptionalRepo: %w", EnableOptionalRepoErr)
		}
	} else {
		return errors.New("unsupported operating system")
	}

	installCommand, err := RetrieveInstallCommandForR(filepath, osType)
	if err != nil {
		return fmt.Errorf("RetrieveInstallCommandForWorkbench: %w", err)
	}

	err = system.RunCommand(installCommand)
	if err != nil {
		return fmt.Errorf("issue installing R: %w", err)
	}

	successMessage := "\nR version " + rVersion + " successfully installed!\n"
	fmt.Println(successMessage)
	return nil
}

// Creates the proper command to install R based on the operating system
func RetrieveInstallCommandForR(filepath string, os string) (string, error) {
	switch os {
	case Ubuntu22, Ubuntu20, Ubuntu18:
		return "sudo gdebi -n " + filepath, nil
	case Redhat7, Redhat8:
		return "sudo yum install -y " + filepath, nil
	default:
		return "", errors.New("operating system not supported")
	}
}

func PopulateInstallerInfo(rVersion string, osType string) (InstallerInfo, error) {
	switch osType {
	case Ubuntu18:
		return InstallerInfo{
			Name:    "r-" + rVersion + "_1_amd64.deb",
			URL:     "https://cdn.rstudio.com/r/ubuntu-1804/pkgs/r-" + rVersion + "_1_amd64.deb",
			Version: rVersion,
		}, nil
	case Ubuntu20:
		return InstallerInfo{
			Name:    "r-" + rVersion + "_1_amd64.deb",
			URL:     "https://cdn.rstudio.com/r/ubuntu-2004/pkgs/r-" + rVersion + "_1_amd64.deb",
			Version: rVersion,
		}, nil
	case Ubuntu22:
		return InstallerInfo{
			Name:    "r-" + rVersion + "_1_amd64.deb",
			URL:     "https://cdn.rstudio.com/r/ubuntu-2204/pkgs/r-" + rVersion + "_1_amd64.deb",
			Version: rVersion,
		}, nil
	case Redhat7:
		return InstallerInfo{
			Name:    "R-" + rVersion + "-1-1.x86_64.rpm",
			URL:     "https://cdn.rstudio.com/r/centos-7/pkgs/R-" + rVersion + "-1-1.x86_64.rpm",
			Version: rVersion,
		}, nil
	case Redhat8:
		return InstallerInfo{
			Name:    "R-" + rVersion + "-1-1.x86_64.rpm",
			URL:     "https://cdn.rstudio.com/r/centos-8/pkgs/R-" + rVersion + "-1-1.x86_64.rpm",
			Version: rVersion,
		}, nil
	default:
		return InstallerInfo{}, errors.New("operating system not supported")
	}
}

// Installs Gdebi Core
func InstallGdebiCore() error {
	gdebiCoreCommand := "sudo apt-get install -y gdebi-core"
	err := system.RunCommand(gdebiCoreCommand)
	if err != nil {
		return fmt.Errorf("issue installing gdebi-core: %w", err)
	}

	fmt.Println("\ngdebi-core has been successfully installed!\n")
	return nil
}

// Upgrades Apt
func UpgradeApt() error {
	aptUpgradeCommand := "sudo apt-get update"
	err := system.RunCommand(aptUpgradeCommand)
	if err != nil {
		return fmt.Errorf("issue upgrading apt: %w", err)
	}

	fmt.Println("\napt has been successfully upgraded!\n")
	return nil
}

// Enable the Extra Packages for Enterprise Linux (EPEL) repository
func EnableEPELRepo(osType string) error {
	var EPELCommand string
	if osType == Redhat8 {
		EPELCommand = "sudo yum install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm"
	} else if osType == Redhat7 {
		EPELCommand = "sudo yum install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm "
	} else {
		return errors.New("operating system not supported")
	}

	err := system.RunCommand(EPELCommand)
	if err != nil {
		return fmt.Errorf("issue enabling EPEL repo: %w", err)
	}

	fmt.Println("\nThe Extra Packages for Enterprise Linux (EPEL) repository has been successfully enabled!\n")
	return nil
}

// Enable the CodeReady Linux Builder repository:
func EnableCodeReadyRepo() error {
	// TODO add support for On Premise as well as cloud (currently only cloud)
	dnfPluginsCoreCommand := "sudo dnf install dnf-plugins-core"
	err := system.RunCommand(dnfPluginsCoreCommand)
	if err != nil {
		return fmt.Errorf("issue installing dnf-plugins-core: %w", err)
	}

	enableCodeReadyCommand := `sudo dnf config-manager --set-enabled "codeready-builder-for-rhel-8-*-rpms"`
	err = system.RunCommand(enableCodeReadyCommand)
	if err != nil {
		return fmt.Errorf("issue enabling the CodeReady Linux Builder repo: %w", err)
	}

	fmt.Println("\nThe CodeReady Linux Builder repository has been successfully enabled!\n")
	return nil
}

// Enable the Optional repository
func EnableOptionalRepo() error {
	// TODO add support for On Premise as well as cloud (currently only cloud)
	yumUtilsCommand := "sudo yum install yum-utils"
	err := system.RunCommand(yumUtilsCommand)
	if err != nil {
		return fmt.Errorf("issue installing yum-utils: %w", err)
	}

	enableOptionalCommand := `sudo yum-config-manager --enable "rhel-*-optional-rpms"`
	err = system.RunCommand(enableOptionalCommand)
	if err != nil {
		return fmt.Errorf("issue enabling the optional repo: %w", err)
	}

	fmt.Println("\nThe Optional repository has been successfully enabled!\n")
	return nil
}
