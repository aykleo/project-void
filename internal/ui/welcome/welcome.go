package welcome

import (
	"fmt"
	"project-void/internal/commands"
	"project-void/internal/config"
	"project-void/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textInput      textinput.Model
	width          int
	height         int
	command        string
	commandError   string
	submitted      bool
	showingHelp    bool
	successMessage string
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter a command (e.g., help)..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return Model{
		textInput:   ti,
		submitted:   false,
		showingHelp: false,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if msg.Width > 60 {
			m.textInput.Width = 50
		} else {
			m.textInput.Width = msg.Width - 10
		}
		return m, nil

	case tea.KeyMsg:
		if m.showingHelp {
			m.showingHelp = false
			m.commandError = ""
			return m, nil
		}

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
				m.successMessage = ""
				m.textInput.SetValue("")
				return m, nil
			}

			if validatedCmd.Action == "jira_status" {
				status, err := config.GetJiraConfigStatus()
				if err != nil {
					m.commandError = fmt.Sprintf("Failed to get JIRA status: %v", err)
				} else {
					m.successMessage = status
					m.commandError = ""
				}
				m.textInput.SetValue("")
				return m, nil
			}

			if validatedCmd.Action == "jira_set_url" {
				key, value := commands.GetJiraConfigValue(validatedCmd.Name)
				if key == "" || value == "" {
					m.commandError = "Invalid JIRA URL command"
					m.textInput.SetValue("")
					return m, nil
				}

				err := config.SetJiraConfig("url", value)
				if err != nil {
					m.commandError = fmt.Sprintf("Failed to set JIRA URL: %v", err)
				} else {
					m.successMessage = fmt.Sprintf("✓ JIRA URL set to: %s", value)
					m.commandError = ""
				}
				m.textInput.SetValue("")
				return m, nil
			}

			if validatedCmd.Action == "jira_set_user" {
				key, value := commands.GetJiraConfigValue(validatedCmd.Name)
				if key == "" || value == "" {
					m.commandError = "Invalid JIRA user command"
					m.textInput.SetValue("")
					return m, nil
				}

				err := config.SetJiraConfig("user", value)
				if err != nil {
					m.commandError = fmt.Sprintf("Failed to set JIRA user: %v", err)
				} else {
					m.successMessage = fmt.Sprintf("✓ JIRA username set to: %s", value)
					m.commandError = ""
				}
				m.textInput.SetValue("")
				return m, nil
			}

			if validatedCmd.Action == "jira_set_token" {
				key, value := commands.GetJiraConfigValue(validatedCmd.Name)
				if key == "" || value == "" {
					m.commandError = "Invalid JIRA token command"
					m.textInput.SetValue("")
					return m, nil
				}

				err := config.SetJiraConfig("token", value)
				if err != nil {
					m.commandError = fmt.Sprintf("Failed to set JIRA token: %v", err)
				} else {
					maskedToken := value
					if len(maskedToken) > 8 {
						maskedToken = maskedToken[:4] + "..." + maskedToken[len(maskedToken)-4:]
					}
					m.successMessage = fmt.Sprintf("✓ JIRA API token set to: %s", maskedToken)
					m.commandError = ""
				}
				m.textInput.SetValue("")
				return m, nil
			}

			if validatedCmd.Action == "jira_set_project" {
				key, value := commands.GetJiraConfigValue(validatedCmd.Name)
				if key == "" || value == "" {
					m.commandError = "Invalid JIRA project command"
					m.textInput.SetValue("")
					return m, nil
				}

				err := config.SetJiraConfig("project", value)
				if err != nil {
					m.commandError = fmt.Sprintf("Failed to set JIRA project: %v", err)
				} else {
					m.successMessage = fmt.Sprintf("✓ JIRA project key(s) set to: %s", value)
					m.commandError = ""
				}
				m.textInput.SetValue("")
				return m, nil
			}

			m.command = validatedCmd.Action
			m.submitted = true
			m.commandError = ""
			m.successMessage = ""
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	if !m.showingHelp {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if m.showingHelp {
		helpText := commands.GetHelpText()
		helpContent := fmt.Sprintf("%s\n\nPress any key to return to command prompt", helpText)

		centerStyle := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1, 2)

		return centerStyle.Render(styles.NeutralStyle.Render(helpContent))
	}

	welcomeText := "Welcome to Project Void"
	welcomeStyled := styles.WelcomeStyle.Render(welcomeText)

	description := "Enter commands below to navigate and interact with the system"
	descriptionStyled := styles.NeutralStyle.Render(description)

	var inputSection string
	if m.submitted {
		inputSection = styles.NeutralStyle.Render(fmt.Sprintf("Command processed: %s\n\nNavigating...", m.command))
	} else if m.commandError != "" {
		errorText := fmt.Sprintf("Error: %s\n\n", m.commandError)
		inputSection = fmt.Sprintf("%sCommand: %s\n\n%s\n\n%s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errorText),
			lipgloss.NewStyle().Align(lipgloss.Center).Render(m.textInput.View()),
			styles.NeutralStyle.Render("Try 'help' to see available commands"),
			styles.QuitStyle.Render("Press Ctrl+C or Esc to quit"))
	} else if m.successMessage != "" {
		successText := fmt.Sprintf("%s\n\n", m.successMessage)
		inputSection = fmt.Sprintf("%sCommand: %s\n\n%s\n\n%s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render(successText),
			lipgloss.NewStyle().Align(lipgloss.Center).Render(m.textInput.View()),
			styles.NeutralStyle.Render("Try 'help' to see available commands"),
			styles.QuitStyle.Render("Press Ctrl+C or Esc to quit"))
	} else {
		helpText := "Type 'help' to see available commands, 'start' to begin"
		inputSection = fmt.Sprintf("Command: %s\n\n%s\n\n%s",
			lipgloss.NewStyle().Align(lipgloss.Center).Render(m.textInput.View()),
			styles.NeutralStyle.Render(helpText),
			styles.QuitStyle.Render("Press Ctrl+C or Esc to quit"))
	}

	content := fmt.Sprintf("%s\n\n%s\n\n%s", welcomeStyled, descriptionStyled, inputSection)

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
	m.commandError = ""
	m.successMessage = ""
	m.textInput.SetValue("")
}
