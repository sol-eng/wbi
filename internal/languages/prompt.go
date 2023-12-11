package languages

import (
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/sol-eng/wbi/internal/prompt"
)

// Prompt asking users which languages they will use
func PromptAndRespond() ([]string, error) {
	promptText := "What languages will you use"

	selectOptions := []string{"R", "python"}

	result, err := prompt.PromptMultiSelect(promptText, selectOptions, selectOptions)
	if err != nil {
		return []string{}, fmt.Errorf("issue occured in languages selection prompt: %w", err)
	}
	if !lo.Contains(result, "R") {
		return []string{}, errors.New("R must be a select language to install Workbench")
	}

	return result, nil
}
