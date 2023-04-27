package jupyter

import (
	"errors"
	"fmt"
	"regexp"

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

// Prompt asking users which additional Python location should be registered as Jupyter kernels
func AdditionalKernelPrompt(pythonPaths []string, defaultPythonPaths []string) ([]string, error) {
	// Allow the user to select multiple versions
	var qs = []*survey.Question{
		{
			Name: "kernelprompt",
			Prompt: &survey.MultiSelect{
				Message: "Which of the remaining Python versions would you like to have automatically registered as Jupyter kernels? (select none to skip this step)",
				Options: pythonPaths,
				Default: defaultPythonPaths,
			},
		},
	}
	kernelAnswers := struct {
		Versions []string `survey:"kernelprompt"`
	}{}

	err := survey.Ask(qs, &kernelAnswers, survey.WithRemoveSelectAll(), survey.WithRemoveSelectNone())
	if err != nil {
		return []string{}, errors.New("there was an issue with the languages prompt")
	}
	return kernelAnswers.Versions, nil
}

func removeString(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
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
				err = InstallJupyter(jupyterPythonTarget)
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

				// prompt and handle additional kernels to be registered
				// remove the primary Jupyter Python version
				pythonVersionsLeft := removeString(pythonVersions, jupyterPythonTarget)
				// remove any non opt locations from the default selections
				defaultPythonVersions := removeNonOptPython(pythonVersionsLeft)
				additionalPythonTargets, err := AdditionalKernelPrompt(pythonVersionsLeft, defaultPythonVersions)
				if err != nil {
					return fmt.Errorf("issue selecting additional Python kernels to register: %w", err)
				}
				// if one or more versions is selected then automatically register them
				if len(additionalPythonTargets) > 0 {
					err = RegisterJupyterKernels(additionalPythonTargets)
					if err != nil {
						return fmt.Errorf("issue registering additional Python kernels: %w", err)
					}
				}
			}
		}
	}
	return nil
}

func removeNonOptPython(pythonPaths []string) []string {
	anyOptLocations := []string{}
	for _, value := range pythonPaths {
		matched, err := regexp.MatchString(".*/opt.*", value)
		if err == nil && matched {
			anyOptLocations = append(anyOptLocations, value)
		}
	}
	return anyOptLocations
}
