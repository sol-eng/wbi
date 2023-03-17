package os

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
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
