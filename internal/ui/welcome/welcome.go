package welcome

import (
	"fmt"
	"project-void/internal/ui/common"
	"project-void/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	commandHandler common.CommandHandler
	width          int
	height         int
	command        string
	submitted      bool
}

func InitialModel() Model {
	return Model{
		commandHandler: common.NewCommandHandler("Enter a command (e.g., help)..."),
		submitted:      false,
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
	}

	return m, cmd
}

func (m Model) View() string {
	if helpView := m.commandHandler.RenderHelp(m.width, m.height); helpView != "" {
		return helpView
	}

	welcomeText := "Welcome to Project Void"
	welcomeStyled := styles.WelcomeStyle.Render(welcomeText)

	description := "Enter commands below to navigate and interact with the system"
	descriptionStyled := styles.NeutralStyle.Render(description)

	var inputSection string
	if m.submitted {
		inputSection = styles.NeutralStyle.Render(fmt.Sprintf("Command processed: %s\n\nNavigating...", m.command))
	} else {
		helpText := "Try 'help' to see available commands"
		inputSection = m.commandHandler.RenderCommandPrompt(helpText)
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
	m.commandHandler.ClearMessages()
}
