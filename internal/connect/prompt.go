package connect

import (
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/prompt"
	"github.com/sol-eng/wbi/internal/system"
	"github.com/sol-eng/wbi/internal/workbench"
)

// Prompt users if they wish to add a default Connect URL to Workbench
func PromptConnectChoice() (bool, error) {
	confirmText := "Would you like to provide a default Connect URL for Workbench? You will need connectivity to the Connect server to use this option."

	result, err := prompt.PromptConfirm(confirmText)
	if err != nil {
		return false, fmt.Errorf("issue occured in Connect URL confirm prompt: %w", err)
	}

	return result, nil
}

func PromptVerifyAndConfigConnect() error {
	var overallSkip bool
	var goodURL bool
	var connectURLFull string
	for {
		rawConnectURL, err := PromptConnectURL()
		if err != nil {
			return fmt.Errorf("issue entering Connect URL: %w", err)
		}
		if strings.Contains(rawConnectURL, "skip") {
			overallSkip = true
			break
		}
		connectURLFull, err = VerifyConnectURL(rawConnectURL)
		if err != nil {
			if !(strings.Contains(err.Error(), "error in HTTP status code") || strings.Contains(err.Error(), "error retrieving JSON data")) {
				return fmt.Errorf("issue with checking the Connect URL: %w", err)
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

	if !overallSkip {
		err := workbench.WriteConnectURLConfig(connectURLFull)
		if err != nil {
			return fmt.Errorf("failed to write Connect URL config: %w", err)
		}
	}
	return nil
}

// Prompt users for a default Connect URL
func PromptConnectURL() (string, error) {
	promptText := "Enter a default Connect URL"

	result, err := prompt.PromptText(promptText)
	if err != nil {
		return "", fmt.Errorf("issue occured in Connect URL text prompt: %w", err)
	}

	return result, nil
}
