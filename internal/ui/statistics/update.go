package statistics

import (
	"project-void/internal/commands"
	commitstable "project-void/internal/ui/statistics/commits-table"
	jiratable "project-void/internal/ui/statistics/jira-table"
	slacktable "project-void/internal/ui/statistics/slack-table"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

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

		if key == "c" && !m.commandHandler.IsShowingCommand() && !m.commandHandler.IsShowingHelp() && !m.commandHandler.IsShowingGitHelp() && !m.commandHandler.IsShowingJiraHelp() {
			updatedHandler, cmd, _ := m.commandHandler.Update(msg)
			m.commandHandler = updatedHandler
			return m, cmd
		}

		if m.commandHandler.IsShowingCommand() || m.commandHandler.IsShowingHelp() || m.commandHandler.IsShowingGitHelp() || m.commandHandler.IsShowingJiraHelp() {
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

						if m.isDev && m.selectedRepoSource != "" {
							tickCmd := m.commitsTable.StartLoadingWithCmd()
							m.authorFilter = authorNames
							m.commitsLoading = true

							var loadCmd tea.Cmd
							if len(m.branchFilter) > 0 {

								loadCmd = loadCommitsByAuthorsAndBranchesCmd(m.selectedRepoSource, m.selectedDate, authorNames, m.branchFilter)
							} else {

								loadCmd = loadCommitsByAuthorsCmd(m.selectedRepoSource, m.selectedDate, authorNames)
							}
							return m, tea.Batch(tickCmd, loadCmd, m.commitsSpinner.Tick)
						} else {
							m.commandHandler.SetError("Author filtering only available in development mode with a repository selected")
							return m, cmd
						}
					}
				}

				if result.Action == "clear_author_filter" {
					if m.isDev && m.selectedRepoSource != "" {
						m.authorFilter = nil
						m.commitsLoading = true
						tickCmd := m.commitsTable.StartLoadingWithCmd()

						var loadCmd tea.Cmd
						if len(m.branchFilter) > 0 {

							loadCmd = loadCommitsByBranchesCmd(m.selectedRepoSource, m.selectedDate, m.branchFilter)
						} else {

							loadCmd = loadCommitsCmd(m.selectedRepoSource, m.selectedDate)
						}
						return m, tea.Batch(tickCmd, loadCmd, m.commitsSpinner.Tick)
					} else {
						m.commandHandler.SetError("Author filtering only available in development mode with a repository selected")
						return m, cmd
					}
				}

				if result.Action == "filter_by_branch" && result.Data != nil {
					if commandData, ok := result.Data["command"].(commands.Command); ok {
						branchNames := commands.GetBranchNamesFromCommand(commandData.Name)
						if len(branchNames) == 0 {
							m.commandHandler.SetError("Invalid branch names in command")
							return m, cmd
						}

						if m.isDev && m.selectedRepoSource != "" {
							tickCmd := m.commitsTable.StartLoadingWithCmd()
							m.branchFilter = branchNames
							m.commitsLoading = true

							var loadCmd tea.Cmd
							if len(m.authorFilter) > 0 {

								loadCmd = loadCommitsByAuthorsAndBranchesCmd(m.selectedRepoSource, m.selectedDate, m.authorFilter, branchNames)
							} else {

								loadCmd = loadCommitsByBranchesCmd(m.selectedRepoSource, m.selectedDate, branchNames)
							}
							return m, tea.Batch(tickCmd, loadCmd, m.commitsSpinner.Tick)
						} else {
							m.commandHandler.SetError("Branch filtering only available in development mode with a repository selected")
							return m, cmd
						}
					}
				}

				if result.Action == "clear_branch_filter" {
					if m.isDev && m.selectedRepoSource != "" {
						m.branchFilter = nil
						m.commitsLoading = true
						tickCmd := m.commitsTable.StartLoadingWithCmd()

						var loadCmd tea.Cmd
						if len(m.authorFilter) > 0 {

							loadCmd = loadCommitsByAuthorsCmd(m.selectedRepoSource, m.selectedDate, m.authorFilter)
						} else {

							loadCmd = loadCommitsCmd(m.selectedRepoSource, m.selectedDate)
						}
						return m, tea.Batch(tickCmd, loadCmd, m.commitsSpinner.Tick)
					} else {
						m.commandHandler.SetError("Branch filtering only available in development mode with a repository selected")
						return m, cmd
					}
				}

				if result.Action == "start" || result.Action == "reset" {
					if m.isDev && m.selectedRepoSource != "" && (len(m.authorFilter) > 0 || len(m.branchFilter) > 0) {
						m.authorFilter = nil
						m.branchFilter = nil
						m.commitsLoading = true
						tickCmd := m.commitsTable.StartLoadingWithCmd()
						loadCmd := loadCommitsCmd(m.selectedRepoSource, m.selectedDate)
						m.command = result.Action
						m.submitted = true
						return m, tea.Batch(tickCmd, loadCmd, m.commitsSpinner.Tick)
					}
					m.command = result.Action
					m.submitted = true
					return m, cmd
				}

				if result.Action == "void_set_date" {
					if dateData, ok := result.Data["date"].(time.Time); ok {
						m.selectedDate = dateData
						m.commandHandler.ClearMessages()

						var cmds []tea.Cmd

						if m.isDev && m.selectedRepoSource != "" {
							m.commitsLoading = true
							tickCmd := m.commitsTable.StartLoadingWithCmd()
							cmds = append(cmds, tickCmd, m.commitsSpinner.Tick)

							var loadCmd tea.Cmd
							if len(m.authorFilter) > 0 && len(m.branchFilter) > 0 {

								loadCmd = loadCommitsByAuthorsAndBranchesCmd(m.selectedRepoSource, m.selectedDate, m.authorFilter, m.branchFilter)
							} else if len(m.authorFilter) > 0 {

								loadCmd = loadCommitsByAuthorsCmd(m.selectedRepoSource, m.selectedDate, m.authorFilter)
							} else if len(m.branchFilter) > 0 {

								loadCmd = loadCommitsByBranchesCmd(m.selectedRepoSource, m.selectedDate, m.branchFilter)
							} else {

								loadCmd = loadCommitsCmd(m.selectedRepoSource, m.selectedDate)
							}
							cmds = append(cmds, loadCmd)
						}

						m.jiraLoading = true
						jiraTickCmd := m.jiraTable.StartLoadingWithCmd()
						jiraLoadCmd := loadJiraCmd(m.selectedDate)
						cmds = append(cmds, jiraTickCmd, jiraLoadCmd, m.jiraSpinner.Tick)

						m.slackLoading = true
						m.slackTable.StartLoading()
						slackLoadCmd := loadSlackCmd()
						cmds = append(cmds, slackLoadCmd, m.slackSpinner.Tick)

						return m, tea.Batch(cmds...)
					}
					return m, cmd
				}

				if result.Action == "jira_filter_on" || result.Action == "jira_filter_off" {
					if result.Success {

						m.jiraLoading = true
						tickCmd := m.jiraTable.StartLoadingWithCmd()
						jiraLoadCmd := loadJiraCmd(m.selectedDate)
						return m, tea.Batch(tickCmd, jiraLoadCmd, m.jiraSpinner.Tick)
					}
					return m, cmd
				}

				if result.Action == "git_set_repo" || result.Action == "git_clear_repo" {
					if result.Success {
						if repoData, ok := result.Data["repoURL"].(string); ok {
							m.selectedRepoSource = repoData

							m.isDev = m.selectedRepoSource != ""

							m.authorFilter = nil
							m.branchFilter = nil

							if m.isDev && m.selectedRepoSource != "" {
								m.commitsLoading = true
								tickCmd := m.commitsTable.StartLoadingWithCmd()
								loadCmd := loadCommitsCmd(m.selectedRepoSource, m.selectedDate)
								return m, tea.Batch(tickCmd, loadCmd, m.commitsSpinner.Tick)
							} else {
								m.commitsTable = commitstable.InitialModel()
								return m, cmd
							}
						}
					}
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
