package commitstable

import (
	"math/rand"
	"time"

	"project-void/internal/ui/styles"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Align(lipgloss.Center)
	focusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(styles.HighlightColor).
			Align(lipgloss.Center)
)

type LoadingState int

const (
	LoadingIdle LoadingState = iota
	LoadingInProgress
	LoadingComplete
	LoadingError
)

type Model struct {
	table         table.Model
	styles        table.Styles
	width         int
	height        int
	borderFocused bool
	loadingState  LoadingState
	progress      progress.Model
	loadError     string
	spinner       *spinner.Model
}

type LoadCommitsProgressMsg struct {
	Percent float64
}

type tickMsg time.Time

type LoadingCompleteMsg struct{}

func getCommitTableColumns(width int) []table.Column {
	branchWidth := 12
	authorWidth := 20
	dateWidth := 12
	repoWidth := 15
	numColumns := 5
	messageWidth := width - branchWidth - authorWidth - dateWidth - repoWidth - 10 - (numColumns - 1)
	if messageWidth < 20 {
		messageWidth = 20
	}
	return []table.Column{
		{Title: "Repo", Width: repoWidth},
		{Title: "Branch", Width: branchWidth},
		{Title: "Author", Width: authorWidth},
		{Title: "Date", Width: dateWidth},
		{Title: "Message", Width: messageWidth},
	}
}

func InitialModel() Model {
	columns := getCommitTableColumns(94)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(6),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		table:         t,
		styles:        s,
		borderFocused: true,
		loadingState:  LoadingIdle,
		progress:      p,
	}
}

func (m Model) Init() tea.Cmd {
	if m.loadingState == LoadingInProgress {
		return tickCmd()
	}
	return nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200+time.Duration(rand.Intn(300))*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) GetSelectedCommit() table.Row {
	return m.table.SelectedRow()
}

func (m Model) TotalCommits() int {
	return len(m.table.Rows())
}

func (m *Model) Focus() {
	m.table.Focus()
}

func (m *Model) Blur() {
	m.table.Blur()
}

func (m *Model) SetFocusedStyle() {
	m.styles.Header = m.styles.Header.BorderForeground(styles.HighlightColor)
	m.styles.Selected = m.styles.Selected.
		BorderForeground(styles.HighlightColor).
		Foreground(lipgloss.Color("15")).
		Background(styles.HighlightColor)
	m.table.SetStyles(m.styles)
	m.borderFocused = true
}

func (m *Model) SetBlurredStyle() {
	m.styles.Header = m.styles.Header.BorderForeground(lipgloss.Color("240"))
	m.styles.Selected = m.styles.Selected.
		BorderForeground(lipgloss.Color("57")).
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.NoColor{})
	m.table.SetStyles(m.styles)
	m.borderFocused = false
}

func (m *Model) IsLoading() bool {
	return m.loadingState == LoadingInProgress
}

func (m *Model) StartLoading() {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)
	m.loadError = ""
}

func (m *Model) StartLoadingWithCmd() tea.Cmd {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)
	m.loadError = ""
	return tickCmd()
}

func (m *Model) UpdateProgress(percent float64) tea.Cmd {
	if m.loadingState == LoadingInProgress {
		return m.progress.SetPercent(percent)
	}
	return nil
}

func (m *Model) SetRows(rows []table.Row) {
	m.table.SetRows(rows)
}

func (m *Model) SetSpinner(s *spinner.Model) {
	m.spinner = s
}
