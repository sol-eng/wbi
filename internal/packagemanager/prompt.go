package packagemanager

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/system"
	"github.com/sol-eng/wbi/internal/workbench"
)

// Prompt users if they wish to add a default Posit Package Manager URL to Workbench
func PromptPackageManagerChoice() (bool, error) {
	name := true
	messageText := "Would you like to setup Posit Package Manager as the default R and/or Python repo in Workbench? You will need connectivity to the Package Manager server to use this option."
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Posit Package Manager choice prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

func InteractivePackageManagerPrompts(osType config.OperatingSystem) error {
	// prompt for which languages to setup
	languageChoices, err := PromptLanguageRepos()
	if err != nil {
		return fmt.Errorf("issue in prompt for Posit Package Manager language choices: %w", err)
	}

	var overallSkip bool

	var goodURL bool
	var cleanURL string
	for {
		// prompt for base URL
		rawPackageManagerURL, err := PromptPackageManagerURL()
		if err != nil {
			return fmt.Errorf("issue entering Posit Package Manager URL: %w", err)
		}
		if strings.Contains(rawPackageManagerURL, "skip") {
			overallSkip = true
			break
		}
		// verify URL
		cleanURL, err = VerifyPackageManagerURL(rawPackageManagerURL)
		if err != nil {
			if !(strings.Contains(err.Error(), "error in HTTP status code") || strings.Contains(err.Error(), "error retrieving JSON data")) {
				return fmt.Errorf("issue verifying Posit Package Manager URL: %w", err)
			}
		} else {
			goodURL = true
		}
		if goodURL {
			break
		} else {
			system.PrintAndLogInfo(`The URL you entered is not valid. Please try again. To skip this section type "skip".`)
		}
	}

	// r repo
	if lo.Contains(languageChoices, "r") && !overallSkip {
		var goodRepoR bool
		var repoPackageManager string
		for {
			repoPackageManager, err = PromptPackageManagerRepo("r")

			if err != nil {
				return fmt.Errorf("issue entering Posit Package Manager repo name: %w", err)
			}
			if strings.Contains(repoPackageManager, "skip") {
				overallSkip = true
				break
			}

			err = VerifyPackageManagerRepo(cleanURL, repoPackageManager, "r")
			if err != nil {
				if !(strings.Contains(err.Error(), "error in HTTP status code") || strings.Contains(err.Error(), "error retrieving JSON data") || strings.Contains(err.Error(), "error finding the "+repoPackageManager+" repository in Posit Package Manager")) {
					return fmt.Errorf("issue verifying Posit Package Manager repo: %w", err)
				}
			} else {
				goodRepoR = true
			}

			if goodRepoR {
				break
			} else {
				system.PrintAndLogInfo(`The repo you entered is not valid. Please try again. To skip this section type "skip".`)
			}
		}

		if !overallSkip {
			packageManagerURLFull, err := BuildPackagemanagerFullURL(cleanURL, repoPackageManager, osType, "r")
			if err != nil {
				return fmt.Errorf("issue building Posit Package Manager URL: %w", err)
			}

			err = workbench.WriteRepoConfig(packageManagerURLFull, "cran")
			if err != nil {
				if strings.Contains(err.Error(), "line already exists in repos.conf") {
					system.PrintAndLogInfo("CRAN repo already exists in /etc/rstudio/repos.conf. Skipping writing to the file.")
				} else {
					return fmt.Errorf("failed to write CRAN repo config: %w", err)
				}
			}
		}
	}

	// python repo
	if lo.Contains(languageChoices, "python") && !overallSkip {
		var goodRepoPython bool
		var repoPackageManagerPython string
		for {
			repoPackageManagerPython, err = PromptPackageManagerRepo("python")
			if err != nil {
				return fmt.Errorf("issue entering Posit Package Manager repo name: %w", err)
			}
			if strings.Contains(repoPackageManagerPython, "skip") {
				overallSkip = true
				break
			}

			err = VerifyPackageManagerRepo(cleanURL, repoPackageManagerPython, "python")
			if err != nil {
				if !(strings.Contains(err.Error(), "error in HTTP status code") || strings.Contains(err.Error(), "error retrieving JSON data") || strings.Contains(err.Error(), "error finding the "+repoPackageManagerPython+" repository in Posit Package Manager")) {
					return fmt.Errorf("issue verifying Posit Package Manager repo: %w", err)
				}
			} else {
				goodRepoPython = true
			}

			if goodRepoPython {
				break
			} else {
				system.PrintAndLogInfo(`The repo you entered is not valid. Please try again. To skip this section type "skip".`)
			}
		}

		if !overallSkip {
			packageManagerURLFull, err := BuildPackagemanagerFullURL(cleanURL, repoPackageManagerPython, osType, "python")
			if err != nil {
				return fmt.Errorf("issue building Posit Package Manager URL: %w", err)
			}
			err = workbench.WriteRepoConfig(packageManagerURLFull, "pypi")
			if err != nil {
				if strings.Contains(err.Error(), "line already exists in pip.conf") {
					system.PrintAndLogInfo("PyPI URL already exists in /etc/pip.conf. Skipping writing to the file.")
				} else {
					return fmt.Errorf("failed to write PyPI repo config: %w", err)
				}
			}
		}
	}
	return nil
}

func VerifyAndBuildPublicPackageManager(osType config.OperatingSystem) error {
	publicPackageManagerURL := "https://packagemanager.posit.co"

	// verify URL
	_, err := VerifyPackageManagerURL(publicPackageManagerURL)
	if err != nil {
		return fmt.Errorf("issue verifying Posit Package Manager URL: %w", err)
	}

	err = VerifyPackageManagerRepo(publicPackageManagerURL, "cran", "r")
	if err != nil {
		return fmt.Errorf("issue verifying Posit Package Manager repo: %w", err)
	}

	packageManagerURLFull, err := BuildPackagemanagerFullURL(publicPackageManagerURL, "cran", osType, "r")
	if err != nil {
		return fmt.Errorf("issue building Posit Public Package Manager URL: %w", err)
	}
	err = workbench.WriteRepoConfig(packageManagerURLFull, "cran")
	if err != nil {
		return fmt.Errorf("failed to write CRAN repo config: %w", err)
	}
	return nil
}

// Prompt users for a default Posit Package Manager URL
func PromptPackageManagerURL() (string, error) {
	target := ""
	messageText := "Enter your Posit Package Manager base URL (for example, https://exampleaddress.com):"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a Posit Package Manager URL: %w", err)
	}
	log.Info(messageText)
	log.Info(target)
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
	messageText := "Enter the name of your " + languageTitle + " repository on Posit Package Manager (for example, " + exampleRepo + ") :"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a Posit Package Manager "+languageTitle+" repo: %w", err)
	}
	log.Info(messageText)
	log.Info(target)
	return target, nil
}

// Prompt users if they wish to add Posit Public Package Manager as the default R repo in Workbench
func PromptPublicPackageManagerChoice() (bool, error) {
	name := true
	messageText := "Would you like to setup Posit Public Package Manager as the default R repo in Workbench? You will need internet accessibility to use this option."
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Posit Public Package Manager R choice prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
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
	messageText := "What language repositories would you like to setup?"
	var qs = []*survey.Question{
		{
			Name: "languages",
			Prompt: &survey.MultiSelect{
				Message: messageText,
				Options: []string{"r", "python"},
				Default: []string{"r", "python"},
			},
		},
	}
	languageAnswers := struct {
		Languages []string `survey:"languages"`
	}{}
	err := survey.Ask(qs, &languageAnswers, survey.WithRemoveSelectAll(), survey.WithRemoveSelectNone())
	if err != nil {
		return []string{}, errors.New("there was an issue with the repo languages prompt")
	}
	log.Info(messageText)
	log.Info(strings.Join(languageAnswers.Languages, ", "))
	return languageAnswers.Languages, nil
}
