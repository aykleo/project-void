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
	if m.isDev && m.selectedRepoSource != "" {
		commitsHeader := fmt.Sprintf("Commits for the repo %s", m.selectedRepoSource)

		if len(m.authorFilter) > 0 {
			authorFilterText := strings.Join(m.authorFilter, ", ")
			commitsHeader += fmt.Sprintf(" (filtered by authors: %s)", authorFilterText)
		}

		header := styles.WelcomeStyle.Render(commitsHeader)

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

		// var slackText string
		// if m.slackLoading {
		// 	slackText = fmt.Sprintf("%s Slack messages (coming soon)", m.slackSpinner.View())
		// } else {
		// 	slackText = "0 Slack messages (coming soon)"
		// }

		dateInfo := fmt.Sprintf("%s, %s since %s", commitsText, jiraText, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Render(dateInfo)

		tableView := m.commitsTable.View()
		jiraView := m.jiraTable.View()

		tableViewCentered := tableView
		jiraViewCentered := jiraView

		mainContent = header + "\n" + dateInfoRendered + "\n\n" + tableViewCentered + "\n\n" + jiraViewCentered
	} else {
		var jiraText string
		if m.jiraLoading {
			jiraText = fmt.Sprintf("%s JIRA issues", m.jiraSpinner.View())
		} else {
			totalIssues := m.jiraTable.TotalIssues()
			jiraText = fmt.Sprintf("%d JIRA issues", totalIssues)
		}

		// var slackText string
		// if m.slackLoading {
		// 	slackText = fmt.Sprintf("%s Slack messages (coming soon)", m.slackSpinner.View())
		// } else {
		// 	slackText = "0 Slack messages (coming soon)"
		// }

		dateInfo := fmt.Sprintf("%s since %s", jiraText, m.selectedDate.Format("January 2, 2006"))
		dateInfoRendered := styles.NeutralStyle.Render(dateInfo)

		jiraView := m.jiraTable.View()

		jiraViewStyled := jiraView

		mainContent = dateInfoRendered + "\n\n" + jiraViewStyled
	}

	fullContent := mainContent + "\n" + commandHeaderCentered

	centerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(1, 2)

	return centerStyle.Render(fullContent)
}
