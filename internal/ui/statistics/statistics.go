package statistics

import (
	"fmt"
	"project-void/internal/commands"
	"project-void/internal/jira"
	"project-void/internal/ui/common"
	commitstable "project-void/internal/ui/statistics/commits-table"
	jiratable "project-void/internal/ui/statistics/jira-table"
	slacktable "project-void/internal/ui/statistics/slack-table"
	"project-void/internal/ui/styles"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	commitsTable   commitstable.Model
	jiraTable      jiratable.Model
	slackTable     slacktable.Model
	commandHandler common.CommandHandler
	selectedFolder string
	selectedDate   time.Time
	isDev          bool
	width          int
	height         int
	loaded         bool
	loadError      string
	focusedTable   int
	command        string
	submitted      bool
	authorFilter   []string
	commitsSpinner spinner.Model
	jiraSpinner    spinner.Model
	slackSpinner   spinner.Model
	commitsLoading bool
	jiraLoading    bool
	slackLoading   bool
}

func InitialModel(selectedFolder string, selectedDate time.Time, isDev bool) Model {
	commitsTable := commitstable.InitialModel()
	jiraTable := jiratable.InitialModel()
	slackTable := slacktable.InitialModel()

	commitsTable.StartLoading()
	jiraTable.StartLoading()
	slackTable.StartLoading()

	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	commitsSpinner := spinner.New()
	commitsSpinner.Style = spinnerStyle
	commitsSpinner.Spinner = spinner.Dot

	jiraSpinner := spinner.New()
	jiraSpinner.Style = spinnerStyle
	jiraSpinner.Spinner = spinner.Dot

	slackSpinner := spinner.New()
	slackSpinner.Style = spinnerStyle
	slackSpinner.Spinner = spinner.Dot

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

	commitsTable.SetSpinner(&commitsSpinner)
	jiraTable.SetSpinner(&jiraSpinner)
	slackTable.SetSpinner(&slackSpinner)

	return Model{
		commitsTable:   commitsTable,
		jiraTable:      jiraTable,
		slackTable:     slackTable,
		commandHandler: common.NewCommandHandler("Enter a command (e.g., git a, help)..."),
		selectedFolder: selectedFolder,
		selectedDate:   selectedDate,
		isDev:          isDev,
		focusedTable:   0,
		commitsSpinner: commitsSpinner,
		jiraSpinner:    jiraSpinner,
		slackSpinner:   slackSpinner,
		commitsLoading: true,
		jiraLoading:    true,
		slackLoading:   true,
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

func loadCommitsByAuthorsCmd(folder string, since time.Time, authorNames []string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if folder == "" {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommitsByAuthors(folder, since, authorNames)
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
		m.commandHandler.Init(),
		m.commitsSpinner.Tick,
		m.jiraSpinner.Tick,
		m.slackSpinner.Tick,
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

		if key == "c" && !m.commandHandler.IsShowingCommand() && !m.commandHandler.IsShowingHelp() {
			updatedHandler, cmd, _ := m.commandHandler.Update(msg)
			m.commandHandler = updatedHandler
			return m, cmd
		}

		if m.commandHandler.IsShowingCommand() || m.commandHandler.IsShowingHelp() {
			updatedHandler, cmd, result := m.commandHandler.Update(msg)
			m.commandHandler = updatedHandler

			if result != nil {
				if result.ShouldQuit {
					return m, tea.Quit
				}

				if result.ShouldNavigate {
					m.command = result.Action
					m.submitted = true
					return m, cmd
				}

				if result.Action == "filter_by_author" && result.Data != nil {
					if commandData, ok := result.Data["command"].(commands.Command); ok {
						authorNames := commands.GetAuthorNamesFromCommand(commandData.Name)
						if len(authorNames) == 0 {
							m.commandHandler.SetError("Invalid author names in command")
							return m, cmd
						}

						if m.isDev && m.selectedFolder != "" {
							tickCmd := m.commitsTable.StartLoadingWithCmd()
							m.authorFilter = authorNames
							m.commitsLoading = true
							loadCmd := loadCommitsByAuthorsCmd(m.selectedFolder, m.selectedDate, authorNames)
							return m, tea.Batch(tickCmd, loadCmd, m.commitsSpinner.Tick)
						} else {
							m.commandHandler.SetError("Author filtering only available in development mode with a repository selected")
							return m, cmd
						}
					}
				}

				if result.Action == "clear_author_filter" {
					if m.isDev && m.selectedFolder != "" {
						m.authorFilter = nil
						m.commitsLoading = true
						tickCmd := m.commitsTable.StartLoadingWithCmd()
						loadCmd := loadCommitsCmd(m.selectedFolder, m.selectedDate)
						return m, tea.Batch(tickCmd, loadCmd, m.commitsSpinner.Tick)
					} else {
						m.commandHandler.SetError("Author filtering only available in development mode with a repository selected")
						return m, cmd
					}
				}

				if result.Action == "start" || result.Action == "reset" {
					if m.isDev && m.selectedFolder != "" && len(m.authorFilter) > 0 {
						m.authorFilter = nil
						m.commitsLoading = true
						tickCmd := m.commitsTable.StartLoadingWithCmd()
						loadCmd := loadCommitsCmd(m.selectedFolder, m.selectedDate)
						m.command = result.Action
						m.submitted = true
						return m, tea.Batch(tickCmd, loadCmd, m.commitsSpinner.Tick)
					}
					m.command = result.Action
					m.submitted = true
					return m, cmd
				}
			}

			return m, cmd
		}

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

		if key == "ctrl+c" || key == "esc" {
			return m, tea.Quit
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
		m.commitsLoading = false
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
		m.jiraLoading = false
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
		m.slackLoading = false
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

	case spinner.TickMsg:
		var cmd tea.Cmd
		if m.commitsLoading {
			m.commitsSpinner, cmd = m.commitsSpinner.Update(msg)
			m.commitsTable.SetSpinner(&m.commitsSpinner)
			cmds = append(cmds, cmd)
		}
		if m.jiraLoading {
			m.jiraSpinner, cmd = m.jiraSpinner.Update(msg)
			m.jiraTable.SetSpinner(&m.jiraSpinner)
			cmds = append(cmds, cmd)
		}
		if m.slackLoading {
			m.slackSpinner, cmd = m.slackSpinner.Update(msg)
			m.slackTable.SetSpinner(&m.slackSpinner)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

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
	if helpView := m.commandHandler.RenderHelp(m.width, m.height); helpView != "" {
		return helpView
	}

	horizontalPadding := 4
	contentWidth := m.width - (horizontalPadding * 2)

	var commandHeader string
	if commandInput := m.commandHandler.RenderCommandInput(contentWidth); commandInput != "" {
		commandHeader = commandInput
	} else {
		navHelp := "w/s: navigate tables • c: commands • esc: exit"
		commandHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(navHelp)
	}

	commandHeaderCentered := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(commandHeader)

	var mainContent string
	if m.isDev && m.selectedFolder != "" {
		commitsHeader := fmt.Sprintf("Commits for the repo %s", m.selectedFolder)

		if len(m.authorFilter) > 0 {
			authorFilterText := strings.Join(m.authorFilter, ", ")
			commitsHeader += fmt.Sprintf(" (filtered by authors: %s)", authorFilterText)
		}

		header := styles.WelcomeStyle.Width(contentWidth).Render(commitsHeader)

		var commitsText string
		if m.commitsLoading {
			commitsText = fmt.Sprintf("%s commits", m.commitsSpinner.View())
		} else {
			totalCommits := m.commitsTable.TotalCommits()
			commitsText = fmt.Sprintf("%d commits", totalCommits)
		}

		var jiraText string
		if m.jiraLoading {
			jiraText = fmt.Sprintf("%s JIRA issues", m.jiraSpinner.View())
		} else {
			totalIssues := m.jiraTable.TotalIssues()
			jiraText = fmt.Sprintf("%d JIRA issues", totalIssues)
		}

		var slackText string
		if m.slackLoading {
			slackText = fmt.Sprintf("%s Slack messages (coming soon)", m.slackSpinner.View())
		} else {
			slackText = "0 Slack messages (coming soon)"
		}

		dateInfo := fmt.Sprintf("%s, %s, %s since %s", commitsText, jiraText, slackText, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		tableView := m.commitsTable.View()
		jiraView := m.jiraTable.View()
		slackView := m.slackTable.View()

		tableViewCentered := styles.NeutralStyle.Width(contentWidth).Render(tableView)
		jiraViewCentered := styles.NeutralStyle.Width(contentWidth).Render(jiraView)
		slackViewCentered := styles.NeutralStyle.Width(contentWidth).Render(slackView)

		mainContent = header + "\n" + dateInfoRendered + "\n\n" + tableViewCentered + "\n\n" + jiraViewCentered + "\n\n" + slackViewCentered
	} else {
		var jiraText string
		if m.jiraLoading {
			jiraText = fmt.Sprintf("%s JIRA issues", m.jiraSpinner.View())
		} else {
			totalIssues := m.jiraTable.TotalIssues()
			jiraText = fmt.Sprintf("%d JIRA issues", totalIssues)
		}

		var slackText string
		if m.slackLoading {
			slackText = fmt.Sprintf("%s Slack messages (coming soon)", m.slackSpinner.View())
		} else {
			slackText = "0 Slack messages (coming soon)"
		}

		dateInfo := fmt.Sprintf("%s, %s since %s", jiraText, slackText, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		jiraView := m.jiraTable.View()
		slackView := m.slackTable.View()

		jiraViewStyled := styles.NeutralStyle.Width(contentWidth).Render(jiraView)
		slackViewStyled := styles.NeutralStyle.Width(contentWidth).Render(slackView)

		mainContent = dateInfoRendered + "\n\n" + jiraViewStyled + "\n\n" + slackViewStyled
	}

	fullContent := mainContent + "\n" + commandHeaderCentered

	centerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(1, 2)

	return centerStyle.Render(fullContent)
}

func (m Model) GetCommand() string {
	return m.command
}

func (m Model) HasCommand() bool {
	return m.submitted
}

func (m *Model) ResetCommand() {
	m.submitted = false
	m.command = ""
	m.commandHandler.ClearMessages()
}
