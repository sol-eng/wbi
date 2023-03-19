package ssl

import (
	"errors"
	"fmt"

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

func PromptVerifyAndConfigSSL() error {
	certPath, err := PromptSSLFilePath()
	if err != nil {
		return fmt.Errorf("issue with the provided SSL cert path: %w", err)
	}
	keyPath, err := PromptSSLKeyFilePath()
	if err != nil {
		return fmt.Errorf("issue with the provided SSL cert key path: %w", err)
	}
	verifySSLCert := VerifySSLCertAndKey(certPath, keyPath)
	if verifySSLCert != nil {
		return fmt.Errorf("could not verify the SSL cert: %w", err)
	}
	fmt.Println("SSL successfully setup and verified")
	return nil
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
