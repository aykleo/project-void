package slacktable

import (
	"fmt"
	"math/rand"
	"project-void/internal/slack"
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

type LoadMessagesProgressMsg struct {
	Percent float64
}

type tickMsg time.Time

type LoadingCompleteMsg struct{}

func InitialModel() Model {
	columns := []table.Column{
		{Title: "From (not yet implemented)", Width: 28},
		{Title: "To", Width: 16},
		{Title: "Time", Width: 10},
		{Title: "Message", Width: 35},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(1),
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
		if tableHeight < 1 {
			tableHeight = 1
		}
		m.table.SetHeight(tableHeight)

		if m.width > 0 {
			fromWidth := 28
			toWidth := 16
			timeWidth := 10
			messageWidth := m.width - fromWidth - toWidth - timeWidth - 12
			if messageWidth < 20 {
				messageWidth = 20
			}
			columns := []table.Column{
				{Title: "From (not yet implemented)", Width: fromWidth},
				{Title: "To", Width: toWidth},
				{Title: "Time", Width: timeWidth},
				{Title: "Message", Width: messageWidth},
			}
			m.table.SetColumns(columns)
		}

	case LoadMessagesProgressMsg:
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
		loadingText := "Loading Slack messages..."
		progressView := m.progress.View()

		content := lipgloss.JoinVertical(lipgloss.Center,
			loadingText,
			progressView,
		)

		return content
	}

	if m.loadingState == LoadingError {
		errorText := fmt.Sprintf("Error loading Slack messages: %s", m.loadError)
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

func (m *Model) LoadMessages(slackClient *slack.SlackClient, since time.Time, config *slack.Config) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	m.progress.SetPercent(0.3)
	messages, err := slackClient.GetMessagesSince(since, config)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load messages: %w", err)
	}

	m.progress.SetPercent(0.6)
	var userMessages []slack.Message
	for _, message := range messages {
		if config.FilterByUser && config.UserID != "" {
			if message.UserID == config.UserID {
				userMessages = append(userMessages, message)
			}
		} else {
			userMessages = append(userMessages, message)
		}
	}

	rows := make([]table.Row, len(userMessages))
	for i, message := range userMessages {
		if len(userMessages) > 0 {
			progress := 0.6 + (0.3 * float64(i) / float64(len(userMessages)))
			m.progress.SetPercent(progress)
		}

		fromName := message.User
		if len(fromName) > 11 {
			fromName = fromName[:8] + "..."
		}

		channelName, err := slackClient.GetChannelName(message.Channel)
		if err != nil {
			channelName = message.Channel
		}

		toName := channelName
		if len(toName) > 15 {
			toName = toName[:12] + "..."
		}

		timestamp := message.Timestamp.Format("15:04")

		text := strings.ReplaceAll(message.Text, "\n", " ")
		text = strings.TrimSpace(text)
		maxTextLength := 60
		if len(text) > maxTextLength {
			text = text[:maxTextLength-3] + "..."
		}

		rows[i] = table.Row{
			fromName,
			toName,
			timestamp,
			text,
		}
	}

	m.progress.SetPercent(0.9)
	m.table.SetRows(rows)
	m.progress.SetPercent(1.0)
	m.loadingState = LoadingComplete
	return nil
}

func (m Model) GetSelectedMessage() table.Row {
	return m.table.SelectedRow()
}

func (m Model) TotalMessages() int {
	return len(m.table.Rows())
}

func (m *Model) SetPlaceholder() {
	placeholderRow := table.Row{
		"Slack",
		"Coming Soon",
		"--:--",
		"Slack integration is not yet implemented",
	}

	m.table.SetRows([]table.Row{placeholderRow})
	m.loadingState = LoadingComplete
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
