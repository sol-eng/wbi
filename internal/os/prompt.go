package os

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/sol-eng/wbi/internal/config"
)

func PromptCloud() (bool, error) {
	name := false
	prompt := &survey.Confirm{
		Message: "Is your instance of Workbench running in a public cloud(AWS, Azure, GCP, etc)?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with determining workbench server location")
	}
	return name, nil
}

func FirewallPrompt() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Posit products are often blocked by local server firewalls, most organizations\n " +
			"do not rely on local firewalls for server security. If your organization controls access\n " +
			"to this server with an external firewall, we recommend disabling the local firewall.\n" +
			" Would you like to disable the local firewall?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the disable local firewall prompt")
	}
	return name, nil
}
func LinuxSecurityPrompt(osType config.OperatingSystem) (bool, error) {
	name := false
	if osType == config.Redhat7 || osType == config.Redhat8 {
		name = true
		prompt := &survey.Confirm{
			Message: "SELinux is often enabled by default on Redhat Linux distributions. \nWe recommend that SELinux be" +
				" disabled, unless you and your organization have \nspecific security requirements that require its use.\n" +
				"Would you like to disable SELinux on this server?",
		}
		err := survey.AskOne(prompt, &name)
		if err != nil {
			return false, errors.New("there was an issue with the disable local firewall prompt")
		}
	}
	return name, nil
}

func PromptInstallPrereqs() (bool, error) {
	var name bool
	prompt := &survey.Confirm{
		Message: "In order to install Workbench from start to finish, you will need the following things\n" +
			"1. Internet access for this server\n" +
			"2. The versions of R and Python you would like to install\n" +
			"3. Your Workbench license key string in this form: XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX\n" +
			"4. The location on this server of your SSL key and certificate files (optional)\n" +
			"5. The URL for your instance of Posit Package Manager (optional)\n" +
			"6. The URL for your instance of Posit Connect (optional)\n\n" +
			"Please confirm that you're ready to install Workbench",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the installation confirmation")
	}
	return name, nil
}
