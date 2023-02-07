package jupyter

import (
	"errors"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dpastoor/wbi/internal/config"
)

func InstallPrompt() bool {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to install Jupyter?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func KernelPrompt(PythonConfig *config.PythonConfig) (string, error) {
	// Allow the user to select a version of Python to target
	target := ""
	prompt := &survey.Select{
		Message: "Select a Python kernel to install Jupyter into:",
		Options: PythonConfig.Paths,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		log.Fatal(err)
	}
	if target == "" {
		return target, errors.New("no Python kernel selected for Jupyter")
	}
	return target, nil
}
