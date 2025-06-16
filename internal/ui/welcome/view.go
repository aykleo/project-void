package welcome

import (
	"fmt"
	"project-void/internal/config"
	"project-void/internal/ui/styles"
	"time"

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

	welcomeText := "Welcome to Project Void"
	welcomeStyled := styles.WelcomeStyle.Render(welcomeText)

	var gitDescription string
	var jiraDescription string
	currentDate := time.Now().Format("2006-01-02")
	if userConfig, err := config.LoadUserConfig(); err == nil && userConfig.Git.RepoURL != "" {
		gitDescription = fmt.Sprintf("Current repo: %s", lipgloss.NewStyle().Foreground(styles.HighlightColor).Render(userConfig.Git.RepoURL))
	} else {
		gitDescription = "If you are a developer, you can use git status to check your repository configuration."
	}

	if userConfig, err := config.LoadUserConfig(); err == nil && userConfig.Jira.Username != "" {
		jiraDescription = fmt.Sprintf("Current JIRA username: %s", lipgloss.NewStyle().Foreground(styles.HighlightColor).Render(userConfig.Jira.Username))
	} else {
		jiraDescription = "If you are a JIRA user, you can use jira status to check your JIRA configuration."
	}
	gitDescriptionStyled := styles.NeutralStyle.Render(gitDescription)
	jiraDescriptionStyled := styles.NeutralStyle.Render(jiraDescription)
	dateDescriptionStyled := styles.NeutralStyle.Render(fmt.Sprintf("Today is: %s", lipgloss.NewStyle().Foreground(styles.HighlightColor).Render(currentDate)))

	var inputSection string
	if m.submitted {
		inputSection = styles.NeutralStyle.Render(fmt.Sprintf("Command processed: %s\n\nNavigating...", m.command))
	} else {
		helpText := "Use 'git repo <url>' to configure a repository, 'void sd <date>' to set analysis date, or 'void st' to begin"
		inputSection = m.commandHandler.RenderCommandPrompt(helpText)
	}

	content := fmt.Sprintf("%s\n\n%s\n%s\n%s\n\n%s", welcomeStyled, gitDescriptionStyled, jiraDescriptionStyled, dateDescriptionStyled, inputSection)

	centerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(1, 2)

	return centerStyle.Render(content)
}
