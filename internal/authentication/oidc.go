package authentication

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dpastoor/wbi/internal/config"
)

func HandleOIDCConfig(OIDCConfig *config.OIDCConfig) {
	OIDCConfig.AuthOpenID = 1
	OIDCConfig.ClientID = PromptOIDCClientID()
	OIDCConfig.ClientSecret = PromptOIDCClientSecret()
	OIDCConfig.AuthOpenIDIssuer = PromptOIDCIssuerURL()
	OIDCConfig.AuthOpenIDUsernameClaim = PromptOIDCUsernameClaim()
}

func PromptOIDCClientID() string {
	name := ""
	prompt := &survey.Input{
		Message: "OpenID Connect IdP provided client-id:",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func PromptOIDCClientSecret() string {
	name := ""
	prompt := &survey.Input{
		Message: "OpenID Connect IdP provided client-secret:",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func PromptOIDCIssuerURL() string {
	name := ""
	prompt := &survey.Input{
		Message: "OpenID Connect IdP provided issuer URL:",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func PromptOIDCUsernameClaim() string {
	name := "preferred_username"
	prompt := &survey.Input{
		Message: "OpenID Connect IdP provided issuer URL:",
		Default: "preferred_username",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}
