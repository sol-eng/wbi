package operatingsystem

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/config"
)

func PromptCloud() (bool, error) {
	name := false
	messageText := "Is your instance of Workbench running in a public cloud(AWS, Azure, GCP, etc)?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with determining workbench server location")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

func FirewallPrompt() (bool, error) {
	name := true
	messageText := "Posit products are often blocked by local server firewalls, most organizations\n " + "do not rely on local firewalls for server security. If your organization controls access\n " + "to this server with an external firewall, we recommend disabling the local firewall.\n" + " Would you like to disable the local firewall?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the disable local firewall prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

func LinuxSecurityPrompt(osType config.OperatingSystem) (bool, error) {
	name := false
	if osType == config.Redhat7 || osType == config.Redhat8 || osType == config.Redhat9 {
		name = true
		messageText := "SELinux is often enabled by default on Redhat Linux distributions. \nWe recommend that SELinux be" + " disabled, unless you and your organization have \nspecific security requirements that require its use.\n" + "Would you like to disable SELinux on this server?"
		prompt := &survey.Confirm{
			Message: messageText,
		}
		err := survey.AskOne(prompt, &name)
		if err != nil {
			return false, errors.New("there was an issue with the disable local firewall prompt")
		}
		log.Info(messageText)
		log.Info(fmt.Sprintf("%v", name))
	}
	return name, nil
}

func PromptInstallPrereqs() (bool, error) {
	var name bool
	messageText := "In order to install Workbench from start to finish, you will need the following things\n" +
		"1. Internet access for this server\n" +
		"2. The versions of R and Python you would like to install\n" +
		"3. The version of R and Python you would like to set as defaults\n" +
		"4. Your Workbench license key string in this form: XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX\n" +
		"5. The location on this server of your SSL key and certificate files (optional)\n" +
		"6. The IdP Metadata URL (optional, only for SAML SSO integration)\n" +
		"7. The IdP Issuer URL, client id and client secret (optional, only for OIDC SSO integration)\n" +
		"8. The URL and repo name for your instance of Posit Package Manager (optional)\n" +
		"9. The URL for your instance of Posit Connect (optional)\n\n" +
		"Please confirm that you're ready to install Workbench"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the installation confirmation")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}
