package statistics

import (
	"fmt"
	commitstable "project-void/internal/ui/statistics/commits-table"
	jiratable "project-void/internal/ui/statistics/jira-table"
	slacktable "project-void/internal/ui/statistics/slack-table"
	"project-void/internal/ui/styles"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	commitsTable   commitstable.Model
	jiraTable      jiratable.Model
	slackTable     slacktable.Model
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
		jiraTable:      jiratable.InitialModel(),
		slackTable:     slacktable.InitialModel(),
		selectedFolder: selectedFolder,
		selectedDate:   selectedDate,
		isDev:          isDev,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.commitsTable.Init(),
		m.jiraTable.Init(),
		m.slackTable.Init(),
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

		// Divide the available height among the three tables
		availableHeight := msg.Height - 12 // leave room for headers, padding, etc.
		tableCount := 3
		tableHeight := availableHeight / tableCount
		if tableHeight < 5 {
			tableHeight = 5
		}

		tableMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
		updatedCommits, cmd1 := m.commitsTable.Update(tableMsg)
		updatedJira, cmd2 := m.jiraTable.Update(tableMsg)
		updatedSlack, cmd3 := m.slackTable.Update(tableMsg)
		m.commitsTable = updatedCommits.(commitstable.Model)
		m.jiraTable = updatedJira.(jiratable.Model)
		m.slackTable = updatedSlack.(slacktable.Model)

		return m, tea.Batch(cmd1, cmd2, cmd3)

	case LoadedMsg:
		m.loaded = true
		m.commitsTable = msg.CommitsTable
		return m, nil

	case LoadErrorMsg:
		m.loaded = true
		m.loadError = msg.Error
		return m, nil

	default:
		updatedCommits, cmd1 := m.commitsTable.Update(msg)
		updatedJira, cmd2 := m.jiraTable.Update(msg)
		updatedSlack, cmd3 := m.slackTable.Update(msg)
		m.commitsTable = updatedCommits.(commitstable.Model)
		m.jiraTable = updatedJira.(jiratable.Model)
		m.slackTable = updatedSlack.(slacktable.Model)
		return m, tea.Batch(cmd1, cmd2, cmd3)
	}
}

func (m Model) View() string {
	horizontalPadding := 4
	contentWidth := m.width - (horizontalPadding * 2)

	var content string

	if m.isDev && m.selectedFolder != "" {

		commitsHeader := fmt.Sprintf("Commits for the repo %s", m.selectedFolder)
		header := styles.WelcomeStyle.Width(contentWidth).Render(commitsHeader)

		totalCommits := m.commitsTable.TotalCommits()

		dateInfo := fmt.Sprintf("%d commits since %s", totalCommits, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		if m.loadError != "" {
			errorInfo := fmt.Sprintf("Error loading commits: %s", m.loadError)
			errorRendered := styles.NeutralStyle.Width(contentWidth).Render(errorInfo)
			content = header + "\n" + dateInfoRendered + "\n\n" + errorRendered
		} else if !m.loaded {
			loadingInfo := "Loading commits..."
			loadingRendered := styles.NeutralStyle.Width(contentWidth).Render(loadingInfo)
			content = header + "\n" + dateInfoRendered + "\n\n" + loadingRendered
		} else {
			tableView := m.commitsTable.View()
			jiraView := m.jiraTable.View()
			slackView := m.slackTable.View()
			tableViewCentered := styles.NeutralStyle.Width(contentWidth).Render(tableView)
			jiraViewCentered := styles.NeutralStyle.Width(contentWidth).Render(jiraView)
			slackViewCentered := styles.NeutralStyle.Width(contentWidth).Render(slackView)
			content = header + "\n" + dateInfoRendered + "\n\n" + tableViewCentered + "\n\n" + jiraViewCentered + "\n\n" + slackViewCentered
		}
	} else {
		generalHeader := "Statistics Dashboard"
		header := styles.WelcomeStyle.Width(contentWidth).Render(generalHeader)

		dateInfo := fmt.Sprintf("Date: %s", m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		jiraView := m.jiraTable.View()
		slackView := m.slackTable.View()
		jiraViewCentered := styles.NeutralStyle.Width(contentWidth).Render(jiraView)
		slackViewCentered := styles.NeutralStyle.Width(contentWidth).Render(slackView)

		content = header + "\n" + dateInfoRendered + "\n\n" + jiraViewCentered + "\n\n" + slackViewCentered
	}

	return styles.DocStyle.Width(m.width).Render(content)
}
