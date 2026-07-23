package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knst0/mdl/config"
	"github.com/knst0/mdl/extractor/registry"
)

type settingsRow struct {
	hostname    string
	cookieInput textinput.Model
	hasCookie   bool
}

func newSettingsRows() []settingsRow {
	hosts := registry.Hostnames()
	rows := make([]settingsRow, 0, len(hosts))
	for _, host := range hosts {
		ti := textinput.New()
		ti.Placeholder = "cookie value"
		ti.CharLimit = 4096
		ti.Width = 60

		var hasCookie bool
		if site, ok := config.Params.File.Sites[host]; ok && site.Cookie != nil {
			ti.SetValue(*site.Cookie)
			hasCookie = true
		}

		rows = append(rows, settingsRow{hostname: host, cookieInput: ti, hasCookie: hasCookie})
	}
	return rows
}

func (m *model) enterSettings(returnMode mode) {
	m.settingsReturnMode = returnMode
	m.settingsRows = newSettingsRows()
	m.settingsSelected = 0
	m.settingsEditing = false
	m.settingsSaved = ""
	m.mode = modeSettings
}

func (m *model) handleSettingsKey(msg tea.KeyMsg) tea.Cmd {
	if len(m.settingsRows) == 0 {
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.mode = m.settingsReturnMode
		}
		return nil
	}

	row := &m.settingsRows[m.settingsSelected]

	if m.settingsEditing {
		return m.handleSettingsEditKey(msg, row)
	}

	return m.handleSettingsBrowseKey(msg, row)
}

func (m *model) handleSettingsEditKey(msg tea.KeyMsg, row *settingsRow) tea.Cmd {
	switch msg.String() {
	case "esc":
		if row.hasCookie {
			if site, ok := config.Params.File.Sites[row.hostname]; ok && site.Cookie != nil {
				row.cookieInput.SetValue(*site.Cookie)
			}
		} else {
			row.cookieInput.SetValue("")
		}
		row.cookieInput.Blur()
		m.settingsEditing = false
		return nil

	case "enter":
		val := strings.TrimSpace(row.cookieInput.Value())
		if val == "" {
			config.SetSiteCookie(row.hostname, nil)
			row.hasCookie = false
		} else {
			config.SetSiteCookie(row.hostname, &val)
			row.hasCookie = true
		}
		if err := config.Save(""); err != nil {
			m.settingsSaved = errStyle.Render(fmt.Sprintf("save failed: %v", err))
		} else {
			m.settingsSaved = okStyle.Render("saved")
		}
		row.cookieInput.Blur()
		m.settingsEditing = false
		return nil

	default:
		var cmd tea.Cmd
		row.cookieInput, cmd = row.cookieInput.Update(msg)
		return cmd
	}
}

func (m *model) handleSettingsBrowseKey(msg tea.KeyMsg, row *settingsRow) tea.Cmd {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		if m.cancel != nil {
			m.cancel()
		}
		if m.downloader != nil {
			go m.downloader.Stop()
		}
		return tea.Quit

	case "esc", "q":
		m.mode = m.settingsReturnMode
		return nil

	case "up", "k":
		if m.settingsSelected > 0 {
			m.settingsSelected--
		}

	case "down", "j":
		if m.settingsSelected < len(m.settingsRows)-1 {
			m.settingsSelected++
		}

	case "enter":
		m.settingsEditing = true
		m.settingsSaved = ""
		return row.cookieInput.Focus()

	case "d":
		config.SetSiteCookie(row.hostname, nil)
		row.hasCookie = false
		row.cookieInput.SetValue("")
		if err := config.Save(""); err != nil {
			m.settingsSaved = errStyle.Render(fmt.Sprintf("save failed: %v", err))
		} else {
			m.settingsSaved = okStyle.Render("cleared")
		}

	case "o":
		if err := openURL("https://" + row.hostname); err != nil {
			m.settingsSaved = errStyle.Render(fmt.Sprintf("open failed: %v", err))
		}
	}

	return nil
}

// hyperlink wraps text in an OSC 8 escape sequence so terminals that support
// it (Windows Terminal, iTerm2, VS Code, etc.) render it as a clickable link
// without mdl handling any mouse events itself.
func hyperlink(url, text string) string {
	return "\x1b]8;;" + url + "\x1b\\" + text + "\x1b]8;;\x1b\\"
}

func (m *model) buildSettingsContent() string {
	var b strings.Builder

	if len(m.settingsRows) == 0 {
		b.WriteString(dimStyle.Render("no extractors registered"))
		return b.String()
	}

	for i, row := range m.settingsRows {
		prefix := "  "
		if i == m.settingsSelected {
			prefix = "> "
		}

		status := dimStyle.Render("[not set]")
		if row.hasCookie {
			status = okStyle.Render("[cookie set]")
		}

		link := hyperlink("https://"+row.hostname, fmt.Sprintf("%-30s", row.hostname))
		line := fmt.Sprintf("%s%s %s", prefix, link, status)
		b.WriteString(line)
		b.WriteString("\n")

		if i == m.settingsSelected && m.settingsEditing {
			b.WriteString("    " + row.cookieInput.View())
			b.WriteString("\n")
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

func (m *model) viewSettings() string {
	content := m.buildSettingsContent()

	maxHeight := m.height - 8
	if maxHeight < 1 {
		maxHeight = 1
	}
	contentHeight := strings.Count(content, "\n") + 1
	if contentHeight > maxHeight {
		contentHeight = maxHeight
	}

	m.viewport.SetContent(content)
	m.viewport.Height = contentHeight

	footer := dimStyle.Render("enter: edit cookie • o: open in browser • d: clear cookie • esc/q: back")
	if m.settingsEditing {
		footer = dimStyle.Render("enter: save • esc: cancel")
	}

	indent := lipgloss.NewStyle().MarginLeft(2)

	return lipgloss.JoinVertical(lipgloss.Left,
		"\n"+indent.Render(m.header()+"\n\nSite cookies:"),
		indent.Render(m.viewport.View()),
		indent.Render(m.settingsSaved),
		indent.Render(footer),
	)
}
