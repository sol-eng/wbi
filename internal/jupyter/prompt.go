package jupyter

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dpastoor/wbi/internal/config"
)

// Prompt asking users if they wish to install Jupyter
func InstallPrompt() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to install Jupyter?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Jupyter install prompt")
	}
	return name, nil
}

// Prompt asking users which Python location should Jupyter be installed into
func KernelPrompt(PythonConfig *config.PythonConfig) (string, error) {
	// Allow the user to select a version of Python to target
	target := ""
	prompt := &survey.Select{
		Message: "Select a Python kernel to install Jupyter into:",
		Options: PythonConfig.Paths,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", errors.New("there was an issue with the Python selection prompt for installing Jupyter")
	}
	if target == "" {
		return target, errors.New("no Python kernel selected for Jupyter")
	}
	return target, nil
}
