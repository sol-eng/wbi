package languages

import (
	"errors"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

// Prompt asking users which languages they will use
func PromptAndRespond() ([]string, error) {
	messageText := "What languages will you use"
	var qs = []*survey.Question{
		{
			Name: "languages",
			Prompt: &survey.MultiSelect{
				Message: messageText,
				Options: []string{"R", "python"},
				Default: []string{"R", "python"},
			},
		},
	}
	languageAnswers := struct {
		Languages []string `survey:"languages"`
	}{}

	err := survey.Ask(qs, &languageAnswers, survey.WithRemoveSelectAll(), survey.WithRemoveSelectNone())
	if err != nil {
		return []string{}, errors.New("there was an issue with the languages prompt")
	}
	if !lo.Contains(languageAnswers.Languages, "R") {
		return []string{}, errors.New("R must be a select language to install Workbench")
	}
	log.Info(messageText)
	log.Info(strings.Join(languageAnswers.Languages, ", "))
	return languageAnswers.Languages, nil
}
