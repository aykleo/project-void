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

	commitsTable.StartLoading()
	jiraTable.StartLoading()
	slackTable.StartLoading()

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

func loadCommitsCmd(folder string, since time.Time) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if folder == "" {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommits(folder, since)
		if err != nil {
			return LoadErrorMsg{Error: err.Error()}
		}

		return LoadedMsg{CommitsTable: commitsTable}
	})
}

func loadJiraCmd(since time.Time) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		var jiraTable jiratable.Model = jiratable.InitialModel()
		jiraTable.StartLoading()

		config, err := jira.LoadConfig()
		if err != nil {
			return JiraLoadErrorMsg{Error: fmt.Sprintf("Failed to load JIRA config: %v", err)}
		}

		client := jira.NewClientFromConfig(config)

		if err := client.TestConnection(); err != nil {
			return JiraLoadErrorMsg{Error: fmt.Sprintf("JIRA connection failed: %v", err)}
		}

		err = jiraTable.LoadIssues(client, since, config)
		if err != nil {
			return JiraLoadErrorMsg{Error: fmt.Sprintf("Failed to load JIRA issues: %v", err)}
		}
		return JiraLoadedMsg{JiraTable: jiraTable}
	})
}

func loadSlackCmd() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		var slackTable slacktable.Model = slacktable.InitialModel()
		slackTable.StartLoading()
		slackTable.SetPlaceholder()
		return SlackLoadedMsg{SlackTable: slackTable}
	})
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.commitsTable.Init(),
		m.jiraTable.Init(),
		m.slackTable.Init(),
		loadCommitsCmd(m.selectedFolder, m.selectedDate),
		loadJiraCmd(m.selectedDate),
		loadSlackCmd(),
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

type SlackLoadedMsg struct {
	SlackTable slacktable.Model
}

type SlackLoadErrorMsg struct {
	Error string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		horizontalPadding := 4
		contentWidth := msg.Width - (horizontalPadding * 2)

		availableHeight := msg.Height - 12

		if m.isDev {
			tableHeight := availableHeight / 3
			if tableHeight < 3 {
				tableHeight = 3
			}

			commitsMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			updatedCommits, cmd1 := m.commitsTable.Update(commitsMsg)
			m.commitsTable = updatedCommits.(commitstable.Model)
			if cmd1 != nil {
				cmds = append(cmds, cmd1)
			}

			jiraMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			slackMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			updatedJira, cmd2 := m.jiraTable.Update(jiraMsg)
			updatedSlack, cmd3 := m.slackTable.Update(slackMsg)
			m.jiraTable = updatedJira.(jiratable.Model)
			m.slackTable = updatedSlack.(slacktable.Model)
			if cmd2 != nil {
				cmds = append(cmds, cmd2)
			}
			if cmd3 != nil {
				cmds = append(cmds, cmd3)
			}
		} else {
			tableHeight := availableHeight / 2
			if tableHeight < 3 {
				tableHeight = 3
			}

			jiraMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			slackMsg := tea.WindowSizeMsg{Width: contentWidth, Height: tableHeight}
			updatedJira, cmd2 := m.jiraTable.Update(jiraMsg)
			updatedSlack, cmd3 := m.slackTable.Update(slackMsg)
			m.jiraTable = updatedJira.(jiratable.Model)
			m.slackTable = updatedSlack.(slacktable.Model)
			if cmd2 != nil {
				cmds = append(cmds, cmd2)
			}
			if cmd3 != nil {
				cmds = append(cmds, cmd3)
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
		return m, tea.Batch(cmds...)

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

		if cmd1 != nil {
			cmds = append(cmds, cmd1)
		}
		if cmd2 != nil {
			cmds = append(cmds, cmd2)
		}
		if cmd3 != nil {
			cmds = append(cmds, cmd3)
		}
		return m, tea.Batch(cmds...)

	case LoadedMsg:
		m.loaded = true
		m.commitsTable = msg.CommitsTable
		updatedCommits, cmd := m.commitsTable.Update(commitstable.LoadingCompleteMsg{})
		m.commitsTable = updatedCommits.(commitstable.Model)

		if m.isDev && m.focusedTable == 0 {
			m.commitsTable.Focus()
			m.commitsTable.SetFocusedStyle()
		} else {
			m.commitsTable.Blur()
			m.commitsTable.SetBlurredStyle()
		}

		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case LoadErrorMsg:
		m.loaded = true
		m.loadError = msg.Error
		return m, nil

	case JiraLoadedMsg:
		m.jiraTable = msg.JiraTable
		updatedJira, cmd := m.jiraTable.Update(jiratable.LoadingCompleteMsg{})
		m.jiraTable = updatedJira.(jiratable.Model)

		if (m.isDev && m.focusedTable == 1) || (!m.isDev && m.focusedTable == 0) {
			m.jiraTable.Focus()
			m.jiraTable.SetFocusedStyle()
		} else {
			m.jiraTable.Blur()
			m.jiraTable.SetBlurredStyle()
		}

		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case JiraLoadErrorMsg:
		return m, nil

	case SlackLoadedMsg:
		m.slackTable = msg.SlackTable
		updatedSlack, cmd := m.slackTable.Update(slacktable.LoadingCompleteMsg{})
		m.slackTable = updatedSlack.(slacktable.Model)

		if (m.isDev && m.focusedTable == 2) || (!m.isDev && m.focusedTable == 1) {
			m.slackTable.Focus()
			m.slackTable.SetFocusedStyle()
		} else {
			m.slackTable.Blur()
			m.slackTable.SetBlurredStyle()
		}

		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case SlackLoadErrorMsg:
		return m, nil

	default:
		updatedCommits, cmd1 := m.commitsTable.Update(msg)
		updatedJira, cmd2 := m.jiraTable.Update(msg)
		updatedSlack, cmd3 := m.slackTable.Update(msg)
		m.commitsTable = updatedCommits.(commitstable.Model)
		m.jiraTable = updatedJira.(jiratable.Model)
		m.slackTable = updatedSlack.(slacktable.Model)

		if cmd1 != nil {
			cmds = append(cmds, cmd1)
		}
		if cmd2 != nil {
			cmds = append(cmds, cmd2)
		}
		if cmd3 != nil {
			cmds = append(cmds, cmd3)
		}
		return m, tea.Batch(cmds...)
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

		dateInfo := fmt.Sprintf("%d commits, %d JIRA issues, %d Slack messages (coming soon) since %s", totalCommits, totalIssues, 0, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		tableView := m.commitsTable.View()
		jiraView := m.jiraTable.View()
		slackView := m.slackTable.View()

		tableViewCentered := styles.NeutralStyle.Width(contentWidth).Render(tableView)
		jiraViewCentered := styles.NeutralStyle.Width(contentWidth).Render(jiraView)
		slackViewCentered := styles.NeutralStyle.Width(contentWidth).Render(slackView)

		content = header + "\n" + dateInfoRendered + "\n\n" + tableViewCentered + "\n\n" + jiraViewCentered + "\n\n" + slackViewCentered
	} else {
		generalHeader := "Your Statistics"
		header := styles.WelcomeStyle.Width(contentWidth).Render(generalHeader)

		totalIssues := m.jiraTable.TotalIssues()
		dateInfo := fmt.Sprintf("%d JIRA issues, %d Slack messages (coming soon) since %s", totalIssues, 0, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		jiraView := m.jiraTable.View()
		slackView := m.slackTable.View()

		jiraViewStyled := styles.NeutralStyle.Width(contentWidth).Render(jiraView)
		slackViewStyled := styles.NeutralStyle.Width(contentWidth).Render(slackView)

		content = header + "\n" + dateInfoRendered + "\n\n" + jiraViewStyled + "\n\n" + slackViewStyled
	}

	return styles.DocStyle.Width(m.width).Render(content)
}
