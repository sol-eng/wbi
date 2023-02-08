package jupyter

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
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

// Remove python or python3 from the end of a path so the directory can be used for commands
func RemovePythonFromPath(pythonPath string) (string, error) {
	if _, err := regexp.MatchString(".*/python.*", pythonPath); err == nil {
		i := strings.LastIndex(pythonPath, "/python")
		excludingLast := pythonPath[:i] + strings.Replace(pythonPath[i:], "/python", "", 1)
		return excludingLast, nil
	} else if _, err := regexp.MatchString(".*/python3.*", pythonPath); err == nil {
		i := strings.LastIndex(pythonPath, "/python3")
		excludingLast := pythonPath[:i] + strings.Replace(pythonPath[i:], "/python3", "", 1)
		return excludingLast, nil
	} else {
		return pythonPath, nil
	}
}

// Install and enable various Jupyter notebook extensions
func InstallAndEnableJupyterNotebookExtensions(pythonPath string) error {

	pythonPathShort, err := RemovePythonFromPath(pythonPath)
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
		installCommand := "sudo " + pythonPathShort + "/" + command
		cmd := exec.Command("/bin/sh", "-c", installCommand)
		stdout, err := cmd.Output()

		fmt.Println(string(stdout))
		if err != nil {
			return fmt.Errorf("issue installing Jupyter notebook extensions: %w", err)
		}
	}

	// TODO add some proper tests to ensure Jupyter notebook extensions are working
	fmt.Println("Jupyter notebook extensions have been successfully installed and enabled!")
	return nil
}

// Install various Jupyter related packages from PyPI
func InstallJupyterAndComponents(pythonPath string) error {
	cmdLicense := "sudo " + pythonPath + " -m pip install jupyter jupyterlab rsp_jupyter rsconnect_jupyter workbench_jupyterlab"
	cmd := exec.Command("/bin/sh", "-c", cmdLicense)
	stdout, err := cmd.Output()

	fmt.Println(string(stdout))
	if err != nil {
		return fmt.Errorf("issue installing Jupyter: %w", err)
	}

	// TODO add some proper tests to ensure Jupyter is working
	fmt.Println("Jupyter has been successfully installed!")
	return nil
}
