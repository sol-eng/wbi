package quarto

import (
	"fmt"

	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/prompt"
)

func ScanAndHandleQuartoVersions(osType config.OperatingSystem) error {
	// check if a Workbench bundled version of Quarto exists
	quartoBundled, err := checkForBundledQuartoVersion()
	if err != nil {
		return fmt.Errorf("issue checking for bundled Quarto version: %w", err)
	}
	var quartoBundledVersion string
	if quartoBundled {
		// detect the bundled Quarto version
		quartoBundledVersion, err = ScanForBundledQuartoVersion()
		if err != nil {
			return fmt.Errorf("issue scanning for bundled Quarto version: %w", err)
		}
	}

	// prompt the user to present the bundled version and ask if they want to install any other versions. If nothing is bundled then just ask if they want to install any versions
	quartoInstall, err := PromptQuartoInstall(quartoBundledVersion)
	if err != nil {
		return fmt.Errorf("there was an issue prompting for Quarto install: %w", err)
	}

	if quartoInstall {
		// retrieve other versions and present them to the user
		validQuartoVersions, err := RetrieveValidQuartoVersions()
		if err != nil {
			return fmt.Errorf("there was an issue retrieving valid Quarto versions: %w", err)
		}
		installQuartoVersions, err := QuartoSelectVersionsPrompt(validQuartoVersions)
		if err != nil {
			return fmt.Errorf("issue selecting Quarto versions: %w", err)
		}
		if len(installQuartoVersions) > 0 {
			// install the version(s)
			err = DownloadAndInstallQuartoVersions(installQuartoVersions, osType)
			if err != nil {
				return fmt.Errorf("there was an issue installing Quarto versions: %w", err)
			}
			// ask which version they want to use for default and symlink it to /usr/local/bin/quarto so Jupyter and VS Code will pick it up
			var quartoPaths []string
			if quartoBundled {
				quartoPaths = append(quartoVersionsToPaths(installQuartoVersions), "/usr/lib/rstudio-server/bin/quarto/bin/quarto")
			} else {
				quartoPaths = quartoVersionsToPaths(installQuartoVersions)
			}

			err = checkPromtAndSetQuartoSymlinks(quartoPaths)
			if err != nil {
				return fmt.Errorf("there was an issue setting Quarto symlinks: %w", err)
			}
		} else {
			// continue with the bundled version and symlink it to /usr/local/bin/quarto so Jupyter and VS Code will pick it up
			err = setQuartoSymlinks("/usr/lib/rstudio-server/bin/quarto/bin/quarto", true)
			if err != nil {
				return fmt.Errorf("issue setting Quarto symlinks: %w", err)
			}
		}
	} else {
		// continue with the bundled version and symlink it to /usr/local/bin/quarto so Jupyter and VS Code will pick it up
		err = setQuartoSymlinks("/usr/lib/rstudio-server/bin/quarto/bin/quarto", false)
		if err != nil {
			return fmt.Errorf("issue setting Quarto symlinks: %w", err)
		}
	}

	return nil
}

// quartoVersionsToPaths converts Quarto versions to full paths in /opt
func quartoVersionsToPaths(quartoVersions []string) []string {
	quartoPaths := []string{}
	for _, version := range quartoVersions {
		quartoPaths = append(quartoPaths, "/opt/quarto/"+version+"/bin/quarto")
	}
	return quartoPaths
}

func PromptQuartoInstall(bundledVersion string) (bool, error) {
	var confirmText string
	if bundledVersion == "" {
		confirmText = "Would you like to install Quarto?"
	} else {
		confirmText = "Workbench bundles Quarto version " + bundledVersion + "Would you like to install any different version(s)?"
	}

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in Quarto install confirm prompt: %w", err)
	}

	return result, nil
}

// QuartoSelectVersionsPrompt Prompt asking users which Quarto version(s) they would like to install
func QuartoSelectVersionsPrompt(availableQuartoVersions []string) ([]string, error) {
	promptText := "Which version(s) of Quarto would you like to install?"

	result, err := prompt.PromptMultiSelect(promptText, availableQuartoVersions, []string{availableQuartoVersions[0]})
	if err != nil {
		return []string{}, fmt.Errorf("issue occured in the Quarto versions selection prompt: %w", err)
	}

	return result, nil
}
