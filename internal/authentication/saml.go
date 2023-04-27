package authentication

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
)

// Prompt asking users to provide a username attribute for SAML
func PromptSAMLAttribute() (string, error) {
	name := "Username"
	messageText := "SAML IdP username attribute:"
	prompt := &survey.Input{
		Message: messageText,
		Default: "Username",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the SAML username attribute prompt")
	}
	log.Info(messageText)
	log.Info(name)
	return name, nil
}

// Prompt asking users to provide a metadata URL for SAML
func PromptSAMLMetadataURL() (string, error) {
	name := ""
	messageText := "SAML IdP metadata URL:"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the SAML IdP URL prompt")
	}
	log.Info(messageText)
	log.Info(name)
	return name, nil
}
