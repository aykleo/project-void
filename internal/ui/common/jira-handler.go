package common

import (
	"fmt"
	"project-void/internal/commands"
	"project-void/internal/config"
)

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

	case "jira_filter_on":
		err := config.SetJiraConfig("filter", "true")
		if err != nil {
			return &CommandResult{
				Action:  "jira_filter_on",
				Success: false,
				Message: fmt.Sprintf("Failed to enable user filtering: %v", err),
			}
		}
		return &CommandResult{
			Action:  "jira_filter_on",
			Success: true,
			Message: "✓ User filtering enabled - showing only your issues",
		}

	case "jira_filter_off":
		err := config.SetJiraConfig("filter", "false")
		if err != nil {
			return &CommandResult{
				Action:  "jira_filter_off",
				Success: false,
				Message: fmt.Sprintf("Failed to disable user filtering: %v", err),
			}
		}
		return &CommandResult{
			Action:  "jira_filter_off",
			Success: true,
			Message: "✓ User filtering disabled - showing all issues",
		}
	}

	return nil
}
