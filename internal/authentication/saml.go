package authentication

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dpastoor/wbi/internal/config"
)

func HandleSAMLConfig(SAMLConfig *config.SAMLConfig) {
	SAMLConfig.AuthSAML = 1
	SAMLConfig.AuthSamlSpAttributeUsername = PromptSAMLAttribute()
	SAMLConfig.AuthSamlMetadataURL = PromptSAMLMetadataURL()
}

func PromptSAMLAttribute() string {
	name := "Username"
	prompt := &survey.Input{
		Message: "SAML IdP username attribute:",
		Default: "Username",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func PromptSAMLMetadataURL() string {
	name := ""
	prompt := &survey.Input{
		Message: "SAML IdP metadata URL:",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}
