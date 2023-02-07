package ssl

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
)

// Prompt asking users if they wish to use SSL
func PromptSSL() (bool, error) {
	name := false
	prompt := &survey.Confirm{
		Message: "Would you like to use SSL?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the SSL prompt")
	}
	return name, nil
}

// Prompt asking users for a filepath to their SSL cert
func PromptSSLFilePath() (string, error) {
	target := ""
	prompt := &survey.Input{
		Message: "Filepath to SSL certificate:",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", errors.New("there was an issue with the SSL cert path prompt")
	}
	return target, nil
}

// Prompt asking users for a filepath to their SSL cert key
func PromptSSLKeyFilePath() (string, error) {
	target := ""
	prompt := &survey.Input{
		Message: "Filepath to SSL certificate key:",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", errors.New("there was an issue with the SSL cert key path prompt")
	}
	return target, nil
}
