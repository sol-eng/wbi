package os

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/install"
	"github.com/sol-eng/wbi/internal/system"
)

func InstallPrereqs(osType config.OperatingSystem) error {
	fmt.Println("Installing prerequisites...\n")
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
	fmt.Println("\nPrerequisites successfully installed!\n")
	return nil
}

// Installs Gdebi Core
func InstallGdebiCore() error {
	gdebiCoreCommand := "apt-get install -y gdebi-core"
	err := system.RunCommand(gdebiCoreCommand)
	if err != nil {
		return fmt.Errorf("issue installing gdebi-core: %w", err)
	}

	fmt.Println("\ngdebi-core has been successfully installed!\n")
	return nil
}

// Upgrades Apt
func UpgradeApt() error {
	aptUpgradeCommand := "apt-get update"
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
	dnfPluginsCoreCommand := "dnf install -y dnf-plugins-core"
	err := system.RunCommand(dnfPluginsCoreCommand)
	if err != nil {
		return fmt.Errorf("issue installing dnf-plugins-core: %w", err)
	}

	enableCodeReadyCommand := `dnf config-manager --set-enabled "codeready-builder-for-rhel-8-*-rpms"`
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
	yumUtilsCommand := "yum install -y yum-utils"
	err := system.RunCommand(yumUtilsCommand)
	if err != nil {
		return fmt.Errorf("issue installing yum-utils: %w", err)
	}

	enableOptionalCommand := `yum-config-manager --enable "rhel-*-optional-rpms"`
	err = system.RunCommand(enableOptionalCommand)
	if err != nil {
		return fmt.Errorf("issue enabling the optional repo: %w", err)
	}

	fmt.Println("\nThe Optional repository has been successfully enabled!\n")
	return nil
}

// Enable the Extra Packages for Enterprise Linux (EPEL) repository
func EnableEPELRepo(osType config.OperatingSystem) error {
	var EPELURL string
	if osType == config.Redhat8 {
		EPELURL = "https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm"
	} else if osType == config.Redhat7 {
		EPELURL = "https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm"
	} else {
		return errors.New("operating system not supported")
	}

	EPELCommand, err := install.RetrieveInstallCommand(EPELURL, osType)
	if err != nil {
		return fmt.Errorf("issue retrieving EPEL install command: %w", err)
	}
	commandOutput, err := system.RunCommandAndCaptureOutput(EPELCommand)
	if err != nil {
		if strings.Contains(commandOutput, "does not update installed package") && osType == config.Redhat7 {
			fmt.Println("\nThe Extra Packages for Enterprise Linux (EPEL) repository was already enabled.\n")
			return nil
		}
		return fmt.Errorf("issue enabling EPEL repo: %w", err)
	}
	if err != nil {
		return fmt.Errorf("issue enabling EPEL repo: %w", err)
	}

	fmt.Println("\nThe Extra Packages for Enterprise Linux (EPEL) repository has been successfully enabled!\n")
	return nil
}

// Disable local firewall on server
func DisableFirewall(osType config.OperatingSystem) error {
	var FWCommand string

	switch osType {
	case config.Ubuntu18, config.Ubuntu20, config.Ubuntu22:
		FWCommand = "ufw disable"
	case config.Redhat7, config.Redhat8:
		FWCommand = "systemctl stop firewalld && systemctl disable firewalld"
	default:
		return errors.New("Unsupported OS, setting FWCommand failed")
	}
	err := system.RunCommand(FWCommand)
	if err != nil {
		return fmt.Errorf("issue disabling system firewall: %w", err)
	}

	fmt.Println("\nThe system firewall has been successfully disabled!\n")
	return nil
}

func DisableLinuxSecurity() error {

	setenforceCommand := "setenforce 0"
	err := system.RunCommand(setenforceCommand)
	if err != nil {
		return fmt.Errorf("issue stopping selinux enforcement via setenforce: %w", err)
	}

	disableSELinuxCommand := "sed -i s/^SELINUX=.*$/SELINUX=disabled/ /etc/selinux/config"
	err = system.RunCommand(disableSELinuxCommand)
	if err != nil {
		return fmt.Errorf("issue disabling selinux: %w", err)
	}

	fmt.Println("\nThe SELinux has been successfully changed to permissive mode, and will be disabled on next reboot!\n")
	return nil
}
