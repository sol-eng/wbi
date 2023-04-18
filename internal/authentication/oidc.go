package authentication

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
)

// Prompt asking users to provide a client-id for OIDC
func PromptOIDCClientID() (string, error) {
	name := ""
	messageText := "OpenID Connect IdP provided client-id:"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the OIDC client-id prompt")
	}
	log.Info(messageText)
	log.Info(name)
	return name, nil
}

// Prompt asking users to provide a client-secret for OIDC
func PromptOIDCClientSecret() (string, error) {
	name := ""
	messageText := "OpenID Connect IdP provided client-secret:"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the OIDC client-secret prompt")
	}
	log.Info(messageText)
	log.Info(name)
	return name, nil
}

// Prompt asking users to provide an issuer URL for OIDC
func PromptOIDCIssuerURL() (string, error) {
	name := ""
	messageText := "OpenID Connect IdP provided issuer URL:"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the OIDC IdP issuer URL prompt")
	}
	log.Info(messageText)
	log.Info(name)
	return name, nil
}

// Prompt asking users to provide a username claim for OIDC
func PromptOIDCUsernameClaim() (string, error) {
	name := "preferred_username"
	messageText := "OpenID Connect IdP provided username claim:"
	prompt := &survey.Input{
		Message: messageText,
		Default: "preferred_username",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the OIDC IdP username claim prompt")
	}
	log.Info(messageText)
	log.Info(name)
	return name, nil
}
