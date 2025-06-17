package statistics

import (
	"fmt"
	"project-void/internal/ui/styles"
	"strings"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if helpView := m.commandHandler.RenderHelp(m.width, m.height); helpView != "" {
		return helpView
	}

	if gitHelpView := m.commandHandler.RenderGitHelp(m.width, m.height); gitHelpView != "" {
		return gitHelpView
	}

	if jiraHelpView := m.commandHandler.RenderJiraHelp(m.width, m.height); jiraHelpView != "" {
		return jiraHelpView
	}

	horizontalPadding := 4
	contentWidth := m.width - (horizontalPadding * 2)

	var commandHeader string
	if commandInput := m.commandHandler.RenderCommandInput(contentWidth); commandInput != "" {
		commandHeader = commandInput
	} else {
		navHelp := "\nw/s: navigate tables • c: commands • esc: exit"
		commandHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(navHelp)
	}

	commandHeaderCentered := commandHeader

	var mainContent string

	if !m.hasGit && !m.hasJira {
		welcomeMessage := "Welcome to Project Void"
		setupMessage := "Please configure at least one service to get started:\n\n• For Git: Use 'git repo <url>' to set up a repository\n• For JIRA: Use 'jira url <url>' to set up JIRA\n\nType 'void help' for more commands"

		welcomeStyled := styles.WelcomeStyle.Render(welcomeMessage)
		setupStyled := styles.NeutralStyle.Render(setupMessage)

		mainContent = welcomeStyled + "\n\n" + setupStyled
	} else {
		var contentParts []string

		if m.hasGit && m.hasJira {
			commitsHeader := fmt.Sprintf("Commits for the repo %s", m.selectedRepoSource)
			jiraHeader := fmt.Sprintf("JIRA Issues for %s", m.selectedJiraSource)

			if len(m.authorFilter) > 0 {
				authorFilterText := strings.Join(m.authorFilter, ", ")
				commitsHeader += fmt.Sprintf(" (filtered by authors: %s)", authorFilterText)
			}

			if len(m.branchFilter) > 0 {
				branchFilterText := strings.Join(m.branchFilter, ", ")
				if len(m.authorFilter) > 0 {
					commitsHeader += fmt.Sprintf(" and branches: %s", branchFilterText)
				} else {
					commitsHeader += fmt.Sprintf(" (filtered by branches: %s)", branchFilterText)
				}
			}

			commitsHeader = styles.WelcomeStyle.Render(commitsHeader)
			jiraHeader = styles.WelcomeStyle.Render(jiraHeader)
			contentParts = append(contentParts, commitsHeader)
			contentParts = append(contentParts, jiraHeader)

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

			dateInfo := fmt.Sprintf("%s, %s since %s", commitsText, jiraText, m.selectedDate.Format("January 2, 2006"))
			contentParts = append(contentParts, styles.NeutralStyle.Render(dateInfo))

			contentParts = append(contentParts, m.commitsTable.View())

			contentParts = append(contentParts, m.jiraTable.View())

		} else if m.hasGit {
			commitsHeader := fmt.Sprintf("Commits for the repo %s", m.selectedRepoSource)

			if len(m.authorFilter) > 0 {
				authorFilterText := strings.Join(m.authorFilter, ", ")
				commitsHeader += fmt.Sprintf(" (filtered by authors: %s)", authorFilterText)
			}

			if len(m.branchFilter) > 0 {
				branchFilterText := strings.Join(m.branchFilter, ", ")
				if len(m.authorFilter) > 0 {
					commitsHeader += fmt.Sprintf(" and branches: %s", branchFilterText)
				} else {
					commitsHeader += fmt.Sprintf(" (filtered by branches: %s)", branchFilterText)
				}
			}

			header := styles.WelcomeStyle.Render(commitsHeader)
			contentParts = append(contentParts, header)

			var commitsText string
			if m.commitsLoading {
				commitsText = fmt.Sprintf("%s commits", m.commitsSpinner.View())
			} else {
				totalCommits := m.commitsTable.TotalCommits()
				commitsText = fmt.Sprintf("%d commits", totalCommits)
			}

			dateInfo := fmt.Sprintf("%s since %s", commitsText, m.selectedDate.Format("January 2, 2006"))
			contentParts = append(contentParts, styles.NeutralStyle.Render(dateInfo))

			contentParts = append(contentParts, m.commitsTable.View())

		} else if m.hasJira {
			jiraHeader := fmt.Sprintf("JIRA Issues for %s", m.selectedJiraSource)
			header := styles.WelcomeStyle.Render(jiraHeader)
			contentParts = append(contentParts, header)

			var jiraText string
			if m.jiraLoading {
				jiraText = fmt.Sprintf("%s JIRA issues", m.jiraSpinner.View())
			} else {
				totalIssues := m.jiraTable.TotalIssues()
				jiraText = fmt.Sprintf("%d JIRA issues", totalIssues)
			}

			dateInfo := fmt.Sprintf("%s since %s", jiraText, m.selectedDate.Format("January 2, 2006"))
			contentParts = append(contentParts, styles.NeutralStyle.Render(dateInfo))

			contentParts = append(contentParts, m.jiraTable.View())
		}

		mainContent = strings.Join(contentParts, "\n")
	}

	fullContent := mainContent + "\n" + commandHeaderCentered

	centerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(1, 2)

	return centerStyle.Render(fullContent)
}
