package prompt

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

const listHeight = 14

var (
	appStyle          = lipgloss.NewStyle().Padding(1, 2)
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// single select model
type modelSingle struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m modelSingle) Init() tea.Cmd {
	return nil
}

func (m modelSingle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelSingle) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(m.choice)
	}
	if m.quitting {
		return quitTextStyle.Render("Nothing chosen")
	}
	return "\n" + m.list.View()
}

// multi select model
type modelMulti struct {
	list     list.Model
	choices  []string
	quitting bool
}

func (m modelMulti) Init() tea.Cmd {
	return nil
}

func (m modelMulti) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "space":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choices = append(m.choices, string(i))
			}
			return m, nil
		case "enter":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelMulti) View() string {
	return appStyle.Render(m.list.View())
}

func stringSliceToListItems(s []string) []list.Item {
	items := make([]list.Item, len(s))
	for i, str := range s {
		items[i] = item(str)
	}
	return items
}

func PromptSingleSelect(promptText string, options []string, defaultOption string) (string, error) {

	items := stringSliceToListItems(options)

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = promptText
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	p := tea.NewProgram(modelSingle{list: l})
	m, err := p.Run()
	if err != nil {
		return "", errors.New("issue occured with the single select prompt")
	}

	log.Info(promptText)
	log.Info(fmt.Sprintf("%v", m.(modelSingle).choice))
	return m.(modelSingle).choice, nil
}

func PromptMultiSelect(promptText string, options []string, defaultOptions []string) ([]string, error) {
	items := stringSliceToListItems(options)

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = promptText
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	p := tea.NewProgram(modelMulti{list: l})
	m, err := p.Run()
	if err != nil {
		return []string{}, errors.New("issue occured with the multi select prompt")
	}

	log.Info(promptText)
	log.Info(fmt.Sprintf("%v", strings.Join(m.(modelMulti).choices, "\n")))
	return m.(modelMulti).choices, nil
}
