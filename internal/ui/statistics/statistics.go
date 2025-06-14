package statistics

import (
	"fmt"
	"project-void/internal/commands"
	"project-void/internal/jira"
	commitstable "project-void/internal/ui/statistics/commits-table"
	jiratable "project-void/internal/ui/statistics/jira-table"
	slacktable "project-void/internal/ui/statistics/slack-table"
	"project-void/internal/ui/styles"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	commitsTable   commitstable.Model
	jiraTable      jiratable.Model
	slackTable     slacktable.Model
	textInput      textinput.Model
	selectedFolder string
	selectedDate   time.Time
	isDev          bool
	width          int
	height         int
	loaded         bool
	loadError      string
	focusedTable   int
	commandError   string
	showingHelp    bool
	showingCommand bool
	command        string
	submitted      bool
	authorFilter   []string
}

func InitialModel(selectedFolder string, selectedDate time.Time, isDev bool) Model {
	commitsTable := commitstable.InitialModel()
	jiraTable := jiratable.InitialModel()
	slackTable := slacktable.InitialModel()

	commitsTable.StartLoading()
	jiraTable.StartLoading()
	slackTable.StartLoading()

	ti := textinput.New()
	ti.Placeholder = "Enter a command (e.g., help)..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

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
		textInput:      ti,
		selectedFolder: selectedFolder,
		selectedDate:   selectedDate,
		isDev:          isDev,
		focusedTable:   0,
		showingCommand: false,
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
		textinput.Blink,
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

		if msg.Width > 60 {
			m.textInput.Width = 50
		} else {
			m.textInput.Width = msg.Width - 10
		}

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
		if m.showingHelp {
			m.showingHelp = false
			m.commandError = ""
			return m, nil
		}

		if m.showingCommand {
			switch msg.Type {
			case tea.KeyEnter:
				inputValue := m.textInput.Value()
				validatedCmd, err := commands.ValidateCommand(inputValue)
				if err != nil {
					m.commandError = err.Error()
					m.textInput.SetValue("")
					return m, nil
				}

				if validatedCmd.Action == "help" {
					m.showingHelp = true
					m.commandError = ""
					m.textInput.SetValue("")
					return m, nil
				}

				if validatedCmd.Action == "quit" {
					return m, tea.Quit
				}

				if validatedCmd.Action == "start" || validatedCmd.Action == "reset" {
					if m.isDev && m.selectedFolder != "" && len(m.authorFilter) > 0 {
						m.authorFilter = nil
						tickCmd := m.commitsTable.StartLoadingWithCmd()
						loadCmd := loadCommitsCmd(m.selectedFolder, m.selectedDate)
						m.commandError = ""
						m.textInput.SetValue("")
						m.showingCommand = false
						return m, tea.Batch(tickCmd, loadCmd)
					}
					m.command = validatedCmd.Action
					m.submitted = true
					m.commandError = ""
					return m, nil
				}

				if validatedCmd.Action == "filter_by_author" {
					authorNames := commands.GetAuthorNamesFromCommand(validatedCmd.Name)
					if len(authorNames) == 0 {
						m.commandError = "Invalid author names in command"
						m.textInput.SetValue("")
						return m, nil
					}

					if m.isDev && m.selectedFolder != "" {
						tickCmd := m.commitsTable.StartLoadingWithCmd()
						m.authorFilter = authorNames
						loadCmd := loadCommitsByAuthorsCmd(m.selectedFolder, m.selectedDate, authorNames)
						m.commandError = ""
						m.textInput.SetValue("")
						m.showingCommand = false
						return m, tea.Batch(tickCmd, loadCmd)
					} else {
						m.commandError = "Author filtering only available in development mode with a repository selected"
						m.textInput.SetValue("")
						return m, nil
					}
				}

				if validatedCmd.Action == "clear_author_filter" {
					if m.isDev && m.selectedFolder != "" {
						m.authorFilter = nil
						tickCmd := m.commitsTable.StartLoadingWithCmd()
						loadCmd := loadCommitsCmd(m.selectedFolder, m.selectedDate)
						m.commandError = ""
						m.textInput.SetValue("")
						m.showingCommand = false
						return m, tea.Batch(tickCmd, loadCmd)
					} else {
						m.commandError = "Author filtering only available in development mode with a repository selected"
						m.textInput.SetValue("")
						return m, nil
					}
				}

				m.commandError = fmt.Sprintf("Unknown command action: %s", validatedCmd.Action)
				m.textInput.SetValue("")
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				if msg.Type == tea.KeyEsc {
					return m, nil
				}
				return m, tea.Quit
			}

			if msg.String() == "'" {
				m.showingCommand = false
				m.commandError = ""
				m.textInput.SetValue("")
				return m, nil
			}

			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		key := msg.String()

		if key == "c" {
			m.showingCommand = true
			m.textInput.Focus()
			return m, nil
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
	if m.showingHelp {
		helpText := commands.GetHelpText()
		helpContent := fmt.Sprintf("%s\n\nPress any key to return to statistics", helpText)

		centerStyle := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1, 2)

		return centerStyle.Render(styles.NeutralStyle.Render(helpContent))
	}

	horizontalPadding := 4
	contentWidth := m.width - (horizontalPadding * 2)

	var commandHeader string
	if m.showingCommand {
		if m.commandError != "" {
			errorText := fmt.Sprintf("Error: %s", m.commandError)
			commandHeader = fmt.Sprintf("%s\nCommand: %s\n%s",
				lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errorText),
				m.textInput.View(),
				styles.QuitStyle.Render("Press ' to exit command mode, Esc to go back to start"))
		} else {
			commandHeader = fmt.Sprintf("Command: %s\n%s",
				m.textInput.View(),
				styles.QuitStyle.Render("Press ' to exit command mode, Esc to go back to start"))
		}
	} else {
		navHelp := "Use w/s to navigate tables, c for commands, q/Esc to go back"
		commandHeader = styles.QuitStyle.Render(navHelp)
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

		mainContent = header + "\n" + dateInfoRendered + "\n\n" + tableViewCentered + "\n\n" + jiraViewCentered + "\n\n" + slackViewCentered
	} else {

		totalIssues := m.jiraTable.TotalIssues()
		dateInfo := fmt.Sprintf("%d JIRA issues, %d Slack messages (coming soon) since %s", totalIssues, 0, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Width(contentWidth).Render(dateInfo)

		jiraView := m.jiraTable.View()
		slackView := m.slackTable.View()

		jiraViewStyled := styles.NeutralStyle.Width(contentWidth).Render(jiraView)
		slackViewStyled := styles.NeutralStyle.Width(contentWidth).Render(slackView)

		mainContent = dateInfoRendered + "\n\n" + jiraViewStyled + "\n\n" + slackViewStyled
	}

	fullContent := commandHeaderCentered + "\n" + mainContent

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
	m.commandError = ""
	m.textInput.SetValue("")
}
