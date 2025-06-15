package welcome

import (
	"fmt"
	"project-void/internal/config"
	"project-void/internal/ui/common"
	"project-void/internal/ui/styles"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	commandHandler common.CommandHandler
	width          int
	height         int
	command        string
	submitted      bool
	selectedDate   *time.Time
}

func InitialModel() Model {
	return Model{
		commandHandler: common.NewCommandHandler("Enter a command (e.g., git repo <url>, void help)..."),
		submitted:      false,
		selectedDate:   nil,
	}
}

func (m Model) Init() tea.Cmd {
	return m.commandHandler.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

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

		if result.Success && result.Action == "void_set_date" {
			if dateData, ok := result.Data["date"].(time.Time); ok {
				m.selectedDate = &dateData
			}
			return m, cmd
		}

		if result.Success {
			return m, cmd
		}
	}

	return m, cmd
}

func (m Model) View() string {
	if helpView := m.commandHandler.RenderHelp(m.width, m.height); helpView != "" {
		return helpView
	}

	if gitHelpView := m.commandHandler.RenderGitHelp(m.width, m.height); gitHelpView != "" {
		return gitHelpView
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

func (m Model) GetSelectedDate() *time.Time {
	return m.selectedDate
}
