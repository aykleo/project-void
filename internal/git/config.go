package git

import (
	"fmt"
	"os"
	"project-void/internal/config"
	"strings"
)

type GitConfig struct {
	RepoURL     string `json:"repo_url"`
	RepoType    string `json:"repo_type"`
	GitHubToken string `json:"github_token,omitempty"`
}

func LoadGitConfig() (*GitConfig, error) {
	userConfig, err := config.LoadUserConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load user config: %w", err)
	}

	gitConfig := &GitConfig{
		RepoURL:     userConfig.Git.RepoURL,
		RepoType:    userConfig.Git.RepoType,
		GitHubToken: userConfig.Git.GitHubToken,
	}

	if gitConfig.RepoURL == "" {
		gitConfig.RepoURL = os.Getenv("GIT_REPO_URL")
	}
	if gitConfig.RepoType == "" {
		if gitConfig.RepoURL != "" {
			if isRemoteURL(gitConfig.RepoURL) {
				gitConfig.RepoType = "remote"
			} else {
				gitConfig.RepoType = "local"
			}
		}
	}

	return gitConfig, nil
}

func SetGitRepo(repoURL string) error {
	userConfig, err := config.LoadUserConfig()
	if err != nil {
		return err
	}

	repoType := "local"
	if isRemoteURL(repoURL) {
		repoType = "remote"
	}

	userConfig.Git.RepoURL = repoURL
	userConfig.Git.RepoType = repoType

	return config.SaveUserConfig(userConfig)
}

func GetGitStatus() (string, error) {
	gitConfig, err := LoadGitConfig()
	if err != nil {
		return "", err
	}

	if gitConfig.RepoURL == "" {
		return "No Git repository configured", nil
	}

	status := fmt.Sprintf("Git Repository: %s\nType: %s", gitConfig.RepoURL, gitConfig.RepoType)
	return status, nil
}

func ShouldEnableDevMode() bool {
	gitConfig, err := LoadGitConfig()
	if err != nil {
		return false
	}
	return gitConfig.RepoURL != ""
}

func GetConfiguredRepoSource() string {
	gitConfig, err := LoadGitConfig()
	if err != nil {
		return ""
	}
	return gitConfig.RepoURL
}

func GetGitConfigStatus() (string, error) {
	gitConfig, err := LoadGitConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load git config: %w", err)
	}

	var status strings.Builder
	status.WriteString("Current Git Configuration:\n\n")

	if gitConfig.RepoURL != "" {
		status.WriteString(fmt.Sprintf("  Repository: %s\n", gitConfig.RepoURL))
		status.WriteString(fmt.Sprintf("  Type: %s\n", gitConfig.RepoType))
	} else {
		status.WriteString("  Repository: (not set)\n")
		status.WriteString("  Type: (not set)\n")
	}

	if gitConfig.GitHubToken != "" {
		masked := gitConfig.GitHubToken
		if len(masked) > 8 {
			masked = masked[:4] + "..." + masked[len(masked)-4:]
		}
		status.WriteString(fmt.Sprintf("  GitHub Token: %s\n", masked))
		status.WriteString("  Rate Limit: 5,000 requests/hour (authenticated)\n")
	} else {
		status.WriteString("  GitHub Token: (not set)\n")
		status.WriteString("  Rate Limit: 60 requests/hour (unauthenticated)\n")
	}

	if gitConfig.RepoURL != "" {
		status.WriteString("\n✓ Git repository is configured!")
		if gitConfig.GitHubToken == "" && strings.Contains(gitConfig.RepoURL, "github.com") {
			status.WriteString("\n💡 Tip: Set a GitHub token with 'git token <token>' for higher rate limits")
		}
	} else {
		status.WriteString("\n⚠ No Git repository configured. Use 'git repo <url>' to set one.")
	}

	return status.String(), nil
}
