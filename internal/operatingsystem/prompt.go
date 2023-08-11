package operatingsystem

import (
	"fmt"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/prompt"
)

func PromptCloud() (bool, error) {
	confirmText := "Is your instance of Workbench running in a public cloud(AWS, Azure, GCP, etc)?"

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in prereq cloud confirm prompt: %w", err)
	}

	return result, nil
}

func FirewallPrompt() (bool, error) {
	messageText := "Posit products are often blocked by local server firewalls, most organizations do not rely on local firewalls for server security. If your organization controls access to this server with an external firewall, we recommend disabling the local firewall."

	pterm.DefaultParagraph.Println(messageText)
	log.Info(messageText)

	confirmText := "Would you like to disable the local firewall?"

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in firewall disable prompt: %w", err)
	}
	return result, nil
}

func LinuxSecurityPrompt(osType config.OperatingSystem) (bool, error) {
	result := false
	if osType == config.Redhat7 || osType == config.Redhat8 || osType == config.Redhat9 {

		messageText := "SELinux is often enabled by default on Redhat Linux distributions. We recommend that SELinux be disabled, unless you and your organization have specific security requirements that require its use."

		pterm.DefaultParagraph.Println(messageText)
		log.Info(messageText)

		confirmText := "Would you like to disable SELinux on this server?"

		result, err := prompt.PromptConfirm(confirmText)
		if err != nil {
			return false, fmt.Errorf("issue occured in SELinux disable prompt: %w", err)
		}
		return result, nil
	}
	return result, nil
}

func PromptInstallPrereqs() (bool, error) {
	messageText := "In order to install Workbench from start to finish, you will need the following things\n\n" +
		"1. Internet access for this server\n" +
		"2. At least one non-root local Linux user account with a home directory\n" +
		"3. The versions of R, Python and Quarto you would like to install\n" +
		"4. The version of R, Python and Quarto you would like to set as defaults\n" +
		"5. Your Workbench license key string in this form: XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX\n" +
		"6. The location on this server of your SSL key and certificate files (optional)\n" +
		"7. The URL and repo name for your instance of Posit Package Manager (optional)\n" +
		"8. The URL for your instance of Posit Connect (optional)"

	pterm.DefaultBasicText.Println(messageText)
	log.Info(messageText)

	confirmText := "Please confirm that you're ready to install Workbench"
	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in prereq confirm prompt: %w", err)
	}

	return result, nil
}

// PromptUserAccount prompts the user for the name of a local Linux user account to use for verifying the installation
func PromptUserAccount() (string, error) {
	promptText := "Enter a non-root local Linux account username to use for testing the Workbench installation"

	result, err := prompt.PromptText(promptText)
	if err != nil {
		return "", fmt.Errorf("issue occured in user account prompt: %w", err)
	}

	return result, nil
}
