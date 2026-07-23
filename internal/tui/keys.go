package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Retry    key.Binding
	Add      key.Binding
	Open     key.Binding
	Settings key.Binding
	Back     key.Binding
	Quit     key.Binding
	Scroll   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.Open, k.Retry, k.Settings, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Add, k.Retry, k.Settings, k.Quit},
	}
}

var keys = keyMap{
	Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Retry:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "retry")),
	Add:      key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add URL")),
	Open:     key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open folder")),
	Settings: key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "settings")),
	Back:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Scroll:   key.NewBinding(key.WithKeys("pgup", "pgdown"), key.WithHelp("pgup/dn", "scroll")),
}
