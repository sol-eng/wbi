package os

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/sol-eng/wbi/internal/config"
)

func FirewallPrompt() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Posit products are often blocked by local server firewalls, most organizations " +
			"do not rely on local firewalls for server security. If your organization controls access " +
			"to this server with an external firewall, we recommend disabling the local firewall. Would you like " +
			"to disable the local firewall?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the disable local firewall prompt")
	}
	return name, nil
}

func SELinuxPrompt(osType config.OperatingSystem) (bool, error) {
	name := false
	if osType == config.Redhat7 || osType == config.Redhat8 {
		name = true
		prompt := &survey.Confirm{
			Message: "SELinux is often enabled by default on Redhat Linux distributions. We recommend that SELinux be" +
				"disabled, unless you and your organization have specific security requirements that require its use." +
				"Would you like to disable SELinux on this server?",
		}
		err := survey.AskOne(prompt, &name)
		if err != nil {
			return false, errors.New("there was an issue with the disable local firewall prompt")
		}
	}
	return name, nil
}
