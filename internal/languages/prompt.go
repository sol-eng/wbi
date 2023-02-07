package languages

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
)

// Prompt asking users which languages they will use
func PromptAndRespond() ([]string, error) {
	var qs = []*survey.Question{
		{
			Name: "languages",
			Prompt: &survey.MultiSelect{
				Message: "What languages will you use",
				Options: []string{"R", "python"},
				Default: []string{"R", "python"},
			},
		},
	}
	languageAnswers := struct {
		Languages []string `survey:"languages"`
	}{}
	err := survey.Ask(qs, &languageAnswers)
	if err != nil {
		return []string{}, errors.New("there was an issue with the languages prompt")
	}
	if !lo.Contains(languageAnswers.Languages, "R") {
		return []string{}, errors.New("R must be a select language to install Workbench")
	}
	return languageAnswers.Languages, nil
}
