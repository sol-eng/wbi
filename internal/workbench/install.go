package workbench

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dpastoor/wbi/internal/system"
)

// InstallerInfo contains the information needed to download and install Workbench
type InstallerInfo struct {
	BaseName string `json:"basename"`
	URL      string `json:"url"`
	Version  string `json:"version"`
	Label    string `json:"label"`
}

// OperatingSystems contains the installer information for each supported operating system
type OperatingSystems struct {
	Bionic  InstallerInfo `json:"bionic"`
	Jammy   InstallerInfo `json:"jammy"`
	Redhat7 InstallerInfo `json:"redhat7_64"`
	Redhat8 InstallerInfo `json:"rhel8"`
}

// Installer contains the installer information for a product
type Installer struct {
	Installer OperatingSystems `json:"installer"`
}

// ProductType contains the installer for each product type
type ProductType struct {
	Server Installer `json:"server"`
}

// Category contains information for stable and preview product types
type Category struct {
	Stable ProductType `json:"stable"`
}

// Product contains information for each RStudio product
type Product struct {
	Pro Category `json:"pro"`
}

// RStudio contains product information
type RStudio struct {
	Rstudio Product `json:"rstudio"`
}

// Retrieves JSON data from Posit, downloads the Workbench installer, and installs Workbench
func DownloadAndInstallWorkbench(osType string) error {
	// Retrieve JSON data
	rstudio, err := RetrieveWorkbenchInstallerInfo()
	if err != nil {
		return fmt.Errorf("RetrieveWorkbenchInstallerInfo: %w", err)
	}
	// Retrieve installer info
	installerInfo, err := rstudio.GetInstallerInfo(osType)
	if err != nil {
		return fmt.Errorf("GetInstallerInfo: %w", err)
	}
	// Download installer
	filepath, err := installerInfo.DownloadWorkbench()
	if err != nil {
		return fmt.Errorf("DownloadWorkbench: %w", err)
	}
	// Install Workbench
	err = InstallWorkbench(filepath, osType)
	if err != nil {
		return fmt.Errorf("InstallWorkbench: %w", err)
	}
	return nil
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

// Installs Workbench in a certain way based on the operating system
func InstallWorkbench(filepath string, osType string) error {
	// Install gdebi-core if an Ubuntu system
	if osType == "ubuntu22" || osType == "ubuntu20" || osType == "ubuntu18" {
		AptErr := UpgradeApt()
		if AptErr != nil {
			return fmt.Errorf("UpgradeApt: %w", AptErr)
		}

		GdebiCoreErr := InstallGdebiCore()
		if GdebiCoreErr != nil {
			return fmt.Errorf("InstallGdebiCore: %w", GdebiCoreErr)
		}
	}

	installCommand, err := RetrieveInstallCommandForWorkbench(filepath, osType)
	if err != nil {
		return fmt.Errorf("RetrieveInstallCommandForWorkbench: %w", err)
	}

	err = system.RunCommand(installCommand)
	if err != nil {
		return fmt.Errorf("issue installing Workbench: %w", err)
	}

	fmt.Println("\nWorkbench has been successfully installed!\n")
	return nil
}

// Creates the proper command to install Workbench based on the operating system
func RetrieveInstallCommandForWorkbench(filepath string, os string) (string, error) {
	switch os {
	case "ubuntu22", "ubuntu20", "ubuntu18":
		return "sudo gdebi -n " + filepath, nil
	case "redhat7", "redhat8":
		return "sudo yum install -y " + filepath, nil
	default:
		return "", errors.New("operating system not supported")
	}
}

// Pulls out the installer information from the JSON data based on the operating system
func (r *RStudio) GetInstallerInfo(os string) (InstallerInfo, error) {
	switch os {
	case "ubuntu18", "ubuntu20":
		return r.Rstudio.Pro.Stable.Server.Installer.Bionic, nil
	case "ubuntu22":
		return r.Rstudio.Pro.Stable.Server.Installer.Jammy, nil
	case "redhat7":
		return r.Rstudio.Pro.Stable.Server.Installer.Redhat7, nil
	case "redhat8":
		return r.Rstudio.Pro.Stable.Server.Installer.Redhat8, nil
	default:
		return InstallerInfo{}, errors.New("operating system not supported")
	}
}

// Retrieves JSON data from Posit
func RetrieveWorkbenchInstallerInfo() (RStudio, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, "https://www.rstudio.com/wp-content/downloads.json", nil)
	if err != nil {
		return RStudio{}, errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return RStudio{}, errors.New("error retrieving JSON data")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return RStudio{}, errors.New("error retrieving JSON data")
	}
	var rstudio RStudio
	err = json.NewDecoder(res.Body).Decode(&rstudio)
	if err != nil {
		return RStudio{}, errors.New("error unmarshalling JSON data")
	}
	return rstudio, nil
}
