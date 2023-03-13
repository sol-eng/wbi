package os

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
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
