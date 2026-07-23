// Package tui provides an opt-in interactive queue view for downloads,
// driven by downloader.ProgressReporter events.
package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sekiju/mdl/downloader"
)

type chapterState struct {
	url         string
	chapterID   string
	totalPages  int
	donePages   int
	failedPages int
	done        bool
	err         error
	errMsg      string
	start       time.Time
	duration    time.Duration
}

type (
	startedMsg struct{ url string }
	pageMsg    struct {
		chapterID string
		index     uint
		total     int
		err       error
	}
	chapterErrMsg struct {
		url, chapterID, msg string
		err                 error
	}
	doneMsg struct {
		chapterID string
		duration  time.Duration
	}
)

// reporter implements downloader.ProgressReporter and forwards every event
// to a running tea.Program as a tea.Msg.
type reporter struct {
	program *tea.Program
}

func (r *reporter) ChapterStarted(url string) {
	r.program.Send(startedMsg{url: url})
}

func (r *reporter) ChapterError(url, chapterID, msg string, err error) {
	r.program.Send(chapterErrMsg{url: url, chapterID: chapterID, msg: msg, err: err})
}

func (r *reporter) PageDownloaded(chapterID string, index uint, total int, err error) {
	r.program.Send(pageMsg{chapterID: chapterID, index: index, total: total, err: err})
}

func (r *reporter) ChapterDone(chapterID string, duration time.Duration) {
	r.program.Send(doneMsg{chapterID: chapterID, duration: duration})
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	pendStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	dimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type model struct {
	cancel   context.CancelFunc
	order    []string
	states   map[string]*chapterState
	quitting bool
	total    int
	finished int
}

func newModel(cancel context.CancelFunc, total int) *model {
	return &model{cancel: cancel, states: make(map[string]*chapterState), total: total}
}

func (m *model) Init() tea.Cmd { return nil }

// stateFor resolves the chapterState an event belongs to, keyed by chapterID
// once known and by URL beforehand. If a chapterID arrives with no existing
// entry and exactly one pending (chapterID-less) entry exists, that entry
// adopts the chapterID.
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
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
	}
	if m.total > 0 && m.finished >= m.total {
		return m, tea.Quit
	}
	return m, nil
}

func (m *model) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("mdl - download queue"))
	b.WriteString("\n\n")

	if len(m.order) == 0 {
		b.WriteString(dimStyle.Render("waiting for downloads..."))
		b.WriteString("\n")
	}

	for _, key := range m.order {
		s := m.states[key]
		if s == nil {
			continue
		}
		label := s.chapterID
		if label == "" {
			label = s.url
		}

		switch {
		case s.err != nil:
			b.WriteString(errStyle.Render(fmt.Sprintf("x %s - %s: %v", label, s.errMsg, s.err)))
		case s.done:
			b.WriteString(okStyle.Render(fmt.Sprintf("v %s - done in %s", label, s.duration.Round(time.Millisecond))))
		case s.totalPages > 0:
			b.WriteString(pendStyle.Render(fmt.Sprintf("~ %s - %d/%d pages (%d failed)", label, s.donePages, s.totalPages, s.failedPages)))
		default:
			b.WriteString(pendStyle.Render(fmt.Sprintf("~ %s - downloading...", label)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("ctrl+c/q: cancel and exit"))
	b.WriteString("\n")

	return b.String()
}

// Run starts the interactive queue view for the given URLs, driving a
// downloader through ctx. It blocks until the TUI exits, either because all
// chapters finished or the user cancelled. cancel is invoked on cancel
// keypress so in-flight downloads can clean up via ctx.
func Run(ctx context.Context, cancel context.CancelFunc, urls []string) error {
	m := newModel(cancel, len(urls))
	program := tea.NewProgram(m)

	loader := downloader.NewDownloader(ctx, &reporter{program: program})
	for _, u := range urls {
		loader.Queue(u)
	}

	done := make(chan struct{})
	go func() {
		loader.Stop()
		close(done)
	}()

	if _, err := program.Run(); err != nil {
		return err
	}

	<-done
	return nil
}
