package tui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/knst0/mdl/downloader"
)

type mode int

const (
	modeInput mode = iota
	modeQueue
	modeSettings
)

type chapterState struct {
	url         string
	chapterID   string
	title       string
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
	titleMsg struct {
		chapterID string
		title     string
	}
	statusMsg struct{ text string }
)

type model struct {
	mode mode

	textInput textinput.Model
	spinner   spinner.Model
	progress  progress.Model
	viewport  viewport.Model
	help      help.Model

	downloader *downloader.Downloader
	reporter   *reporter
	cancel     context.CancelFunc
	ctx        context.Context

	order    []string
	states   map[string]*chapterState
	selected int
	total    int
	finished int

	settingsRows       []settingsRow
	settingsSelected   int
	settingsEditing    bool
	settingsSaved      string
	settingsReturnMode mode

	width  int
	height int

	statusLine  string
	version     string
	statusMsgs  []string
	statusIndex int

	quitting bool
	ready    bool
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	pendStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	dimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
