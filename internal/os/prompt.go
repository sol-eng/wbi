package os

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/sol-eng/wbi/internal/config"
)

func PromptCloud(osType config.OperatingSystem) (bool, error) {
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
	name := false
	prompt := &survey.Confirm{
		Message: "Please confirm that you're ready to install Workbench.",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the installation confirmation")
	}
	return name, nil
}
