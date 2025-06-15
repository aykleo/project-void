package common

import (
	"fmt"
	"project-void/internal/commands"
	"project-void/internal/config"
	"strings"
)

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
