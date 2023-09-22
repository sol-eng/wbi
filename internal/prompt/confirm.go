package prompt

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

type (
	errMsgConfirm error
)

type modelConfirm struct {
	text      string
	textInput textinput.Model
	err       error
}

func initialModelConfirm(text string, defaultText string) modelConfirm {
	ti := textinput.New()
	ti.Placeholder = defaultText
	ti.Focus()
	ti.CharLimit = 1

	return modelConfirm{
		text:      text,
		textInput: ti,
		err:       nil,
	}
}

func (m modelConfirm) Init() tea.Cmd {
	return textinput.Blink
}

func (m modelConfirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsgConfirm:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m modelConfirm) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.text,
		m.textInput.View(),
		"(esc or ctrl+c to quit)",
	) + "\n\n"
}

// Generic prompt for yes/no questions with prompt text and answer logged
func PromptConfirm(promptText string) (bool, error) {
	displayTextFull := fmt.Sprintf("%s: %s", promptText, "[y/n]")

	p := tea.NewProgram(initialModelConfirm(displayTextFull, "y"))
	m, err := p.Run()
	if err != nil {
		return false, errors.New("issue occured with the confirm prompt")
	}

	log.Info(displayTextFull)
	log.Info(fmt.Sprintf("%v", m.(modelConfirm).textInput.Value()))

	// return true if the user input is "y", false if "n"
	if m.(modelConfirm).textInput.Value() == "y" {
		return true, nil
	} else if m.(modelConfirm).textInput.Value() == "n" {
		return false, nil
	} else {
		return false, errors.New("invalid input")
	}
}
