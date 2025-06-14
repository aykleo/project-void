package jiratable

import (
	"project-void/internal/ui/styles"

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
			BorderForeground(styles.HighlightColor). // blue
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
		{Title: "Issue", Width: 8},
		{Title: "Status", Width: 8},
		{Title: "Assignee", Width: 10},
	}

	rows := []table.Row{
		{"PV-101", "In Progress", "Alice"},
		{"PV-102", "Done", "Bob"},
		{"PV-103", "To Do", "Charlie"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
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
		if tableHeight < 5 {
			tableHeight = 5
		}
		m.table.SetHeight(tableHeight)

		if m.width > 0 {
			issueWidth := 8
			statusWidth := 8
			assigneeWidth := m.width - issueWidth - statusWidth - 10
			if assigneeWidth < 10 {
				assigneeWidth = 10
			}
			columns := []table.Column{
				{Title: "Issue", Width: issueWidth},
				{Title: "Status", Width: statusWidth},
				{Title: "Assignee", Width: assigneeWidth},
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
