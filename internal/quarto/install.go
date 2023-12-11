package quarto

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/samber/lo"
	"github.com/sol-eng/wbi/internal/config"
	cmdlog "github.com/sol-eng/wbi/internal/logging"
	"github.com/sol-eng/wbi/internal/prompt"
	"github.com/sol-eng/wbi/internal/system"
)

func RetrieveValidQuartoVersions() ([]string, error) {
	// TODO automate the retrieving the list of valid versions
	return []string{"1.3.340", "1.2.475", "1.1.189", "1.0.38"}, nil
}

func ValidateQuartoVersions(quartoVersions []string) error {
	availQuartoVersions, err := RetrieveValidQuartoVersions()
	if err != nil {
		return fmt.Errorf("error retrieving valid Quarto versions: %w", err)
	}
	for _, quartoVersion := range quartoVersions {
		if !lo.Contains(availQuartoVersions, quartoVersion) {
			return errors.New("version " + quartoVersion + " is not a valid Quarto version")
		}
	}
	return nil
}

func DownloadAndInstallQuartoVersions(quartoVersions []string, osType config.OperatingSystem) error {
	for _, quartoVersion := range quartoVersions {
		err := DownloadAndInstallQuarto(quartoVersion, osType)
		if err != nil {
			return fmt.Errorf("issue installing Quarto version: %w", err)
		}
	}
	return nil
}

func DownloadAndInstallQuarto(quartoVersion string, osType config.OperatingSystem) error {
	// Find URL
	quartoURL := generateQuartoInstallURL(quartoVersion, osType)
	// Download installer
	installerPath, err := downloadFileQuarto(quartoURL, quartoVersion, osType)
	if err != nil {
		return fmt.Errorf("DownloadFileQuarto: %w", err)
	}
	// Install Quarto
	err = installQuarto(installerPath, osType, quartoVersion, true)
	if err != nil {
		return fmt.Errorf("InstallQuarto: %w", err)
	}
	// save to command log
	quartoPath := fmt.Sprintf("/opt/quarto/%s", quartoVersion)
	cmdlog.Info("curl -o quarto.tar.gz -L " + quartoURL)
	cmdlog.Info("mkdir -p " + quartoPath)
	cmdlog.Info(fmt.Sprintf(`tar -zxvf quarto.tar.gz -C "%s" --strip-components=1`, quartoPath))
	cmdlog.Info("rm quarto.tar.gz")
	return nil
}

func generateQuartoInstallURL(quartoVersion string, osType config.OperatingSystem) string {
	// treat RHEL 7 differently as specified here: https://docs.posit.co/resources/install-quarto/#specify-quarto-version-tar
	var url string
	if osType == config.Redhat7 {
		url = fmt.Sprintf("https://github.com/quarto-dev/quarto-cli/releases/download/v%s/quarto-%s-linux-rhel7-amd64.tar.gz", quartoVersion, quartoVersion)
	} else {
		url = fmt.Sprintf("https://github.com/quarto-dev/quarto-cli/releases/download/v%s/quarto-%s-linux-amd64.tar.gz", quartoVersion, quartoVersion)
	}
	return url
}

func downloadFileQuarto(url string, version string, osType config.OperatingSystem) (string, error) {
	system.PrintAndLogInfo("Downloading Quarto Version: " + version + " installer from: " + url)

	// Create the file
	filename := "*_" + fmt.Sprintf("quarto-%s-linux-amd64.tar.gz", version)
	tmpFile, err := os.CreateTemp("", filename)
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
		return "", errors.New("error downloading " + filename + " installer")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("error retrieving " + filename + " installer")
	}

	// Writer the body to file
	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// Installs Quarto
func installQuarto(filepath string, osType config.OperatingSystem, version string, save bool) error {
	// create the /opt/quarto directory if it doesn't exist
	path := fmt.Sprintf("/opt/quarto/%s", version)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
	}

	installCommand := fmt.Sprintf(`tar -zxvf "%s" -C "%s" --strip-components=1`, filepath, path)

	err := system.RunCommand(installCommand, false, 0, false)
	if err != nil {
		return fmt.Errorf("the command '%s' failed to run: %w", installCommand, err)
	}

	successMessage := "\nQuarto version " + version + " successfully installed!\n"
	system.PrintAndLogInfo(successMessage)
	return nil
}

// promptAndSetQuartoSymlinks prompts user to set the Quarto symlink
func promptAndSetQuartoSymlink(quartoPaths []string) error {
	setQuartoSymlinkChoice, err := quartoSymlinkPrompt()
	if err != nil {
		return fmt.Errorf("an issue occured during the selection of Quarto symlink choice: %w", err)
	}
	if setQuartoSymlinkChoice {
		quartoPathChoice, err := quartoLocationSymlinksPrompt(quartoPaths)
		if err != nil {
			return fmt.Errorf("issue selecting Quarto binary to add symlinks: %w", err)
		}
		err = setQuartoSymlinks(quartoPathChoice, true)
		if err != nil {
			return fmt.Errorf("issue setting Quarto symlinks: %w", err)
		}

		system.PrintAndLogInfo("\n " + quartoPathChoice + " has been successfully symlinked and will be available on the default system PATH.\n")
	}
	return nil
}

// setQuartoSymlinks sets the Quarto symlink
func setQuartoSymlinks(quartoPath string, display bool) error {
	quartoCommand := "ln -s " + quartoPath + " /usr/local/bin/quarto"
	err := system.RunCommand(quartoCommand, display, 0, true)
	if err != nil {
		return fmt.Errorf("error setting Quarto symlink with the command '%s': %w", quartoCommand, err)
	}
	return nil
}

// quartoLocationSymlinksPrompt asks users which Quarto binary they want to symlink
func quartoLocationSymlinksPrompt(quartoPaths []string) (string, error) {
	messageText := "Select a Quarto binary to symlink:"

	result, err := prompt.PromptSingleSelect(messageText, quartoPaths, quartoPaths[0])
	if err != nil {
		return "", fmt.Errorf("issue occured in Quarto selection prompt for symlinking: %w", err)
	}
	if result == "" {
		return result, errors.New("no Quarto binary selected to be symlinked")
	}

	return result, nil
}

// quartoSymlinkPrompt asks users if they would like to set the quarto symlink
func quartoSymlinkPrompt() (bool, error) {
	confirmText := `Would you like to symlink a Quarto version to make it available on PATH? This is recommended so Workbench can default to this version of Quarto in each of the IDEs and users can type "quarto" in the terminal.`

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in Quarto symlink confirm prompt: %w", err)
	}

	return result, nil
}
