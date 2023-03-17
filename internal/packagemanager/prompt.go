package packagemanager

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sol-eng/wbi/internal/config"
)

// Prompt users if they wish to add a default Posit Package Manager URL to Workbench
func PromptPackageManagerChoice() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to setup Posit Package Manager as the default R and/or Python repo in Workbench? You will need connectivity to the Package Manager server to use this option.",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Posit Package Manager choice prompt")
	}
	return name, nil
}

// Prompt users for a default Posit Package Manager URL
func PromptPackageManagerURL() (string, error) {
	target := ""
	prompt := &survey.Input{
		Message: "Enter your Posit Package Manager base URL (for example, https://exampleaddress.com):",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a Posit Package Manager URL: %w", err)
	}
	return target, nil
}

// Prompt users for a Posit Package Manager repo name
func PromptPackageManagerRepo(language string) (string, error) {
	var exampleRepo string
	if language == "r" {
		exampleRepo = "prod-cran"
	} else if language == "python" {
		exampleRepo = "pypi"
	} else {
		return "", errors.New("language not supported for Posit Package Manager")
	}

	languageTitle := strings.Title(language)

	target := ""
	prompt := &survey.Input{
		Message: "Enter the name of your " + languageTitle + " repository on Posit Package Manager (for example, " + exampleRepo + ") :",
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a Posit Package Manager "+languageTitle+" repo: %w", err)
	}
	return target, nil
}

// Prompt users if they wish to add Posit Public Package Manager as the default R repo in Workbench
func PromptPublicPackageManagerChoice() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to setup Posit Public Package Manager as the default R repo in Workbench? You will need internet accessibility to use this option.",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Posit Public Package Manager R choice prompt")
	}
	return name, nil
}

// PromptPackageManagerNameAndBuildURL prompts users for a Posit Package Manager repo name and builds the full URL
func PromptPackageManagerNameAndBuildURL(cleanURL string, osType config.OperatingSystem, language string) (string, error) {
	repoPackageManager, err := PromptPackageManagerRepo(language)
	if err != nil {
		return "", fmt.Errorf("issue entering Posit Package Manager repo name: %w", err)
	}

	err = VerifyPackageManagerRepo(cleanURL, repoPackageManager, language)
	if err != nil {
		return "", fmt.Errorf("issue with checking the Posit Package Manager repo: %w", err)
	}

	fullRepoURL, err := BuildPackagemanagerFullURL(cleanURL, repoPackageManager, osType, language)
	if err != nil {
		return "", fmt.Errorf("issue with creating the full Posit Package Manager URL: %w", err)
	}
	return fullRepoURL, nil
}

// Prompt asking users which language repos they will use
func PromptLanguageRepos() ([]string, error) {
	var qs = []*survey.Question{
		{
			Name: "languages",
			Prompt: &survey.MultiSelect{
				Message: "What language repositories would you like to setup?",
				Options: []string{"r", "python"},
				Default: []string{"r", "python"},
			},
		},
	}
	languageAnswers := struct {
		Languages []string `survey:"languages"`
	}{}
	err := survey.Ask(qs, &languageAnswers)
	if err != nil {
		return []string{}, errors.New("there was an issue with the repo languages prompt")
	}

	return languageAnswers.Languages, nil
}
