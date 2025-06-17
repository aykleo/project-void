package git

import (
	"strings"
)

func isRemoteURL(input string) bool {
	return strings.HasPrefix(input, "http://") ||
		strings.HasPrefix(input, "https://") ||
		strings.HasPrefix(input, "git@")
}

func detectProvider(repoURL string) (GitProvider, error) {
	if strings.Contains(repoURL, "github.com") {
		provider := NewGitHubProvider()

		if gitConfig, err := LoadGitConfig(); err == nil && gitConfig.GitHubToken != "" {
			provider.SetToken(gitConfig.GitHubToken)
		}

		return provider, nil
	}

	provider := NewGitHubProvider()

	if gitConfig, err := LoadGitConfig(); err == nil && gitConfig.GitHubToken != "" {
		provider.SetToken(gitConfig.GitHubToken)
	}

	return provider, nil
}
