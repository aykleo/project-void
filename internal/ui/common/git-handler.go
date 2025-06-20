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

		if len(gitConfig.Git.RepoURLs) == 0 {
			return &CommandResult{
				Action:  "git_status",
				Success: true,
				Message: "No Git repositories configured\nUse 'git repo <url-or-path>' to add repositories",
			}
		}

		var status strings.Builder
		status.WriteString(fmt.Sprintf("Git Repositories (%d configured):\n", len(gitConfig.Git.RepoURLs)))
		for i, repo := range gitConfig.Git.RepoURLs {
			status.WriteString(fmt.Sprintf("  %d. %s\n", i+1, repo))
		}
		status.WriteString(fmt.Sprintf("Type: %s", gitConfig.Git.RepoType))

		return &CommandResult{
			Action:  "git_status",
			Success: true,
			Message: status.String(),
		}

	case "git_clear_repo":
		err := config.ClearGitRepo()
		if err != nil {
			return &CommandResult{
				Action:  "git_clear_repo",
				Success: false,
				Message: fmt.Sprintf("Failed to clear Git repositories: %v", err),
			}
		}

		return &CommandResult{
			Action:  "git_clear_repo",
			Success: true,
			Message: "✓ All Git repositories cleared",
			Data:    map[string]interface{}{"repoURLs": []string{}, "repoType": "local"},
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
				Message: fmt.Sprintf("Failed to add Git repository: %v", err),
			}
		}

		repoType := "local"
		if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "git@") {
			repoType = "remote"
		}

		gitConfig, _ := config.LoadUserConfig()
		repoURLs := gitConfig.Git.RepoURLs

		return &CommandResult{
			Action:  "git_set_repo",
			Success: true,
			Message: fmt.Sprintf("✓ Git repository added: %s (%s)\nTotal repositories: %d", value, repoType, len(repoURLs)),
			Data:    map[string]interface{}{"repoURLs": repoURLs, "repoType": repoType},
		}

	case "git_remove_repo":
		key, value := commands.GetGitConfigValue(cmd.Name)
		if key == "" || value == "" {
			return &CommandResult{
				Action:  "git_remove_repo",
				Success: false,
				Message: "Invalid git repo remove command",
			}
		}

		err := config.RemoveGitRepo(value)
		if err != nil {
			return &CommandResult{
				Action:  "git_remove_repo",
				Success: false,
				Message: fmt.Sprintf("Failed to remove Git repository: %v", err),
			}
		}

		gitConfig, _ := config.LoadUserConfig()
		repoURLs := gitConfig.Git.RepoURLs

		return &CommandResult{
			Action:  "git_remove_repo",
			Success: true,
			Message: fmt.Sprintf("✓ Git repository removed: %s\nRemaining repositories: %d", value, len(repoURLs)),
			Data:    map[string]interface{}{"repoURLs": repoURLs, "repoType": gitConfig.Git.RepoType},
		}

	case "git_list_repos":
		gitConfig, err := config.LoadUserConfig()
		if err != nil {
			return &CommandResult{
				Action:  "git_list_repos",
				Success: false,
				Message: fmt.Sprintf("Failed to load Git config: %v", err),
			}
		}

		if len(gitConfig.Git.RepoURLs) == 0 {
			return &CommandResult{
				Action:  "git_list_repos",
				Success: true,
				Message: "No Git repositories configured",
			}
		}

		var status strings.Builder
		status.WriteString(fmt.Sprintf("Configured Git Repositories (%d):\n", len(gitConfig.Git.RepoURLs)))
		for i, repo := range gitConfig.Git.RepoURLs {
			status.WriteString(fmt.Sprintf("  %d. %s\n", i+1, repo))
		}

		return &CommandResult{
			Action:  "git_list_repos",
			Success: true,
			Message: status.String(),
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
