package quarto

import (
	"fmt"
	"os"

	"github.com/sol-eng/wbi/internal/system"
)

// ScanForBundledQuartoVersion scans for the bundled version of Quarto
func ScanForBundledQuartoVersion() (string, error) {
	quartoPath := "/usr/lib/rstudio-server/bin/quarto/bin/quarto"
	versionCommand := quartoPath + " --version"
	quartoVersion, err := system.RunCommandAndCaptureOutput(versionCommand, false, 0)
	if err != nil {
		return "", fmt.Errorf("issue finding Quarto version: %w", err)
	}
	return quartoVersion, nil
}

func checkForBundledQuartoVersion() (bool, error) {
	// check if /usr/lib/rstudio-server/bin/quarto/bin/quarto exists
	_, err := os.Stat("/usr/lib/rstudio-server/bin/quarto/bin/quarto")
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func CheckAndSetQuartoSymlink(quartoPath string) error {
	// check if Quarto has already been symlinked
	quartoSymlinked := checkIfQuartoSymlinkExists()
	if !quartoSymlinked {
		err := setQuartoSymlinks(quartoPath, true)
		if err != nil {
			return fmt.Errorf("issue setting Quarto symlinks: %w", err)
		}
	} else {
		system.PrintAndLogInfo("Quarto symlink already exist, skipping symlink creation")
	}
	return nil
}

func checkPromtAndSetQuartoSymlinks(quartoPaths []string) error {
	// check if Quarto has already been symlinked
	quartoSymlinked := checkIfQuartoSymlinkExists()
	if (len(quartoPaths) > 0) && !quartoSymlinked {
		err := promptAndSetQuartoSymlink(quartoPaths)
		if err != nil {
			return fmt.Errorf("issue setting Quarto symlinks: %w", err)
		}
	}
	return nil
}

// checkIfQuartoSymlinkExists checks if the Quarto symlink exists
func checkIfQuartoSymlinkExists() bool {
	_, err := os.Stat("/usr/local/bin/quarto")
	if err != nil {
		return false
	}

	system.PrintAndLogInfo("\nAn existing Quarto symlink has been detected (/usr/local/bin/quarto)")
	return true
}
