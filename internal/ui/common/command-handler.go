package common

import (
	"fmt"
	"project-void/internal/commands"
	"project-void/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CommandResult struct {
	Action         string
	Success        bool
	Message        string
	ShouldQuit     bool
	ShouldNavigate bool
	NavigateTo     string
	Data           map[string]interface{}
}

type CommandHandler struct {
	textInput      textinput.Model
	commandError   string
	successMessage string
	showingHelp    bool
	showingGitHelp bool
	showingCommand bool
	enabled        bool
}

func NewCommandHandler(placeholder string) CommandHandler {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.HighlightColor)

	return CommandHandler{
		textInput: ti,
		enabled:   true,
	}
}

func (h CommandHandler) Init() tea.Cmd {
	return textinput.Blink
}

func (h CommandHandler) processCommand() (CommandHandler, tea.Cmd, *CommandResult) {
	input := h.textInput.Value()
	command, err := commands.ValidateCommand(input)
	if err != nil {
		h.commandError = err.Error()
		h.successMessage = ""
		h.textInput.SetValue("")
		return h, nil, nil
	}

	navigationCommands := map[string]string{
		"start": "statistics",
		"reset": "welcome",
	}

	if navigateTo, isNavigation := navigationCommands[command.Action]; isNavigation {
		h.commandError = ""
		h.successMessage = ""
		h.textInput.SetValue("")
		return h, nil, &CommandResult{
			Action:         command.Action,
			Success:        true,
			ShouldNavigate: true,
			NavigateTo:     navigateTo,
		}
	}

	if result := h.handleJiraCommands(command); result != nil {
		h.commandError = ""
		if result.Success {
			h.successMessage = result.Message
		} else {
			h.commandError = result.Message
		}
		h.textInput.SetValue("")
		return h, nil, result
	}

	if result := h.handleGitCommands(command); result != nil {
		h.commandError = ""
		if result.Success {
			h.successMessage = result.Message
		} else {
			h.commandError = result.Message
		}
		h.textInput.SetValue("")
		return h, nil, result
	}

	if command.Action == "help" || command.Name == "help" {
		h.showingHelp = true
		h.showingCommand = false
		h.commandError = ""
		h.successMessage = ""
		h.textInput.SetValue("")
		return h, nil, &CommandResult{
			Action:  "help",
			Success: true,
		}
	}

	if command.Action == "git_help" {
		h.showingGitHelp = true
		h.commandError = ""
		h.successMessage = ""
		h.textInput.SetValue("")
		return h, nil, &CommandResult{
			Action:  "git_help",
			Success: true,
		}
	}

	if command.Action == "quit" {
		h.textInput.SetValue("")
		return h, tea.Quit, &CommandResult{
			Action:     "quit",
			ShouldQuit: true,
		}
	}

	if command.Action == "void_set_date" {
		date, err := commands.GetDateFromCommand(command.Name)
		if err != nil {
			h.commandError = fmt.Sprintf("Error parsing date: %v", err)
			h.successMessage = ""
			h.textInput.SetValue("")
			return h, nil, nil
		}

		h.commandError = ""
		h.successMessage = fmt.Sprintf("✓ Analysis date set to: %s", date.Format("January 2, 2006"))
		h.textInput.SetValue("")
		return h, nil, &CommandResult{
			Action:  "void_set_date",
			Success: true,
			Message: h.successMessage,
			Data:    map[string]interface{}{"date": date},
		}
	}

	h.commandError = ""
	h.successMessage = ""
	h.textInput.SetValue("")

	return h, nil, &CommandResult{
		Action:         command.Action,
		Success:        true,
		ShouldNavigate: false,
		Data:           map[string]interface{}{"command": command},
	}
}

func (h CommandHandler) IsShowingHelp() bool {
	return h.showingHelp
}

func (h CommandHandler) IsShowingGitHelp() bool {
	return h.showingGitHelp
}

func (h CommandHandler) IsShowingCommand() bool {
	return h.showingCommand
}

func (h CommandHandler) IsEnabled() bool {
	return h.enabled
}

func (h *CommandHandler) SetEnabled(enabled bool) {
	h.enabled = enabled
	if !enabled {
		h.showingCommand = false
		h.showingHelp = false
	}
}

func (h *CommandHandler) ClearMessages() {
	h.commandError = ""
	h.successMessage = ""
}

func (h *CommandHandler) SetError(message string) {
	h.commandError = message
	h.successMessage = ""
}

func (h *CommandHandler) SetSuccess(message string) {
	h.successMessage = message
	h.commandError = ""
}
