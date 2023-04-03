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

func PromptAndVerifySSL() error {
	certPath, err := PromptSSLFilePath()
	if err != nil {
		return fmt.Errorf("issue with the provided SSL cert path: %w", err)
	}
	keyPath, err := PromptSSLKeyFilePath()
	if err != nil {
		return fmt.Errorf("issue with the provided SSL cert key path: %w", err)
	}
	_, certHostMatch, err := VerifySSLCertAndKey(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("could not verify the SSL cert: %w", err)
	}
	if !certHostMatch {
		proceed, err := PromptMisMatchedHostName()
		if err != nil {
			return fmt.Errorf("Hostname mismatch error: %w", err)
		}
		if proceed {
			fmt.Println("SSL successfully verified, with accepted hostname/cert" +
				"mismatch")
			return nil
		} else {
			return fmt.Errorf("Hostname mismatch error, exit without proceeding: %w", err)
		}

	}
	fmt.Println("SSL successfully verified")
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

func PromptMisMatchedHostName() (bool, error) {
	name := false
	prompt := &survey.Confirm{
		Message: "The hostname of your server and the subject name in the certificate" +
			"don't match. This is common in configurations that include a load balancer" +
			"or a proxy. Please confirm that you want to proceed?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the SSL prompt")
	}
	return name, nil
}
