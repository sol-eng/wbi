package languages

import (
	"errors"
	"strings"

	"github.com/pterm/pterm"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

// Prompt asking users which languages they will use
func PromptAndRespond() ([]string, error) {
	promptText := "What languages will you use"

	selectOptions := []string{"R", "python"}

	result, err := pterm.DefaultInteractiveMultiselect.
		WithDefaultText(promptText).
		WithDefaultOptions(selectOptions).
		WithOptions(selectOptions).
		WithFilter(false).
		Show()
	if err != nil {
		return []string{}, errors.New("there was an issue with the languages prompt")
	}
	if !lo.Contains(result, "R") {
		return []string{}, errors.New("R must be a select language to install Workbench")
	}

	log.Info(promptText)
	log.Info(strings.Join(result, ", "))
	return result, nil
}
