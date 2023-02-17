package install

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dpastoor/wbi/internal/config"
	"github.com/dpastoor/wbi/internal/system"
)

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

// Enable the Extra Packages for Enterprise Linux (EPEL) repository
func EnableEPELRepo(osType config.OperatingSystem) error {
	var EPELCommand string
	if osType == config.Redhat8 {
		EPELCommand = "sudo yum install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm"
	} else if osType == config.Redhat7 {
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

// Installs R/Python in a certain way based on the operating system
func InstallLanguage(language string, filepath string, osType config.OperatingSystem, version string) error {
	languageTitleCase := strings.Title(language)

	// Update apt and install gdebi-core if an Ubuntu system
	if osType == config.Ubuntu22 || osType == config.Ubuntu20 || osType == config.Ubuntu18 {
		AptErr := UpgradeApt()
		if AptErr != nil {
			return fmt.Errorf("UpgradeApt: %w", AptErr)
		}

		GdebiCoreErr := InstallGdebiCore()
		if GdebiCoreErr != nil {
			return fmt.Errorf("InstallGdebiCore: %w", GdebiCoreErr)
		}
	} else if osType == config.Redhat8 {
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
	} else if osType == config.Redhat7 {
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

	installCommand, err := RetrieveInstallCommand(filepath, osType)
	if err != nil {
		return fmt.Errorf("RetrieveInstallCommand: %w", err)
	}

	err = system.RunCommand(installCommand)
	if err != nil {
		return fmt.Errorf("RunCommand: %w", err)
	}

	successMessage := "\n" + languageTitleCase + " version " + version + " successfully installed!\n"
	fmt.Println(successMessage)
	return nil
}

// Creates the proper command to install R/Python based on the operating system
func RetrieveInstallCommand(filepath string, osType config.OperatingSystem) (string, error) {
	switch osType {
	case config.Ubuntu22, config.Ubuntu20, config.Ubuntu18:
		return "sudo gdebi -n " + filepath, nil
	case config.Redhat7, config.Redhat8:
		return "sudo yum install -y " + filepath, nil
	default:
		return "", errors.New("operating system not supported")
	}
}
