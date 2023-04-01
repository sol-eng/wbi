package prodrivers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/install"
	"github.com/sol-eng/wbi/internal/system"
)

// InstallerInfo contains the information needed to download and install Posit Pro Drivers
type InstallerInfo struct {
	BaseName string `json:"basename"`
	URL      string `json:"url"`
	Version  string `json:"version"`
	Label    string `json:"label"`
}

// OperatingSystems contains the installer information for each supported operating system
type OperatingSystems struct {
	// Posit Pro Drivers are the same for all Ubuntu versions so we only need one
	Focal   InstallerInfo `json:"focal"`
	Redhat7 InstallerInfo `json:"redhat7_64"`
	Redhat8 InstallerInfo `json:"rhel8"`
}

// Installer contains the installer information for a product
type Installer struct {
	Installer OperatingSystems `json:"installer"`
}

// ProDrivers contains product information
type ProDrivers struct {
	ProDrivers Installer `json:"pro_drivers"`
}

// Retrieves JSON data from Posit, downloads the Pro Drivers installer, and installs Pro Drivers
func DownloadAndInstallProDrivers(osType config.OperatingSystem) error {
	// Retrieve JSON data
	rstudio, err := RetrieveProDriversInstallerInfo()
	if err != nil {
		return fmt.Errorf("RetrieveProDriversInstallerInfo: %w", err)
	}
	// Retrieve installer info
	installerInfo, err := rstudio.GetInstallerInfo(osType)
	if err != nil {
		return fmt.Errorf("GetInstallerInfo: %w", err)
	}
	// Install prerequisites
	err = InstallUnixODBC(osType)
	if err != nil {
		return fmt.Errorf("InstallUnixODBC: %w", err)
	}
	// Download installer
	filepath, err := install.DownloadFile("Pro Drivers", installerInfo.URL, installerInfo.BaseName)
	if err != nil {
		return fmt.Errorf("DownloadFile: %w", err)
	}
	// Install Pro Drivers
	err = InstallProDrivers(filepath, osType)
	if err != nil {
		return fmt.Errorf("InstallProDrivers: %w", err)
	}
	// Configure ODBC driver name and locations
	err = BackupAndAppendODBCConfiguration()
	if err != nil {
		return fmt.Errorf("BackupAndAppendODBCConfiguration: %w", err)
	}
	fmt.Println("\nPosit Pro Drivers next steps:\nNow that the Pro Drivers are installed and /etc/odbcinst.ini is setup, the next step is to test database connectivity and/or create DSNs in your /etc/odbc.ini file.\n\n More information about each of these steps can be found here: https://docs.posit.co/pro-drivers/workbench-connect/#step-4-testing-database-connectivity")
	return nil
}

// Installs Posit Pro Drivers in a certain way based on the operating system
func InstallProDrivers(filepath string, osType config.OperatingSystem) error {
	installCommand, err := install.RetrieveInstallCommand(filepath, osType)
	if err != nil {
		return fmt.Errorf("RetrieveInstallCommand: %w", err)
	}

	err = system.RunCommand(installCommand)
	if err != nil {
		return fmt.Errorf("issue installing Pro Drivers: %w", err)
	}

	fmt.Println("\nPosit Pro Drivers have been successfully installed!")
	return nil
}

// Pulls out the installer information from the JSON data based on the operating system
func (pd *ProDrivers) GetInstallerInfo(osType config.OperatingSystem) (InstallerInfo, error) {
	switch osType {
	// Posit Pro Drivers are the same for all Ubuntu versions
	case config.Ubuntu18, config.Ubuntu20, config.Ubuntu22:
		return pd.ProDrivers.Installer.Focal, nil
	case config.Redhat7:
		return pd.ProDrivers.Installer.Redhat7, nil
	// Posit Pro Drivers are the same for RHEL 8 and RHEL 9
	case config.Redhat8, config.Redhat9:
		return pd.ProDrivers.Installer.Redhat8, nil
	default:
		return InstallerInfo{}, errors.New("operating system not supported")
	}
}

// Retrieves JSON data from Posit
func RetrieveProDriversInstallerInfo() (ProDrivers, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, "https://www.rstudio.com/wp-content/downloads.json", nil)
	if err != nil {
		return ProDrivers{}, errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return ProDrivers{}, errors.New("error retrieving JSON data")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return ProDrivers{}, errors.New("error retrieving JSON data")
	}
	var proDrivers ProDrivers
	err = json.NewDecoder(res.Body).Decode(&proDrivers)
	if err != nil {
		return ProDrivers{}, errors.New("error unmarshalling JSON data")
	}
	return proDrivers, nil
}

// Installs unixODBC and unixODBC-devel
func InstallUnixODBC(osType config.OperatingSystem) error {
	if osType == config.Ubuntu22 || osType == config.Ubuntu20 || osType == config.Ubuntu18 {
		err := system.RunCommand("apt-get -y install unixodbc unixodbc-dev")
		if err != nil {
			return fmt.Errorf("issue installing unixodbc and unixodbc-dev: %w", err)
		}
	} else if osType == config.Redhat7 || osType == config.Redhat8 || osType == config.Redhat9 {
		err := system.RunCommand("yum -y install unixODBC unixODBC-devel")
		if err != nil {
			return fmt.Errorf("issue installing unixodbc and unixodbc-dev: %w", err)
		}
	} else {
		return errors.New("operating system not supported")
	}

	fmt.Println("\nunixodbc and unixodbc-dev has been successfully installed!")
	return nil
}

func BackupAndAppendODBCConfiguration() error {
	// backup odbcinst.ini if one already exists
	if _, err := os.Stat("/etc/odbcinst.ini"); err == nil {
		fmt.Println("Backing up /etc/odbcinst.ini to /etc/odbcinst.ini.bak")
		err := system.RunCommand("cp /etc/odbcinst.ini /etc/odbcinst.ini.bak")
		if err != nil {
			return fmt.Errorf("issue backing up /etc/odbcinst.ini: %w", err)
		}
	}
	// append sample ODBC configuration to odbcinst.ini
	_, err := system.RunCommandAndCaptureOutput("cat /opt/rstudio-drivers/odbcinst.ini.sample | tee -a /etc/odbcinst.ini")
	if err != nil {
		return fmt.Errorf("issue appending sample configuration to /etc/odbcinst.ini: %w", err)
	}

	fmt.Println("\nThe sample preconfigured odbcinst.ini has been appended to /etc/odbcinst.ini")
	return nil
}
