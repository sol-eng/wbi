package prompt

import (
	"errors"
	"strings"

	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

func PromptSingleSelect(promptText string, options []string, defaultOption string, filter bool) (string, error) {
	result, err := pterm.DefaultInteractiveSelect.
		WithDefaultText(promptText).
		WithDefaultOption(defaultOption).
		WithFilter(filter).
		WithOptions(options).
		Show()
	if err != nil {
		return "", errors.New("issue occured with the single select prompt")
	}

	log.Info(promptText)
	log.Info(result)
	return result, nil
}

func PromptMultiSelect(promptText string, options []string, defaultOptions []string, filter bool) ([]string, error) {
	result, err := pterm.DefaultInteractiveMultiselect.
		WithDefaultText(promptText).
		WithDefaultOptions(defaultOptions).
		WithOptions(options).
		WithFilter(filter).
		WithKeyConfirm(keys.Enter).
		WithKeySelect(keys.Space).
		WithCheckmark(&pterm.Checkmark{Checked: "+", Unchecked: " "}).
		Show()
	if err != nil {
		return []string{}, errors.New("issue occured with the multi select prompt")
	}

	log.Info(promptText)
	log.Info(strings.Join(result, ", "))
	return result, nil
}
