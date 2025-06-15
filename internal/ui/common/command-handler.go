package common

import (
	"fmt"
	"project-void/internal/commands"
	"project-void/internal/config"
	"project-void/internal/ui/styles"

	"strings"

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

func (h CommandHandler) Update(msg tea.Msg) (CommandHandler, tea.Cmd, *CommandResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > 60 {
			h.textInput.Width = 50
		} else {
			h.textInput.Width = msg.Width - 10
		}
		return h, nil, nil

	case tea.KeyMsg:
		if !h.enabled {
			return h, nil, nil
		}

		if h.showingHelp {
			h.showingHelp = false
			h.commandError = ""
			return h, nil, nil
		}

		if h.showingCommand {
			switch msg.Type {
			case tea.KeyEnter:
				return h.processCommand()
			case tea.KeyCtrlC:
				return h, tea.Quit, &CommandResult{ShouldQuit: true}
			case tea.KeyEsc:
				h.showingCommand = false
				h.commandError = ""
				h.textInput.SetValue("")
				return h, nil, nil
			}

			if msg.String() == "'" {
				h.showingCommand = false
				h.commandError = ""
				h.textInput.SetValue("")
				return h, nil, nil
			}

			var cmd tea.Cmd
			h.textInput, cmd = h.textInput.Update(msg)
			return h, cmd, nil
		}

		if msg.String() == "c" {
			h.showingCommand = true
			h.textInput.Focus()
			return h, nil, nil
		}

		switch msg.Type {
		case tea.KeyEnter:
			return h.processCommand()
		case tea.KeyCtrlC, tea.KeyEsc:
			return h, tea.Quit, &CommandResult{ShouldQuit: true}
		}
	}

	if h.enabled && !h.showingHelp && !h.showingCommand {
		var cmd tea.Cmd
		h.textInput, cmd = h.textInput.Update(msg)
		return h, cmd, nil
	}

	return h, nil, nil
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
		"start": "home",
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

func (h CommandHandler) handleJiraCommands(cmd commands.Command) *CommandResult {
	switch cmd.Action {
	case "jira_status":
		status, err := config.GetJiraConfigStatus()
		if err != nil {
			return &CommandResult{
				Action:  "jira_status",
				Success: false,
				Message: fmt.Sprintf("Failed to get JIRA status: %v", err),
			}
		}
		return &CommandResult{
			Action:  "jira_status",
			Success: true,
			Message: status,
		}

	case "jira_set_url":
		key, value := commands.GetJiraConfigValue(cmd.Name)
		if key == "" || value == "" {
			return &CommandResult{
				Action:  "jira_set_url",
				Success: false,
				Message: "Invalid JIRA URL command",
			}
		}

		err := config.SetJiraConfig("url", value)
		if err != nil {
			return &CommandResult{
				Action:  "jira_set_url",
				Success: false,
				Message: fmt.Sprintf("Failed to set JIRA URL: %v", err),
			}
		}
		return &CommandResult{
			Action:  "jira_set_url",
			Success: true,
			Message: fmt.Sprintf("✓ JIRA URL set to: %s", value),
		}

	case "jira_set_user":
		key, value := commands.GetJiraConfigValue(cmd.Name)
		if key == "" || value == "" {
			return &CommandResult{
				Action:  "jira_set_user",
				Success: false,
				Message: "Invalid JIRA user command",
			}
		}

		err := config.SetJiraConfig("user", value)
		if err != nil {
			return &CommandResult{
				Action:  "jira_set_user",
				Success: false,
				Message: fmt.Sprintf("Failed to set JIRA user: %v", err),
			}
		}
		return &CommandResult{
			Action:  "jira_set_user",
			Success: true,
			Message: fmt.Sprintf("✓ JIRA username set to: %s", value),
		}

	case "jira_set_token":
		key, value := commands.GetJiraConfigValue(cmd.Name)
		if key == "" || value == "" {
			return &CommandResult{
				Action:  "jira_set_token",
				Success: false,
				Message: "Invalid JIRA token command",
			}
		}

		err := config.SetJiraConfig("token", value)
		if err != nil {
			return &CommandResult{
				Action:  "jira_set_token",
				Success: false,
				Message: fmt.Sprintf("Failed to set JIRA token: %v", err),
			}
		}

		maskedToken := value
		if len(maskedToken) > 8 {
			maskedToken = maskedToken[:4] + "..." + maskedToken[len(maskedToken)-4:]
		}
		return &CommandResult{
			Action:  "jira_set_token",
			Success: true,
			Message: fmt.Sprintf("✓ JIRA API token set to: %s", maskedToken),
		}

	case "jira_set_project":
		key, value := commands.GetJiraConfigValue(cmd.Name)
		if key == "" || value == "" {
			return &CommandResult{
				Action:  "jira_set_project",
				Success: false,
				Message: "Invalid JIRA project command",
			}
		}

		err := config.SetJiraConfig("project", value)
		if err != nil {
			return &CommandResult{
				Action:  "jira_set_project",
				Success: false,
				Message: fmt.Sprintf("Failed to set JIRA project: %v", err),
			}
		}
		return &CommandResult{
			Action:  "jira_set_project",
			Success: true,
			Message: fmt.Sprintf("✓ JIRA project key(s) set to: %s", value),
		}
	}

	return nil
}

func (h CommandHandler) handleGitCommands(cmd commands.Command) *CommandResult {
	switch cmd.Action {
	case "git_status":
		gitConfig, err := config.LoadUserConfig()
		if err != nil {
			return &CommandResult{
				Action:  "git_status",
				Success: false,
				Message: fmt.Sprintf("Failed to load Git config: %v", err),
			}
		}

		if gitConfig.Git.RepoURL == "" {
			return &CommandResult{
				Action:  "git_status",
				Success: true,
				Message: "No Git repository configured\nUse 'git repo <url-or-path>' to set a repository",
			}
		}

		status := fmt.Sprintf("Git Repository: %s\nType: %s", gitConfig.Git.RepoURL, gitConfig.Git.RepoType)
		return &CommandResult{
			Action:  "git_status",
			Success: true,
			Message: status,
		}

	case "git_clear_repo":
		err := config.ClearGitRepo()
		if err != nil {
			return &CommandResult{
				Action:  "git_clear_repo",
				Success: false,
				Message: fmt.Sprintf("Failed to clear Git repository: %v", err),
			}
		}

		return &CommandResult{
			Action:  "git_clear_repo",
			Success: true,
			Message: "✓ Git repository configuration cleared",
			Data:    map[string]interface{}{"repoURL": "", "repoType": ""},
		}

	case "git_set_repo":
		key, value := commands.GetGitConfigValue(cmd.Name)
		if key == "" || value == "" {
			return &CommandResult{
				Action:  "git_set_repo",
				Success: false,
				Message: "Invalid git repo command",
			}
		}

		err := config.SetGitConfig("repo", value)
		if err != nil {
			return &CommandResult{
				Action:  "git_set_repo",
				Success: false,
				Message: fmt.Sprintf("Failed to set Git repository: %v", err),
			}
		}

		repoType := "local"
		if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "git@") {
			repoType = "remote"
		}

		return &CommandResult{
			Action:  "git_set_repo",
			Success: true,
			Message: fmt.Sprintf("✓ Git repository set to: %s (%s)", value, repoType),
			Data:    map[string]interface{}{"repoURL": value, "repoType": repoType},
		}

	case "git_set_token":
		key, value := commands.GetGitConfigValue(cmd.Name)
		if key == "" || value == "" {
			return &CommandResult{
				Action:  "git_set_token",
				Success: false,
				Message: "Invalid git token command",
			}
		}

		err := config.SetGitConfig("token", value)
		if err != nil {
			return &CommandResult{
				Action:  "git_set_token",
				Success: false,
				Message: fmt.Sprintf("Failed to set GitHub token: %v", err),
			}
		}

		maskedToken := value
		if len(maskedToken) > 8 {
			maskedToken = maskedToken[:4] + "..." + maskedToken[len(maskedToken)-4:]
		}

		return &CommandResult{
			Action:  "git_set_token",
			Success: true,
			Message: fmt.Sprintf("✓ GitHub API token set: %s\nRate limit increased from 60 to 5,000 requests per hour!", maskedToken),
		}
	}

	return nil
}

func (h CommandHandler) RenderCommandInput(width int) string {
	if h.showingCommand {
		if h.commandError != "" {
			errorText := fmt.Sprintf("Error: %s", h.commandError)
			return fmt.Sprintf("%s\n%s\n%s",
				lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errorText),
				h.textInput.View(),
				lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press ' to exit command mode, esc to cancel"))
		} else {
			return fmt.Sprintf("%s\n%s",
				h.textInput.View(),
				lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press ' to exit command mode, esc to cancel"))
		}
	}

	return ""
}

func (h CommandHandler) RenderCommandPrompt(helpText string) string {
	if h.commandError != "" {
		errorText := fmt.Sprintf("Error: %s\n\n", h.commandError)
		return fmt.Sprintf("%s%s\n\n%s\n\n%s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errorText),
			lipgloss.NewStyle().Align(lipgloss.Center).Render(h.textInput.View()),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(helpText),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press c for commands, Ctrl+C or Esc to quit"))
	} else if h.successMessage != "" {
		successText := fmt.Sprintf("%s\n\n", h.successMessage)
		return fmt.Sprintf("%s%s\n\n%s\n\n%s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render(successText),
			lipgloss.NewStyle().Align(lipgloss.Center).Render(h.textInput.View()),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(helpText),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press c for commands, Ctrl+C or Esc to quit"))
	} else {
		return fmt.Sprintf("%s\n\n%s\n\n%s",
			lipgloss.NewStyle().Align(lipgloss.Center).Render(h.textInput.View()),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(helpText),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press c for commands, Ctrl+C or Esc to quit"))
	}
}

func (h CommandHandler) RenderHelp(width, height int) string {
	if h.showingHelp {
		helpText := commands.GetHelpText()
		helpContent := fmt.Sprintf("%s\n\nPress any key to return", helpText)

		centerStyle := lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1, 2)

		return centerStyle.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(helpContent))
	}
	return ""
}

func (h CommandHandler) IsShowingHelp() bool {
	return h.showingHelp
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
		"start": "home",
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

	if validatedCmd.Action == "filter_by_author" || validatedCmd.Action == "clear_author_filter" {
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

func (h StatisticsCommandHandler) Update(msg tea.Msg) (StatisticsCommandHandler, tea.Cmd, *CommandResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > 60 {
			h.textInput.Width = 50
		} else {
			h.textInput.Width = msg.Width - 10
		}
		return h, nil, nil

	case tea.KeyMsg:
		if !h.enabled {
			return h, nil, nil
		}

		if h.showingHelp {
			h.showingHelp = false
			h.commandError = ""
			return h, nil, nil
		}

		if h.showingCommand {
			switch msg.Type {
			case tea.KeyEnter:
				return h.processCommand()
			case tea.KeyCtrlC:
				return h, tea.Quit, &CommandResult{ShouldQuit: true}
			case tea.KeyEsc:
				h.showingCommand = false
				h.commandError = ""
				h.textInput.SetValue("")
				return h, nil, nil
			}

			if msg.String() == "'" {
				h.showingCommand = false
				h.commandError = ""
				h.textInput.SetValue("")
				return h, nil, nil
			}

			var cmd tea.Cmd
			h.textInput, cmd = h.textInput.Update(msg)
			return h, cmd, nil
		}

		if msg.String() == "c" {
			h.showingCommand = true
			h.textInput.Focus()
			return h, nil, nil
		}

		switch msg.Type {
		case tea.KeyEnter:
			return h.processCommand()
		case tea.KeyCtrlC, tea.KeyEsc:
			return h, tea.Quit, &CommandResult{ShouldQuit: true}
		}
	}

	if h.enabled && !h.showingHelp && !h.showingCommand {
		var cmd tea.Cmd
		h.textInput, cmd = h.textInput.Update(msg)
		return h, cmd, nil
	}

	return h, nil, nil
}
