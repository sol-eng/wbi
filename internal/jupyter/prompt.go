package jupyter

import (
	"fmt"
	"regexp"

	"github.com/sol-eng/wbi/internal/languages"
	"github.com/sol-eng/wbi/internal/prompt"
	"github.com/sol-eng/wbi/internal/workbench"
)

// Prompt asking users if they wish to install Jupyter
func InstallPrompt() (bool, error) {
	confirmText := "Would you like to install Jupyter?"

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in Jupyter install confirm prompt: %w", err)
	}

	return result, nil
}

// Prompt asking users which Python location should Jupyter be installed into
func KernelPrompt(pythonPaths []string) (string, error) {
	promptText := "Select a Python kernel to install Jupyter into"

	result, err := prompt.PromptSingleSelect(promptText, pythonPaths, pythonPaths[0])
	if err != nil {
		return "", fmt.Errorf("issue occured in Jupyter kernel selection prompt: %w", err)
	}

	return result, nil
}

// Prompt asking users which additional Python location should be registered as Jupyter kernels
func AdditionalKernelPrompt(pythonPaths []string, defaultPythonPaths []string) ([]string, error) {

	promptText := "Which of the remaining Python versions would you like to have automatically registered as Jupyter kernels? (select none to skip this step)"

	result, err := prompt.PromptMultiSelect(promptText, pythonPaths, defaultPythonPaths)
	if err != nil {
		return []string{}, fmt.Errorf("issue occured in Jupyter kernel registration selection prompt: %w", err)
	}

	return result, nil
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
				if len(pythonVersionsLeft) > 0 {
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
