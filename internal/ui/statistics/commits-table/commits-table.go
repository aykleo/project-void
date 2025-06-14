package commitstable

import (
	"fmt"
	"project-void/internal/git"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	Align(lipgloss.Center)

type Model struct {
	table  table.Model
	width  int
	height int
}

func InitialModel() Model {
	columns := []table.Column{
		{Title: "Hash", Width: 12},
		{Title: "Author", Width: 20},
		{Title: "Date", Width: 12},
		{Title: "Message", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
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
		table: t,
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
		if tableHeight < 5 {
			tableHeight = 5
		}
		m.table.SetHeight(tableHeight)

		if m.width > 0 {
			hashWidth := 12
			authorWidth := 20
			dateWidth := 12
			messageWidth := m.width - hashWidth - authorWidth - dateWidth - 10
			if messageWidth < 20 {
				messageWidth = 20
			}

			columns := []table.Column{
				{Title: "Hash", Width: hashWidth},
				{Title: "Author", Width: authorWidth},
				{Title: "Date", Width: dateWidth},
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

	if m.width > 0 {
		return lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(tableView)
	}

	return tableView
}

func (m *Model) LoadCommits(repoPath string, since time.Time) error {
	commits, err := git.GetCommitsSince(repoPath, since)
	if err != nil {
		return fmt.Errorf("failed to load commits: %w", err)
	}

	rows := make([]table.Row, len(commits))
	for i, commit := range commits {

		shortHash := commit.Hash
		if len(shortHash) > 8 {
			shortHash = shortHash[:8]
		}

		dateStr := commit.Timestamp.Format("2006-01-02")

		message := strings.ReplaceAll(commit.Message, "\n", " ")
		message = strings.TrimSpace(message)
		if len(message) > 80 {
			message = message[:77] + "..."
		}

		rows[i] = table.Row{
			shortHash,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}

func (m Model) GetSelectedCommit() table.Row {
	return m.table.SelectedRow()
}
