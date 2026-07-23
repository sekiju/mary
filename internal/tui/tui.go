package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sekiju/mdl/downloader"
)

type reporter struct {
	program *tea.Program
}

func (r *reporter) ChapterStarted(url string) {
	if r.program != nil {
		r.program.Send(startedMsg{url: url})
	}
}

func (r *reporter) ChapterError(url, chapterID, msg string, err error) {
	if r.program != nil {
		r.program.Send(chapterErrMsg{url: url, chapterID: chapterID, msg: msg, err: err})
	}
}

func (r *reporter) PageDownloaded(chapterID string, index uint, total int, err error) {
	if r.program != nil {
		r.program.Send(pageMsg{chapterID: chapterID, index: index, total: total, err: err})
	}
}

func (r *reporter) ChapterDone(chapterID string, duration time.Duration) {
	if r.program != nil {
		r.program.Send(doneMsg{chapterID: chapterID, duration: duration})
	}
}

func newModel(ctx context.Context, cancel context.CancelFunc, version string, statusMsgs []string) *model {
	ti := textinput.New()
	ti.Placeholder = "https://example.com/chapter-1"
	ti.Focus()
	ti.CharLimit = 2048
	ti.Width = 60

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = pendStyle

	return &model{
		mode:        modeInput,
		textInput:   ti,
		spinner:     s,
		progress:    progress.New(progress.WithoutPercentage()),
		help:        help.New(),
		ctx:         ctx,
		cancel:      cancel,
		states:      make(map[string]*chapterState),
		version:     version,
		statusLine:  strings.Join(statusMsgs, " · "),
		statusMsgs:  statusMsgs,
		statusIndex: -1,
	}
}

func (m *model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, textinput.Blink, m.spinner.Tick)
	if len(m.statusMsgs) > 0 {
		cmds = append(cmds, m.nextStatus)
	}
	return tea.Batch(cmds...)
}

func (m *model) nextStatus() tea.Msg {
	m.statusIndex++
	if m.statusIndex < len(m.statusMsgs) {
		return statusMsg{text: m.statusMsgs[m.statusIndex]}
	}
	return nil
}

func (m *model) startDownloader() {
	if m.downloader != nil {
		return
	}
	m.downloader = downloader.NewDownloader(m.ctx, m.reporter)
}

func (m *model) queueURLs(urls []string) {
	m.startDownloader()
	for _, u := range urls {
		if u == "" {
			continue
		}
		m.downloader.Queue(u)
		s := &chapterState{url: u, start: time.Now()}
		m.states[u] = s
		m.order = append(m.order, u)
	}
	m.total = len(m.order)
}

func (m *model) stateFor(url, chapterID string) *chapterState {
	if chapterID != "" {
		if s, ok := m.states[chapterID]; ok {
			return s
		}
		var pending *chapterState
		pendingCount := 0
		for _, s := range m.states {
			if s.chapterID == "" {
				pending = s
				pendingCount++
			}
		}
		if pendingCount == 1 {
			pending.chapterID = chapterID
			m.states[chapterID] = pending
			return pending
		}
		s := &chapterState{url: url, chapterID: chapterID, start: time.Now()}
		m.states[chapterID] = s
		m.order = append(m.order, chapterID)
		return s
	}

	if s, ok := m.states[url]; ok {
		return s
	}
	s := &chapterState{url: url, start: time.Now()}
	m.states[url] = s
	m.order = append(m.order, url)
	return s
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-5)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 5
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		switch m.mode {
		case modeInput:
			cmds = append(cmds, m.handleInputKey(msg))
		case modeQueue:
			cmds = append(cmds, m.handleQueueKey(msg))
		}

	case startedMsg:
		m.stateFor(msg.url, "").start = time.Now()

	case pageMsg:
		s := m.stateFor("", msg.chapterID)
		s.totalPages = msg.total
		if msg.err != nil {
			s.failedPages++
		} else {
			s.donePages++
		}

	case chapterErrMsg:
		s := m.stateFor(msg.url, msg.chapterID)
		if !s.done {
			s.done = true
			s.err = msg.err
			s.errMsg = msg.msg
			m.finished++
		}

	case doneMsg:
		s := m.stateFor("", msg.chapterID)
		if !s.done {
			s.done = true
			s.duration = msg.duration
			m.finished++
		}

	case statusMsg:
		m.statusLine = msg.text
		cmds = append(cmds, m.nextStatus)
	}

	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	if m.total > 0 && m.finished >= m.total {
		return m, tea.Quit
	}
	return m, tea.Batch(cmds...)
}

func (m *model) handleInputKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		if m.cancel != nil {
			m.cancel()
		}
		if m.downloader != nil {
			go m.downloader.Stop()
		}
		return tea.Quit

	case "enter":
		raw := strings.TrimSpace(m.textInput.Value())
		if raw == "" {
			return nil
		}
		urls := strings.Split(raw, " ")
		m.queueURLs(urls)
		m.mode = modeQueue
		m.textInput.Reset()
		m.textInput.Blur()
		m.selected = 0
		return nil

	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return cmd
	}
}

func (m *model) handleQueueKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		if m.cancel != nil {
			m.cancel()
		}
		if m.downloader != nil {
			go m.downloader.Stop()
		}
		return tea.Quit

	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}

	case "down", "j":
		if m.selected < len(m.order)-1 {
			m.selected++
		}

	case "r":
		if m.selected >= 0 && m.selected < len(m.order) {
			key := m.order[m.selected]
			s := m.states[key]
			if s != nil && s.err != nil {
				s.done = false
				s.err = nil
				s.errMsg = ""
				s.donePages = 0
				s.failedPages = 0
				s.totalPages = 0
				s.start = time.Now()
				m.finished--
				m.downloader.Queue(s.url)
			}
		}

	case "a":
		m.mode = modeInput
		m.textInput.Focus()
		m.textInput.Reset()
		return textinput.Blink
	}

	return nil
}

func (m *model) renderRow(s *chapterState, selected bool) string {
	label := s.chapterID
	if label == "" {
		label = s.url
	}

	prefix := "  "
	if selected {
		prefix = "> "
	}

	switch {
	case s.err != nil:
		return errStyle.Render(fmt.Sprintf("%sx %s — %s: %v", prefix, label, s.errMsg, s.err))
	case s.done:
		return okStyle.Render(fmt.Sprintf("%sv %s — done in %s", prefix, label, s.duration.Round(time.Millisecond)))
	case s.totalPages > 0:
		pct := float64(s.donePages) / float64(s.totalPages)
		bar := m.progress.ViewAs(pct)
		return pendStyle.Render(fmt.Sprintf("%s~ %s", prefix, label)) + " " + bar +
			dimStyle.Render(fmt.Sprintf(" %d/%d (%d failed)", s.donePages, s.totalPages, s.failedPages))
	default:
		return pendStyle.Render(fmt.Sprintf("%s%s %s — waiting...", prefix, m.spinner.View(), label))
	}
}

func (m *model) renderDetail() string {
	if m.selected < 0 || m.selected >= len(m.order) {
		return ""
	}
	s := m.states[m.order[m.selected]]
	if s == nil {
		return ""
	}
	var b strings.Builder
	label := s.chapterID
	if label == "" {
		label = s.url
	}
	b.WriteString(dimStyle.Render("── detail ──────────────────────────────"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("ID:       %s\n", label))
	b.WriteString(fmt.Sprintf("URL:      %s\n", s.url))
	b.WriteString(fmt.Sprintf("Pages:    %d/%d (%d failed)\n", s.donePages, s.totalPages, s.failedPages))
	if s.done {
		b.WriteString(fmt.Sprintf("Duration: %s\n", s.duration.Round(time.Millisecond)))
	}
	if s.err != nil {
		b.WriteString(errStyle.Render(fmt.Sprintf("Error:    %s: %v\n", s.errMsg, s.err)))
	}
	b.WriteString(dimStyle.Render("──────────────────────────────────────────"))
	return b.String()
}

func (m *model) buildContent() string {
	var b strings.Builder
	for i, key := range m.order {
		s := m.states[key]
		if s == nil {
			continue
		}
		b.WriteString(m.renderRow(s, i == m.selected))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(m.renderDetail())
	return b.String()
}

func (m *model) header() string {
	mdlStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4"))
	return mdlStyle.Render("mdl") + dimStyle.Render(" v"+m.version)
}

func (m *model) View() string {
	if !m.ready {
		return "initializing..."
	}

	switch m.mode {
	case modeInput:
		return m.viewInput()
	case modeQueue:
		return m.viewQueue()
	default:
		return ""
	}
}

func (m *model) statusBar() string {
	if m.statusLine == "" {
		return ""
	}
	return dimStyle.Render(m.statusLine)
}

func (m *model) viewInput() string {
	var b strings.Builder
	b.WriteString(m.header())
	b.WriteString("\n\n")
	b.WriteString("Enter one or more chapter URLs (space-separated):\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render("enter: confirm • q/ctrl+c: quit"))
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Margin(1, 2).Render(b.String()),
		m.statusBar(),
	)
}

func (m *model) viewQueue() string {
	content := m.buildContent()
	m.viewport.SetContent(content)

	return lipgloss.JoinVertical(lipgloss.Left,
		m.header(),
		m.viewport.View(),
		m.help.View(keys),
		m.statusBar(),
	)
}

func Run(ctx context.Context, cancel context.CancelFunc, urls []string, version string, statusMsgs []string) error {
	m := newModel(ctx, cancel, version, statusMsgs)

	rep := &reporter{}
	m.reporter = rep

	program := tea.NewProgram(m, tea.WithAltScreen())
	rep.program = program

	if len(urls) > 0 {
		m.mode = modeQueue
		loader := downloader.NewDownloader(ctx, rep)
		m.downloader = loader
		for _, u := range urls {
			loader.Queue(u)
			s := &chapterState{url: u, start: time.Now()}
			m.states[u] = s
			m.order = append(m.order, u)
		}
		m.total = len(m.order)

		done := make(chan struct{})
		go func() {
			loader.Stop()
			close(done)
		}()
		defer func() { <-done }()
	}

	if _, err := program.Run(); err != nil {
		return err
	}

	if m.downloader != nil {
		m.downloader.Stop()
	}
	return nil
}
