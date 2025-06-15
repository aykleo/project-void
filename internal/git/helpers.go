package git

import (
	"strings"
	"time"
)

func ToMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

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
