package jupyter

import (
	"fmt"

	"github.com/sol-eng/wbi/internal/languages"
	"github.com/sol-eng/wbi/internal/system"
	"github.com/sol-eng/wbi/internal/workbench"
)

// InstallJupyter installs jupypter pip stuff
func InstallJupyter(pythonPath string) error {
	JupyterComponentsErr := InstallJupyterAndComponents(pythonPath)
	if JupyterComponentsErr != nil {
		return fmt.Errorf("InstallJupyterAndComponents: %w", JupyterComponentsErr)
	}
	JupyterNotebookExtensionsErr := InstallAndEnableJupyterNotebookExtensions(pythonPath)
	if JupyterNotebookExtensionsErr != nil {
		return fmt.Errorf("InstallAndEnableJupyterNotebookExtensions: %w", JupyterNotebookExtensionsErr)
	}
	return nil
}

// Install various Jupyter related packages from PyPI
func InstallJupyterAndComponents(pythonPath string) error {
	licenseCommand := pythonPath + " -m pip install jupyter jupyterlab rsp_jupyter rsconnect_jupyter workbench_jupyterlab"
	err := system.RunCommand(licenseCommand, true, 2)
	if err != nil {
		return fmt.Errorf("issue installing Jupyter with the command '%s': %w", licenseCommand, err)
	}
	// TODO add some proper tests to ensure Jupyter is working
	fmt.Println("\nJupyter has been successfully installed!")
	return nil
}

// Install and enable various Jupyter notebook extensions
func InstallAndEnableJupyterNotebookExtensions(pythonPath string) error {

	pythonPathShort, err := languages.RemovePythonFromPath(pythonPath)
	if err != nil {
		return fmt.Errorf("issue shortening Python path: %w", err)
	}

	commands := []string{
		"jupyter-nbextension install --sys-prefix --py rsp_jupyter",
		"jupyter-nbextension enable --sys-prefix --py rsp_jupyter",
		"jupyter-nbextension install --sys-prefix --py rsconnect_jupyter",
		"jupyter-nbextension enable --sys-prefix --py rsconnect_jupyter",
		"jupyter-serverextension enable --sys-prefix --py rsconnect_jupyter",
	}

	for _, command := range commands {
		installCommand := pythonPathShort + "/" + command
		err := system.RunCommand(installCommand, true, 0)
		if err != nil {
			return fmt.Errorf("issue installing Jupyter notebook extension with the command '%s': %w", installCommand, err)
		}
	}

	// TODO add some proper tests to ensure Jupyter notebook extensions are working
	fmt.Println("\nJupyter notebook extensions have been successfully installed and enabled!")
	return nil
}

func InstallAndConfigJupyter(pythonPath string) error {
	err := InstallJupyter(pythonPath)
	if err != nil {
		return fmt.Errorf("issue installing Jupyter: %w", err)
	}
	// the path to jupyter must be set in the config, not python
	pythonSubPath, err := languages.RemovePythonFromPath(pythonPath)
	if err != nil {
		return fmt.Errorf("issue removing Python from path: %w", err)
	}
	jupyterPath := pythonSubPath + "/jupyter"
	err = workbench.WriteJupyterConfig(jupyterPath)
	if err != nil {
		return fmt.Errorf("issue writing Jupyter config: %w", err)
	}
	return nil
}
