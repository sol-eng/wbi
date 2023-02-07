package license

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
)

func PromptLicense() string {
	target := ""
	prompt := &survey.Input{
		Message: "Workbench license key:",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		log.Fatal(err)
	}
	return target
}

func ActivateLicenseKey(licenseKey string) {
	cmdLicense := "sudo rstudio-server license-manager activate " + licenseKey
	cmd := exec.Command("/bin/sh", "-c", cmdLicense)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(stdout))
}
