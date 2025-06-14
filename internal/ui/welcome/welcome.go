package welcome

import (
	"fmt"
	"project-void/internal/ui/styles"

	"project-void/internal/commands"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textInput    textinput.Model
	width        int
	height       int
	command      string
	commandError string
	submitted    bool
	showingHelp  bool
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
				m.textInput.SetValue("")
				return m, nil
			}

			m.command = validatedCmd.Action
			m.submitted = true
			m.commandError = ""
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
	m.textInput.SetValue("")
}
