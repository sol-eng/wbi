package license

import (
	"fmt"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
)

// Prompt users for a Workbench license key
func PromptLicense() (string, error) {
	target := ""
	prompt := &survey.Input{
		Message: "Workbench license key:",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a license key: %w", err)
	}
	return target, nil
}

// Activate Workbench based on a license key
func ActivateLicenseKey(licenseKey string) error {
	cmdLicense := "sudo rstudio-server license-manager activate " + licenseKey
	cmd := exec.Command("/bin/sh", "-c", cmdLicense)
	stdout, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("issue activating Workbench license: %w", err)
	}

	fmt.Print(string(stdout))
	// TODO add a real check that Workbench is activated
	fmt.Println("Workbench has been successfully activated!")
	return nil
}
