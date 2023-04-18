package jupyter

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/languages"
	"github.com/sol-eng/wbi/internal/workbench"
)

// Prompt asking users if they wish to install Jupyter
func InstallPrompt() (bool, error) {
	name := true
	messageText := "Would you like to install Jupyter?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Jupyter install prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

// Prompt asking users which Python location should Jupyter be installed into
func KernelPrompt(pythonPaths []string) (string, error) {
	// Allow the user to select a version of Python to target
	target := ""
	messageText := "Select a Python kernel to install Jupyter into:"
	prompt := &survey.Select{
		Message: messageText,
		Options: pythonPaths,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", errors.New("there was an issue with the Python selection prompt for installing Jupyter")
	}
	if target == "" {
		return target, errors.New("no Python kernel selected for Jupyter")
	}
	log.Info(messageText)
	log.Info(target)
	return target, nil
}

func ScanPromptInstallAndConfigJupyter() error {
	// scan for Python versions
	pythonVersions, err := languages.ScanForPythonVersions()
	if err != nil {
		return fmt.Errorf("issue occured in scanning for Python versions: %w", err)
	}

	if len(pythonVersions) > 0 {
		jupyterChoice, err := InstallPrompt()
		if err != nil {
			return fmt.Errorf("issue selecting Jupyter: %w", err)
		}

		if jupyterChoice {
			jupyterPythonTarget, err := KernelPrompt(pythonVersions)
			if err != nil {
				return fmt.Errorf("issue selecting Python location for Jupyter: %w", err)
			}
			if jupyterPythonTarget != "" {
				err := InstallJupyter(jupyterPythonTarget)
				if err != nil {
					return fmt.Errorf("issue installing Jupyter: %w", err)
				}

				// the path to jupyter must be set in the config, not python
				pythonSubPath, err := languages.RemovePythonFromPath(jupyterPythonTarget)
				if err != nil {
					return fmt.Errorf("issue removing Python from path: %w", err)
				}
				jupyterPath := pythonSubPath + "/jupyter"
				err = workbench.WriteJupyterConfig(jupyterPath)
				if err != nil {
					return fmt.Errorf("issue writing Jupyter config: %w", err)
				}
			}
		}
	}
	return nil
}
