package common

import (
	"fmt"
	"project-void/internal/commands"

	tea "github.com/charmbracelet/bubbletea"
)

type StatisticsCommandHandler struct {
	CommandHandler
	selectedFolder string
	isDev          bool
}

func NewStatisticsCommandHandler(placeholder, selectedFolder string, isDev bool) StatisticsCommandHandler {
	return StatisticsCommandHandler{
		CommandHandler: NewCommandHandler(placeholder),
		selectedFolder: selectedFolder,
		isDev:          isDev,
	}
}

func (h StatisticsCommandHandler) processCommand() (StatisticsCommandHandler, tea.Cmd, *CommandResult) {
	inputValue := h.textInput.Value()
	validatedCmd, err := commands.ValidateCommand(inputValue)
	if err != nil {
		h.commandError = err.Error()
		h.textInput.SetValue("")
		return h, nil, nil
	}

	h.textInput.SetValue("")
	h.showingCommand = false

	if validatedCmd.Action == "help" {
		h.showingHelp = true
		h.commandError = ""
		h.successMessage = ""
		return h, nil, nil
	}

	if validatedCmd.Action == "quit" {
		return h, tea.Quit, &CommandResult{
			Action:     "quit",
			ShouldQuit: true,
		}
	}

	if validatedCmd.Action == "void_set_date" {
		date, err := commands.GetDateFromCommand(validatedCmd.Name)
		if err != nil {
			h.commandError = fmt.Sprintf("Error parsing date: %v", err)
			return h, nil, nil
		}

		return h, nil, &CommandResult{
			Action:  "void_set_date",
			Success: true,
			Message: fmt.Sprintf("✓ Analysis date set to: %s", date.Format("January 2, 2006")),
			Data:    map[string]interface{}{"date": date},
		}
	}

	navigationCommands := map[string]string{
		"start": "statistics",
		"reset": "welcome",
	}

	if navigateTo, isNavigation := navigationCommands[validatedCmd.Action]; isNavigation {
		return h, nil, &CommandResult{
			Action:         validatedCmd.Action,
			Success:        true,
			ShouldNavigate: true,
			NavigateTo:     navigateTo,
		}
	}

	if result := h.CommandHandler.handleJiraCommands(validatedCmd); result != nil {
		h.CommandHandler.commandError = ""
		h.CommandHandler.successMessage = result.Message
		return h, nil, result
	}

	if result := h.CommandHandler.handleGitCommands(validatedCmd); result != nil {
		h.CommandHandler.commandError = ""
		h.CommandHandler.successMessage = result.Message
		return h, nil, result
	}

	if validatedCmd.Action == "git_help" {
		h.showingGitHelp = true
		h.commandError = ""
		h.successMessage = ""
		return h, nil, nil
	}

	if validatedCmd.Action == "filter_by_author" || validatedCmd.Action == "clear_author_filter" || validatedCmd.Action == "filter_by_branch" || validatedCmd.Action == "clear_branch_filter" {
		return h, nil, &CommandResult{
			Action:  validatedCmd.Action,
			Success: true,
			Data:    map[string]interface{}{"command": validatedCmd},
		}
	}

	return h, nil, &CommandResult{
		Action:  validatedCmd.Action,
		Success: true,
		Data:    map[string]interface{}{"command": validatedCmd},
	}
}
