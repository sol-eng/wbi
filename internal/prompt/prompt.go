package prompt

import (
	"errors"
	"fmt"
	"strings"

	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

// Generic prompt for yes/no questions with prompt text and answer logged
func PromptConfirm(promptText string) (bool, error) {
	result, err := pterm.DefaultInteractiveConfirm.WithDefaultText(promptText).Show()
	if err != nil {
		return false, errors.New("issue occured with the confirm prompt")
	}

	log.Info(promptText)
	log.Info(fmt.Sprintf("%v", result))
	return result, nil
}

// Generic prompt for text input with prompt text and answer logged
func PromptText(promptText string) (string, error) {
	result, err := pterm.DefaultInteractiveTextInput.WithDefaultText(promptText).Show()
	if err != nil {
		return "", errors.New("issue occured with the text prompt")
	}

	log.Info(promptText)
	log.Info(result)
	return result, nil
}

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
