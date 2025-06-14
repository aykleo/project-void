package statistics

import (
	"fmt"
	commitstable "project-void/internal/ui/statistics/commits-table"
	"project-void/internal/ui/styles"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	commitsTable   commitstable.Model
	selectedFolder string
	selectedDate   time.Time
	isDev          bool
	width          int
	height         int
	loaded         bool
	loadError      string
}

func InitialModel(selectedFolder string, selectedDate time.Time, isDev bool) Model {
	return Model{
		commitsTable:   commitstable.InitialModel(),
		selectedFolder: selectedFolder,
		selectedDate:   selectedDate,
		isDev:          isDev,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.commitsTable.Init(),
		func() tea.Msg {
			if m.isDev && m.selectedFolder != "" {
				var commitsTable commitstable.Model = m.commitsTable
				err := commitsTable.LoadCommits(m.selectedFolder, m.selectedDate)
				if err != nil {
					return LoadErrorMsg{Error: err.Error()}
				}
				return LoadedMsg{CommitsTable: commitsTable}
			}
			return LoadedMsg{CommitsTable: m.commitsTable}
		},
	)
}

type LoadedMsg struct {
	CommitsTable commitstable.Model
}

type LoadErrorMsg struct {
	Error string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		horizontalPadding := 4
		contentWidth := msg.Width - (horizontalPadding * 2)

		tableHeight := msg.Height - 6
		if tableHeight < 5 {
			tableHeight = 5
		}

		tableMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
		updatedTable, cmd := m.commitsTable.Update(tableMsg)
		m.commitsTable = updatedTable.(commitstable.Model)

		return m, cmd

	case LoadedMsg:
		m.loaded = true
		m.commitsTable = msg.CommitsTable
		return m, nil

	case LoadErrorMsg:
		m.loaded = true
		m.loadError = msg.Error
		return m, nil

	default:
		updatedTable, cmd := m.commitsTable.Update(msg)
		m.commitsTable = updatedTable.(commitstable.Model)
		return m, cmd
	}
}

func (m Model) View() string {
	horizontalPadding := 4
	contentWidth := m.width - (horizontalPadding * 2)

	welcomeMessage := "Project Void"
	welcome := styles.WelcomeStyle.Width(contentWidth).Render(welcomeMessage)

	quitMessage := "Q or Esc to quit"
	quit := styles.QuitStyle.Width(contentWidth).Render(quitMessage)

	var content string

	if m.isDev && m.selectedFolder != "" {

		commitsHeader := fmt.Sprintf("Commits for the repo %s", m.selectedFolder)
		header := styles.WelcomeStyle.Width(contentWidth).Render(commitsHeader)

		dateInfo := fmt.Sprintf("Since: %s", m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		if m.loadError != "" {
			errorInfo := fmt.Sprintf("Error loading commits: %s", m.loadError)
			errorRendered := styles.NeutralStyle.Width(contentWidth).Render(errorInfo)
			content = welcome + "\n" + quit + "\n\n" + header + "\n" + dateInfoRendered + "\n\n" + errorRendered
		} else if !m.loaded {
			loadingInfo := "Loading commits..."
			loadingRendered := styles.NeutralStyle.Width(contentWidth).Render(loadingInfo)
			content = welcome + "\n" + quit + "\n\n" + header + "\n" + dateInfoRendered + "\n\n" + loadingRendered
		} else {
			tableView := m.commitsTable.View()
			content = welcome + "\n" + quit + "\n\n" + header + "\n" + dateInfoRendered + "\n\n" + tableView
		}
	} else {
		generalHeader := "Statistics Dashboard"
		header := styles.WelcomeStyle.Width(contentWidth).Render(generalHeader)

		dateInfo := fmt.Sprintf("Date: %s", m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		otherStatsInfo := "Other statistics will be displayed here..."
		otherStatsRendered := styles.NeutralStyle.Width(contentWidth).Render(otherStatsInfo)

		content = welcome + "\n" + quit + "\n\n" + header + "\n" + dateInfoRendered + "\n\n" + otherStatsRendered
	}

	return styles.DocStyle.Width(m.width).Render(content)
}
