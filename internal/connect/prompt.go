package connect

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/system"
	"github.com/sol-eng/wbi/internal/workbench"
)

// Prompt users if they wish to add a default Connect URL to Workbench
func PromptConnectChoice() (bool, error) {
	name := true
	messageText := "Would you like to provide a default Connect URL for Workbench? You will need connectivity to the Connect server to use this option."
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Connect URL prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
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
	target := ""
	messageText := "Enter a default Connect URL:"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", fmt.Errorf("issue prompting for a Connect URL: %w", err)
	}
	log.Info(messageText)
	log.Info(target)
	return target, nil
}
