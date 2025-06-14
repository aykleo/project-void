package jiratable

import (
	"fmt"
	"math/rand"
	"project-void/internal/jira"
	"project-void/internal/ui/styles"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
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
}

type LoadIssuesProgressMsg struct {
	Percent float64
}

type tickMsg time.Time

type LoadingCompleteMsg struct{}

func InitialModel() Model {
	columns := []table.Column{
		{Title: "Issue", Width: 12},
		{Title: "Status", Width: 15},
		{Title: "Action", Width: 12},
		{Title: "Date", Width: 10},
		{Title: "Summary", Width: 40},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(8),
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
		table:        t,
		styles:       s,
		loadingState: LoadingIdle,
		progress:     p,
	}
}

func (m Model) Init() tea.Cmd {
	if m.loadingState == LoadingInProgress {
		return tickCmd()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.progress.Width = msg.Width - 20
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}

		tableHeight := m.height - 4
		if tableHeight < 8 {
			tableHeight = 8
		}
		m.table.SetHeight(tableHeight)

		if m.width > 0 {

			issueWidth := 12
			statusWidth := 15
			actionWidth := 12
			dateWidth := 10

			usedWidth := issueWidth + statusWidth + actionWidth + dateWidth + 10
			summaryWidth := m.width - usedWidth

			if summaryWidth < 20 {
				summaryWidth = 20
			}
			if issueWidth < 8 {
				issueWidth = 8
			}
			if statusWidth < 10 {
				statusWidth = 10
			}
			if actionWidth < 8 {
				actionWidth = 8
			}
			if dateWidth < 8 {
				dateWidth = 8
			}

			columns := []table.Column{
				{Title: "Issue", Width: issueWidth},
				{Title: "Status", Width: statusWidth},
				{Title: "Action", Width: actionWidth},
				{Title: "Date", Width: dateWidth},
				{Title: "Summary", Width: summaryWidth},
			}
			m.table.SetColumns(columns)
		}

	case LoadIssuesProgressMsg:
		if m.loadingState == LoadingInProgress {
			cmd = m.progress.SetPercent(msg.Percent)
			return m, cmd
		}

	case tickMsg:
		if m.loadingState == LoadingInProgress {
			currentPercent := m.progress.Percent()

			if currentPercent < 0.95 {
				maxIncrement := 0.1 * (1.0 - currentPercent)
				increment := rand.Float64() * maxIncrement

				newPercent := currentPercent + increment
				if newPercent > 0.95 {
					newPercent = 0.95
				}

				return m, tea.Batch(tickCmd(), m.progress.SetPercent(newPercent))
			}

			return m, tickCmd()
		}

	case LoadingCompleteMsg:
		if m.loadingState == LoadingInProgress {
			m.loadingState = LoadingComplete
			return m, m.progress.SetPercent(1.0)
		}

	case progress.FrameMsg:
		if m.loadingState == LoadingInProgress || m.progress.Percent() < 1.0 {
			progressModel, cmd := m.progress.Update(msg)
			m.progress = progressModel.(progress.Model)
			return m, cmd
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200+time.Duration(rand.Intn(300))*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) View() string {
	if m.loadingState == LoadingInProgress {
		loadingText := "Loading JIRA issues..."
		progressView := m.progress.View()

		content := lipgloss.JoinVertical(lipgloss.Center,
			loadingText,
			progressView,
		)

		return content
	}

	if m.loadingState == LoadingError {
		errorText := fmt.Sprintf("Error loading JIRA issues: %s", m.loadError)
		content := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errorText)

		return content
	}

	tableView := baseStyle.Render(m.table.View())

	if m.borderFocused {
		tableView = focusedStyle.Render(m.table.View())
	}

	if m.width > 0 {
		return lipgloss.NewStyle().Width(m.width).Render(tableView)
	}

	return tableView
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

func (m *Model) LoadIssues(jiraClient *jira.JiraClient, since time.Time, config *jira.Config) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	m.progress.SetPercent(0.2)
	issues, err := jiraClient.GetIssuesSince(since, config)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load issues: %w", err)
	}

	rows := make([]table.Row, len(issues))
	for i, issue := range issues {
		if len(issues) > 0 {
			progress := 0.2 + (0.7 * float64(i) / float64(len(issues)))
			m.progress.SetPercent(progress)
		}

		action := issue.UserAction
		if len(action) > 11 {
			action = action[:8] + "..."
		}

		status := issue.Status
		if len(status) > 14 {
			status = status[:11] + "..."
		}

		actionDate := issue.ActionDate.Format("2006-01-02")

		summary := strings.ReplaceAll(issue.Summary, "\n", " ")
		summary = strings.TrimSpace(summary)
		maxSummaryLength := 80
		if len(summary) > maxSummaryLength {
			summary = summary[:maxSummaryLength-3] + "..."
		}

		rows[i] = table.Row{
			issue.Key,
			status,
			action,
			actionDate,
			summary,
		}
	}

	m.progress.SetPercent(0.9)
	m.table.SetRows(rows)
	m.progress.SetPercent(1.0)
	m.loadingState = LoadingComplete
	return nil
}

func (m Model) GetSelectedIssue() table.Row {
	return m.table.SelectedRow()
}

func (m Model) TotalIssues() int {
	return len(m.table.Rows())
}

func (m *Model) IsLoading() bool {
	return m.loadingState == LoadingInProgress
}

func (m *Model) StartLoading() {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)
	m.loadError = ""
}

func (m *Model) UpdateProgress(percent float64) tea.Cmd {
	if m.loadingState == LoadingInProgress {
		return m.progress.SetPercent(percent)
	}
	return nil
}
