package authentication

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
)

// Prompt asking users to provide a client-id for OIDC
func PromptOIDCClientID() (string, error) {
	name := ""
	prompt := &survey.Input{
		Message: "OpenID Connect IdP provided client-id:",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the OIDC client-id prompt")
	}
	return name, nil
}

// Prompt asking users to provide a client-secret for OIDC
func PromptOIDCClientSecret() (string, error) {
	name := ""
	prompt := &survey.Input{
		Message: "OpenID Connect IdP provided client-secret:",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the OIDC client-secret prompt")
	}
	return name, nil
}

// Prompt asking users to provide an issuer URL for OIDC
func PromptOIDCIssuerURL() (string, error) {
	name := ""
	prompt := &survey.Input{
		Message: "OpenID Connect IdP provided issuer URL:",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the OIDC IdP issuer URL prompt")
	}
	return name, nil
}
