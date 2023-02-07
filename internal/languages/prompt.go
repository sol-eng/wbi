package languages

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
)

func PromptAndRespond() []string {
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
		log.Fatal(err)
	}
	if !lo.Contains(languageAnswers.Languages, "R") {
		log.Fatal("R must be a selected language")
	}
	fmt.Println("You just chose the languages: ", languageAnswers.Languages)
	return languageAnswers.Languages
}
