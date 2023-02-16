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

var availablePythonVersions = []string{
	"3.11.2",
	"3.11.1",
	"3.11.0",
	"3.10.10",
	"3.10.9",
	"3.10.8",
	"3.10.7",
	"3.10.6",
	"3.10.5",
	"3.10.4",
	"3.10.3",
	"3.10.2",
	"3.10.1",
	"3.10.0",
	"3.9.16",
	"3.9.15",
	"3.9.14",
	"3.9.13",
	"3.9.12",
	"3.9.11",
	"3.9.10",
	"3.9.9",
	"3.9.8",
	"3.9.7",
	"3.9.6",
	"3.9.5",
	"3.9.4",
	"3.9.3",
	"3.9.2",
	"3.9.1",
	"3.9.0",
	"3.8.16",
	"3.8.15",
	"3.8.14",
	"3.8.13",
	"3.8.12",
	"3.8.11",
	"3.8.10",
	"3.8.9",
	"3.8.8",
	"3.8.7",
	"3.8.6",
	"3.8.5",
	"3.8.4",
	"3.8.3",
	"3.8.2",
	"3.8.1",
	"3.8.0",
	"3.7.16",
	"3.7.15",
	"3.7.14",
	"3.7.13",
	"3.7.12",
	"3.7.11",
	"3.7.10",
	"3.7.9",
	"3.7.8",
	"3.7.7",
	"3.7.6",
	"3.7.5",
	"3.7.4",
	"3.7.3",
	"3.7.2",
	"3.7.1",
	"3.7.0",
}

// InstallerInfoPython contains the information needed to download and install Python
type InstallerInfoPython struct {
	Name    string
	URL     string
	Version string
}

// Prompt users if they would like to install Python versions
func PythonInstallPrompt() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to install version(s) of Python?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Python install prompt")
	}
	return name, nil
}

func RetrieveValidPythonVersions() ([]string, error) {
	// TODO make this dynamic based on https://cdn.posit.co/python/versions.json
	return availablePythonVersions, nil
}

// Prompt asking users which Python version(s) they would like to install
func PythonSelectVersionsPrompt(availablePythonVersions []string) ([]string, error) {
	var qs = []*survey.Question{
		{
			Name: "pythonVersions",
			Prompt: &survey.MultiSelect{
				Message: "Which version(s) of Python would you like to install?",
				Options: availablePythonVersions,
				Default: availablePythonVersions[0],
			},
		},
	}
	pythonVersionsAnswers := struct {
		PythonVersions []string `survey:"pythonVersions"`
	}{}
	err := survey.Ask(qs, &pythonVersionsAnswers)
	if err != nil {
		return []string{}, errors.New("there was an issue with the Python versions selection prompt")
	}
	if len(pythonVersionsAnswers.PythonVersions) == 0 {
		return []string{}, errors.New("At least one Python version must be selected")
	}
	return pythonVersionsAnswers.PythonVersions, nil
}

// Downloads the Python installer, and installs Python
func DownloadAndInstallPython(pythonVersion string, osType string) error {
	// Create InstallerInfoPython with the proper information
	installerInfo, err := PopulateInstallerInfoPython(pythonVersion, osType)
	if err != nil {
		return fmt.Errorf("PopulateInstallerInfoPython: %w", err)
	}
	// Download installer
	filepath, err := installerInfo.DownloadPython()
	if err != nil {
		return fmt.Errorf("DownloadPython: %w", err)
	}
	// Install Python
	err = InstallPython(filepath, osType, pythonVersion)
	if err != nil {
		return fmt.Errorf("InstallPython: %w", err)
	}
	// Upgrade pip, setuptools, and wheel
	err = UpgradePythonTools(pythonVersion)
	if err != nil {
		return fmt.Errorf("UpgradePythonTools: %w", err)
	}

	return nil
}

// Create a temporary file and download the Python installer to it.
func (installerInfoPython *InstallerInfoPython) DownloadPython() (string, error) {

	url := installerInfoPython.URL
	name := installerInfoPython.Name

	fmt.Println("Downloading Python installer from: " + url)

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
		return "", errors.New("error downloading Python installer")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("error retrieving Python installer")
	}

	// Writer the body to file
	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// Installs Python in a certain way based on the operating system
func InstallPython(filepath string, osType string, pythonVersion string) error {
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

	installCommand, err := RetrieveInstallCommandForPython(filepath, osType)
	if err != nil {
		return fmt.Errorf("RetrieveInstallCommandForPython: %w", err)
	}

	err = system.RunCommand(installCommand)
	if err != nil {
		return fmt.Errorf("issue installing Python: %w", err)
	}

	successMessage := "\nPython version " + pythonVersion + " successfully installed!\n"
	fmt.Println(successMessage)
	return nil
}

func UpgradePythonTools(pythonVersion string) error {
	upgradeCommand := "/opt/python/" + pythonVersion + "/bin/pip install --upgrade pip setuptools wheel"
	err := system.RunCommand(upgradeCommand)
	if err != nil {
		return fmt.Errorf("issue upgrading pip, setuptools and wheel for Python: %w", err)
	}

	successMessage := "\npip, setuptools and wheel have been upgraded for Python version " + pythonVersion + "\n"
	fmt.Println(successMessage)

	return nil
}

// Creates the proper command to install Python based on the operating system
func RetrieveInstallCommandForPython(filepath string, os string) (string, error) {
	switch os {
	case Ubuntu22, Ubuntu20, Ubuntu18:
		return "sudo gdebi -n " + filepath, nil
	case Redhat7, Redhat8:
		return "sudo yum install -y " + filepath, nil
	default:
		return "", errors.New("operating system not supported")
	}
}

func PopulateInstallerInfoPython(pythonVersion string, osType string) (InstallerInfoPython, error) {
	switch osType {
	case Ubuntu18:
		return InstallerInfoPython{
			Name:    "python-" + pythonVersion + "_1_amd64.deb",
			URL:     "https://cdn.rstudio.com/python/ubuntu-1804/pkgs/python-" + pythonVersion + "_1_amd64.deb",
			Version: pythonVersion,
		}, nil
	case Ubuntu20:
		return InstallerInfoPython{
			Name:    "python-" + pythonVersion + "_1_amd64.deb",
			URL:     "https://cdn.rstudio.com/python/ubuntu-2004/pkgs/python-" + pythonVersion + "_1_amd64.deb",
			Version: pythonVersion,
		}, nil
	case Ubuntu22:
		return InstallerInfoPython{
			Name:    "python-" + pythonVersion + "_1_amd64.deb",
			URL:     "https://cdn.rstudio.com/python/ubuntu-2204/pkgs/python-" + pythonVersion + "_1_amd64.deb",
			Version: pythonVersion,
		}, nil
	case Redhat7:
		return InstallerInfoPython{
			Name:    "python-" + pythonVersion + "-1-1.x86_64.rpm",
			URL:     "https://cdn.rstudio.com/python/centos-7/pkgs/python-" + pythonVersion + "-1-1.x86_64.rpm",
			Version: pythonVersion,
		}, nil
	case Redhat8:
		return InstallerInfoPython{
			Name:    "python-" + pythonVersion + "-1-1.x86_64.rpm",
			URL:     "https://cdn.rstudio.com/python/centos-8/pkgs/python-" + pythonVersion + "-1-1.x86_64.rpm",
			Version: pythonVersion,
		}, nil
	default:
		return InstallerInfoPython{}, errors.New("operating system not supported")
	}
}
