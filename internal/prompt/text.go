package prompt

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

type (
	errMsgText error
)

type modelText struct {
	text      string
	textInput textinput.Model
	err       error
}

func initialModelText(text string, defaultText string) modelText {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 500
	ti.SetValue(defaultText)

	return modelText{
		text:      text,
		textInput: ti,
		err:       nil,
	}
}

func (m modelText) Init() tea.Cmd {
	return textinput.Blink
}

func (m modelText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsgText:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m modelText) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.text,
		m.textInput.View(),
		"(esc or ctrl+c to quit)",
	) + "\n\n"
}

// Generic prompt for text input with prompt text and answer logged
func PromptText(promptText string) (string, error) {
	displayTextFull := fmt.Sprintf("%s:", promptText)

	p := tea.NewProgram(initialModelText(displayTextFull, ""))
	m, err := p.Run()
	if err != nil {
		return "", errors.New("issue occured with the text prompt")
	}

	log.Info(displayTextFull)
	log.Info(m.(modelText).textInput.Value())
	return m.(modelText).textInput.Value(), nil
}
