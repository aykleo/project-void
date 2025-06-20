package git

import (
	"fmt"
	"os"
	"project-void/internal/config"
	"strings"
)

type GitConfig struct {
	RepoURLs    []string `json:"repo_urls"`
	RepoType    string   `json:"repo_type"`
	GitHubToken string   `json:"github_token,omitempty"`
}

func LoadGitConfig() (*GitConfig, error) {
	userConfig, err := config.LoadUserConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load user config: %w", err)
	}

	gitConfig := &GitConfig{
		RepoURLs:    userConfig.Git.RepoURLs,
		RepoType:    userConfig.Git.RepoType,
		GitHubToken: userConfig.Git.GitHubToken,
	}

	if len(gitConfig.RepoURLs) == 0 {
		if envRepo := os.Getenv("GIT_REPO_URL"); envRepo != "" {
			gitConfig.RepoURLs = []string{envRepo}
		}
	}
	if gitConfig.RepoType == "" {
		if len(gitConfig.RepoURLs) > 0 {
			hasRemote := false
			for _, repo := range gitConfig.RepoURLs {
				if isRemoteURL(repo) {
					hasRemote = true
					break
				}
			}
			if hasRemote {
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

	for _, existing := range userConfig.Git.RepoURLs {
		if existing == repoURL {
			return nil
		}
	}

	repoType := "local"
	if isRemoteURL(repoURL) {
		repoType = "remote"
	}

	userConfig.Git.RepoURLs = append(userConfig.Git.RepoURLs, repoURL)
	if repoType == "remote" {
		userConfig.Git.RepoType = "remote"
	} else if userConfig.Git.RepoType != "remote" {
		userConfig.Git.RepoType = "local"
	}

	return config.SaveUserConfig(userConfig)
}

func GetGitStatus() (string, error) {
	gitConfig, err := LoadGitConfig()
	if err != nil {
		return "", err
	}

	if len(gitConfig.RepoURLs) == 0 {
		return "No Git repositories configured", nil
	}

	var status strings.Builder
	status.WriteString(fmt.Sprintf("Git Repositories (%d configured):\n", len(gitConfig.RepoURLs)))
	for i, repo := range gitConfig.RepoURLs {
		status.WriteString(fmt.Sprintf("  %d. %s\n", i+1, repo))
	}
	status.WriteString(fmt.Sprintf("Type: %s", gitConfig.RepoType))
	return status.String(), nil
}

func ShouldEnableDevMode() bool {
	gitConfig, err := LoadGitConfig()
	if err != nil {
		return false
	}
	return len(gitConfig.RepoURLs) > 0
}

func GetConfiguredRepoSource() string {
	gitConfig, err := LoadGitConfig()
	if err != nil {
		return ""
	}
	if len(gitConfig.RepoURLs) == 0 {
		return ""
	}
	return gitConfig.RepoURLs[0]
}

func GetGitConfigStatus() (string, error) {
	gitConfig, err := LoadGitConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load git config: %w", err)
	}

	var status strings.Builder
	status.WriteString("Current Git Configuration:\n\n")

	if len(gitConfig.RepoURLs) > 0 {
		status.WriteString(fmt.Sprintf("  Repositories (%d configured):\n", len(gitConfig.RepoURLs)))
		for i, repo := range gitConfig.RepoURLs {
			status.WriteString(fmt.Sprintf("    %d. %s\n", i+1, repo))
		}
		status.WriteString(fmt.Sprintf("  Type: %s\n", gitConfig.RepoType))
	} else {
		status.WriteString("  Repositories: (none configured)\n")
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

	if len(gitConfig.RepoURLs) > 0 {
		status.WriteString("\n✓ Git repositories are configured!")
		hasGitHub := false
		for _, repo := range gitConfig.RepoURLs {
			if strings.Contains(repo, "github.com") {
				hasGitHub = true
				break
			}
		}
		if gitConfig.GitHubToken == "" && hasGitHub {
			status.WriteString("\n💡 Tip: Set a GitHub token with 'git token <token>' for higher rate limits")
		}
	} else {
		status.WriteString("\n⚠ No Git repositories configured. Use 'git repo <url>' to add repositories.")
	}

	return status.String(), nil
}
