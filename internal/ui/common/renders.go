package common

import (
	"fmt"

	"project-void/internal/commands"

	"github.com/charmbracelet/lipgloss"
)

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
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("ctrl+c/esc: exit"))
	} else if h.successMessage != "" {
		successText := fmt.Sprintf("%s\n\n", h.successMessage)
		return fmt.Sprintf("%s%s\n\n%s\n\n%s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render(successText),
			lipgloss.NewStyle().Align(lipgloss.Center).Render(h.textInput.View()),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(helpText),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("ctrl+c/esc: exit"))
	} else {
		return fmt.Sprintf("%s\n\n%s\n\n%s",
			lipgloss.NewStyle().Align(lipgloss.Center).Render(h.textInput.View()),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(helpText),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("ctrl+c/esc: exit"))
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

func (h CommandHandler) RenderGitHelp(width, height int) string {
	if h.showingGitHelp {
		gitHelpText := commands.GetGitHelpText()
		gitHelpContent := fmt.Sprintf("%s\n\nPress any key to return", gitHelpText)

		centerStyle := lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1, 2)

		return centerStyle.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(gitHelpContent))
	}
	return ""
}

func (h CommandHandler) RenderJiraHelp(width, height int) string {
	if h.showingJiraHelp {
		jiraHelpText := commands.GetJiraHelpText()
		jiraHelpContent := fmt.Sprintf("%s\n\nPress any key to return", jiraHelpText)

		centerStyle := lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1, 2)

		return centerStyle.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(jiraHelpContent))
	}
	return ""
}
