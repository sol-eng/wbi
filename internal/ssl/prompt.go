package ssl

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
)

func PromptSSL() bool {
	name := false
	prompt := &survey.Confirm{
		Message: "Would you like to use SSL?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func PromptSSLFilePath() string {
	target := ""
	prompt := &survey.Input{
		Message: "Filepath to SSL certificate:",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		log.Fatal(err)
	}
	return target
}

func PromptSSLKeyFilePath() string {
	target := ""
	prompt := &survey.Input{
		Message: "Filepath to SSL certificate key:",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		log.Fatal(err)
	}
	return target
}
