package authentication

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sol-eng/wbi/internal/config"
)

// Run functions and store values in the OIDCConfig
func HandleOIDCConfig(OIDCConfig *config.OIDCConfig) error {
	OIDCConfig.AuthOpenID = 1

	ClientID, err := PromptOIDCClientID()
	OIDCConfig.ClientID = ClientID
	if err != nil {
		return fmt.Errorf("PromptOIDCClientID: %w", err)
	}

	ClientSecret, err := PromptOIDCClientSecret()
	OIDCConfig.ClientSecret = ClientSecret
	if err != nil {
		return fmt.Errorf("PromptOIDCClientSecret: %w", err)
	}

	AuthOpenIDIssuer, err := PromptOIDCIssuerURL()
	OIDCConfig.AuthOpenIDIssuer = AuthOpenIDIssuer
	if err != nil {
		return fmt.Errorf("PromptOIDCIssuerURL: %w", err)
	}

	AuthOpenIDUsernameClaim, err := PromptOIDCUsernameClaim()
	OIDCConfig.AuthOpenIDUsernameClaim = AuthOpenIDUsernameClaim
	if err != nil {
		return fmt.Errorf("PromptOIDCUsernameClaim: %w", err)
	}
	return nil
}

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

// Prompt asking users to provide a username claim for OIDC
func PromptOIDCUsernameClaim() (string, error) {
	name := "preferred_username"
	prompt := &survey.Input{
		Message: "OpenID Connect IdP provided username claim:",
		Default: "preferred_username",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the OIDC IdP username claim prompt")
	}
	return name, nil
}
