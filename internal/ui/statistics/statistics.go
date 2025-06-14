package statistics

import (
	"fmt"
	"project-void/internal/jira"
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
	focusedTable   int
}

func InitialModel(selectedFolder string, selectedDate time.Time, isDev bool) Model {
	commitsTable := commitstable.InitialModel()
	jiraTable := jiratable.InitialModel()
	slackTable := slacktable.InitialModel()

	if isDev {
		commitsTable.Focus()
		commitsTable.SetFocusedStyle()
		jiraTable.Blur()
		jiraTable.SetBlurredStyle()
		slackTable.Blur()
		slackTable.SetBlurredStyle()
	} else {
		jiraTable.Focus()
		jiraTable.SetFocusedStyle()
		slackTable.Blur()
		slackTable.SetBlurredStyle()
	}

	return Model{
		commitsTable:   commitsTable,
		jiraTable:      jiraTable,
		slackTable:     slackTable,
		selectedFolder: selectedFolder,
		selectedDate:   selectedDate,
		isDev:          isDev,
		focusedTable:   0,
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
		func() tea.Msg {
			config, err := jira.LoadConfig()
			if err != nil {
				return JiraLoadErrorMsg{Error: fmt.Sprintf("Failed to load JIRA config: %v", err)}
			}

			client := jira.NewClientFromConfig(config)

			if err := client.TestConnection(); err != nil {
				return JiraLoadErrorMsg{Error: fmt.Sprintf("JIRA connection failed: %v", err)}
			}

			var jiraTable jiratable.Model = m.jiraTable
			err = jiraTable.LoadIssues(client, m.selectedDate, config)
			if err != nil {
				return JiraLoadErrorMsg{Error: fmt.Sprintf("Failed to load JIRA issues: %v", err)}
			}
			return JiraLoadedMsg{JiraTable: jiraTable}
		},
	)
}

type LoadedMsg struct {
	CommitsTable commitstable.Model
}

type LoadErrorMsg struct {
	Error string
}

type JiraLoadedMsg struct {
	JiraTable jiratable.Model
}

type JiraLoadErrorMsg struct {
	Error string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		horizontalPadding := 4
		contentWidth := msg.Width - (horizontalPadding * 2)

		availableHeight := msg.Height - 12

		if m.isDev {
			tableHeight := availableHeight / 3
			if tableHeight < 5 {
				tableHeight = 5
			}

			commitsMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			updatedCommits, _ := m.commitsTable.Update(commitsMsg)
			m.commitsTable = updatedCommits.(commitstable.Model)

			jiraMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			slackMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			updatedJira, _ := m.jiraTable.Update(jiraMsg)
			updatedSlack, _ := m.slackTable.Update(slackMsg)
			m.jiraTable = updatedJira.(jiratable.Model)
			m.slackTable = updatedSlack.(slacktable.Model)
		} else {
			tableHeight := availableHeight / 2
			if tableHeight < 5 {
				tableHeight = 5
			}

			jiraMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			slackMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			updatedJira, _ := m.jiraTable.Update(jiraMsg)
			updatedSlack, _ := m.slackTable.Update(slackMsg)
			m.jiraTable = updatedJira.(jiratable.Model)
			m.slackTable = updatedSlack.(slacktable.Model)
		}

		if m.isDev {
			if m.focusedTable == 0 {
				m.commitsTable.Focus()
				m.commitsTable.SetFocusedStyle()
				m.jiraTable.Blur()
				m.jiraTable.SetBlurredStyle()
				m.slackTable.Blur()
				m.slackTable.SetBlurredStyle()
			} else if m.focusedTable == 1 {
				m.commitsTable.Blur()
				m.commitsTable.SetBlurredStyle()
				m.jiraTable.Focus()
				m.jiraTable.SetFocusedStyle()
				m.slackTable.Blur()
				m.slackTable.SetBlurredStyle()
			} else {
				m.commitsTable.Blur()
				m.commitsTable.SetBlurredStyle()
				m.jiraTable.Blur()
				m.jiraTable.SetBlurredStyle()
				m.slackTable.Focus()
				m.slackTable.SetFocusedStyle()
			}
		} else {
			if m.focusedTable == 0 {
				m.jiraTable.Focus()
				m.jiraTable.SetFocusedStyle()
				m.slackTable.Blur()
				m.slackTable.SetBlurredStyle()
			} else {
				m.jiraTable.Blur()
				m.jiraTable.SetBlurredStyle()
				m.slackTable.Focus()
				m.slackTable.SetFocusedStyle()
			}
		}
		return m, tea.Batch()

	case tea.KeyMsg:
		key := msg.String()

		if key == "w" || key == "s" {
			if m.isDev {
				if key == "w" {
					m.focusedTable = (m.focusedTable + 2) % 3
				} else {
					m.focusedTable = (m.focusedTable + 1) % 3
				}
			} else {
				if key == "w" {
					m.focusedTable = (m.focusedTable + 1) % 2
				} else {
					m.focusedTable = (m.focusedTable + 1) % 2
				}
			}

			if m.isDev {
				if m.focusedTable == 0 {
					m.commitsTable.Focus()
					m.commitsTable.SetFocusedStyle()
					m.jiraTable.Blur()
					m.jiraTable.SetBlurredStyle()
					m.slackTable.Blur()
					m.slackTable.SetBlurredStyle()
				} else if m.focusedTable == 1 {
					m.commitsTable.Blur()
					m.commitsTable.SetBlurredStyle()
					m.jiraTable.Focus()
					m.jiraTable.SetFocusedStyle()
					m.slackTable.Blur()
					m.slackTable.SetBlurredStyle()
				} else {
					m.commitsTable.Blur()
					m.commitsTable.SetBlurredStyle()
					m.jiraTable.Blur()
					m.jiraTable.SetBlurredStyle()
					m.slackTable.Focus()
					m.slackTable.SetFocusedStyle()
				}
			} else {
				if m.focusedTable == 0 {
					m.jiraTable.Focus()
					m.jiraTable.SetFocusedStyle()
					m.slackTable.Blur()
					m.slackTable.SetBlurredStyle()
				} else {
					m.jiraTable.Blur()
					m.jiraTable.SetBlurredStyle()
					m.slackTable.Focus()
					m.slackTable.SetFocusedStyle()
				}
			}
			return m, nil
		}

		rowKeys := map[string]bool{"up": true, "down": true, "k": true, "j": true, "pgup": true, "pgdown": true, "home": true, "end": true}
		if rowKeys[key] {
			if m.isDev {
				if m.focusedTable == 0 {
					updated, cmd := m.commitsTable.Update(msg)
					m.commitsTable = updated.(commitstable.Model)
					return m, cmd
				} else if m.focusedTable == 1 {
					updated, cmd := m.jiraTable.Update(msg)
					m.jiraTable = updated.(jiratable.Model)
					return m, cmd
				} else {
					updated, cmd := m.slackTable.Update(msg)
					m.slackTable = updated.(slacktable.Model)
					return m, cmd
				}
			} else {
				if m.focusedTable == 0 {
					updated, cmd := m.jiraTable.Update(msg)
					m.jiraTable = updated.(jiratable.Model)
					return m, cmd
				} else {
					updated, cmd := m.slackTable.Update(msg)
					m.slackTable = updated.(slacktable.Model)
					return m, cmd
				}
			}
		}

		updatedCommits, cmd1 := m.commitsTable.Update(msg)
		updatedJira, cmd2 := m.jiraTable.Update(msg)
		updatedSlack, cmd3 := m.slackTable.Update(msg)
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

	case JiraLoadedMsg:
		m.jiraTable = msg.JiraTable
		return m, nil

	case JiraLoadErrorMsg:
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
		totalIssues := m.jiraTable.TotalIssues()

		dateInfo := fmt.Sprintf("%d commits, %d JIRA issues since %s", totalCommits, totalIssues, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		if m.loadError != "" {
			errorInfo := fmt.Sprintf("Error loading commits: %s", m.loadError)
			errorRendered := styles.NeutralStyle.Width(contentWidth).Render(errorInfo)
			content = header + "\n" + dateInfoRendered + "\n\n" + errorRendered
		} else if !m.loaded {
			loadingInfo := "Loading commits and JIRA issues..."
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

		totalIssues := m.jiraTable.TotalIssues()
		dateInfo := fmt.Sprintf("%d JIRA issues since %s", totalIssues, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		jiraView := m.jiraTable.View()
		slackView := m.slackTable.View()

		jiraViewStyled := styles.NeutralStyle.Width(contentWidth).Render(jiraView)
		slackViewStyled := styles.NeutralStyle.Width(contentWidth).Render(slackView)

		content = header + "\n" + dateInfoRendered + "\n\n" + jiraViewStyled + "\n\n" + slackViewStyled
	}

	return styles.DocStyle.Width(m.width).Render(content)
}
