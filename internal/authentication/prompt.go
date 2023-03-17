package authentication

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sol-eng/wbi/internal/config"
)

// Prompt asking users if they wish to setup Authentication
func PromptAuth() (bool, error) {
	name := false
	prompt := &survey.Confirm{
		Message: "Would you like to setup Authentication?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Authentication prompt")
	}
	return name, nil
}

func PromptAndConvertAuthType() (config.AuthType, error) {

	authChoiceRaw, err := PromptAuthentication()
	if err != nil {
		return config.Other, fmt.Errorf("PromptAuthentication: %w", err)
	}
	authChoice, err := ConvertAuthType(authChoiceRaw)
	if err != nil {
		return config.Other, fmt.Errorf("PromptAuthentication: %w", err)
	}
	return authChoice, nil
}

// Prompt asking user for their desired authentication method
func PromptAuthentication() (string, error) {
	name := ""
	prompt := &survey.Select{
		Message: "Choose an authentication method:",
		Options: []string{"SAML", "OIDC", "Active Directory/LDAP", "PAM", "Other"},
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the authentication prompt")
	}
	return name, nil
}

// Convert authChoice from a string to a proper AuthType type
func ConvertAuthType(authChoice string) (config.AuthType, error) {
	switch authChoice {
	case "SAML":
		return config.SAML, nil
	case "OIDC":
		return config.OIDC, nil
	case "Active Directory/LDAP":
		return config.LDAP, nil
	case "PAM":
		return config.PAM, nil
	case "Other":
		return config.Other, nil
	}
	return config.Other, errors.New("unknown AuthType was selected")
}
