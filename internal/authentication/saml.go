package authentication

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sol-eng/wbi/internal/config"
)

// Run functions and store values in the SAMLConfig
func HandleSAMLConfig(SAMLConfig *config.SAMLConfig) error {
	SAMLConfig.AuthSAML = 1

	AuthSamlSpAttributeUsername, err := PromptSAMLAttribute()
	SAMLConfig.AuthSamlSpAttributeUsername = AuthSamlSpAttributeUsername
	if err != nil {
		return fmt.Errorf("PromptSAMLAttribute: %w", err)
	}

	AuthSamlMetadataURL, err := PromptSAMLMetadataURL()
	SAMLConfig.AuthSamlMetadataURL = AuthSamlMetadataURL
	if err != nil {
		return fmt.Errorf("PromptSAMLMetadataURL: %w", err)
	}
	return nil
}

// Prompt asking users to provide a username attribute for SAML
func PromptSAMLAttribute() (string, error) {
	name := "Username"
	prompt := &survey.Input{
		Message: "SAML IdP username attribute:",
		Default: "Username",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the SAML username attribute prompt")
	}
	return name, nil
}

// Prompt asking users to provide a metadata URL for SAML
func PromptSAMLMetadataURL() (string, error) {
	name := ""
	prompt := &survey.Input{
		Message: "SAML IdP metadata URL:",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return "", errors.New("there was an issue with the SAML IdP URL prompt")
	}
	return name, nil
}
