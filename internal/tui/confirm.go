package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	message   string
	confirmed bool
	done      bool
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "y", "Y":
		m.confirmed = true
		m.done = true
		return m, tea.Quit
	case "n", "N", "esc", "ctrl+c":
		m.confirmed = false
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

func (m confirmModel) View() string {
	if m.done {
		return ""
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	promptStyle := lipgloss.NewStyle().Faint(true)

	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n",
		titleStyle.Render("mdl: invalid config file"),
		m.message,
		promptStyle.Render("Overwrite it with defaults? (y/n)"),
	)

	return lipgloss.NewStyle().Padding(1, 2).Render(content)
}

// Confirm shows a full-screen TUI prompt with message and returns whether the
// user confirmed with "y".
func Confirm(message string) (bool, error) {
	m := confirmModel{message: message}
	program := tea.NewProgram(m, tea.WithAltScreen())

	result, err := program.Run()
	if err != nil {
		return false, err
	}

	final, ok := result.(confirmModel)
	if !ok {
		return false, nil
	}

	return final.confirmed, nil
}
