package slacktable

import (
	"fmt"
	"project-void/internal/slack"
	"project-void/internal/ui/styles"
	"strings"
	"time"

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

type Model struct {
	table         table.Model
	styles        table.Styles
	width         int
	height        int
	borderFocused bool
}

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

	return Model{
		table:  t,
		styles: s,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

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
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
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
	messages, err := slackClient.GetMessagesSince(since, config)
	if err != nil {
		return fmt.Errorf("failed to load messages: %w", err)
	}

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

	m.table.SetRows(rows)
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
}
