package install

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/system"
)

// Installs R/Python in a certain way based on the operating system
func InstallLanguage(language string, filepath string, osType config.OperatingSystem, version string) error {
	languageTitleCase := strings.Title(language)

	installCommand, err := RetrieveInstallCommand(filepath, osType)
	if err != nil {
		return fmt.Errorf("RetrieveInstallCommand: %w", err)
	}

	err = system.RunCommand(installCommand, false, 0)
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
		return "gdebi -n " + filepath, nil
	case config.Redhat7, config.Redhat8:
		return "yum install -y " + filepath, nil
	default:
		return "", errors.New("operating system not supported")
	}
}
