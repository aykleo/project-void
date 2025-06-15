package git

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Commit struct {
	Hash      string
	Branch    string
	Author    string
	Message   string
	Timestamp time.Time
}

type GitProvider interface {
	GetCommitsSince(repoURL string, since time.Time) ([]Commit, error)
	GetCommitsSinceByAuthors(repoURL string, since time.Time, authorNames []string) ([]Commit, error)
}

func ToMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func isRemoteURL(input string) bool {
	return strings.HasPrefix(input, "http://") ||
		strings.HasPrefix(input, "https://") ||
		strings.HasPrefix(input, "git@")
}

func GetCommitsSince(repoPathOrURL string, since time.Time) ([]Commit, error) {
	if isRemoteURL(repoPathOrURL) {
		provider, err := detectProvider(repoPathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to detect Git provider: %w", err)
		}
		return provider.GetCommitsSince(repoPathOrURL, since)
	}

	return getCommitsSinceLocal(repoPathOrURL, since)
}

func GetCommitsSinceByAuthors(repoPathOrURL string, since time.Time, authorNames []string) ([]Commit, error) {
	if isRemoteURL(repoPathOrURL) {
		provider, err := detectProvider(repoPathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to detect Git provider: %w", err)
		}
		return provider.GetCommitsSinceByAuthors(repoPathOrURL, since, authorNames)
	}

	return getCommitsSinceByAuthorsLocal(repoPathOrURL, since, authorNames)
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

func getCommitsSinceLocal(repoPath string, since time.Time) ([]Commit, error) {
	since = ToMidnight(since)

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	refs, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed to get references: %w", err)
	}

	uniqueCommits := make(map[string]Commit)

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() == plumbing.HEAD {
			return nil
		}

		commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			return nil
		}
		defer commitIter.Close()

		err = commitIter.ForEach(func(c *object.Commit) error {
			if c.Author.When.UTC().After(since) {
				hash := c.Hash.String()
				if _, exists := uniqueCommits[hash]; !exists {
					uniqueCommits[hash] = Commit{
						Hash:      hash,
						Branch:    ref.Name().Short(),
						Author:    c.Author.Name,
						Message:   c.Message,
						Timestamp: c.Author.When,
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating references: %w", err)
	}

	commits := make([]Commit, 0, len(uniqueCommits))
	for _, commit := range uniqueCommits {
		commits = append(commits, commit)
	}

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.After(commits[j].Timestamp)
	})

	return commits, nil
}

func getCommitsSinceByAuthorsLocal(repoPath string, since time.Time, authorNames []string) ([]Commit, error) {
	since = ToMidnight(since)

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	refs, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed to get references: %w", err)
	}

	uniqueCommits := make(map[string]Commit)

	var lowerAuthorNames []string
	for _, name := range authorNames {
		lowerAuthorNames = append(lowerAuthorNames, strings.ToLower(name))
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() == plumbing.HEAD {
			return nil
		}

		commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			return nil
		}
		defer commitIter.Close()

		err = commitIter.ForEach(func(c *object.Commit) error {
			if c.Author.When.UTC().After(since) {
				authorLower := strings.ToLower(c.Author.Name)
				matchesAuthor := false

				for _, targetAuthor := range lowerAuthorNames {
					if strings.Contains(authorLower, targetAuthor) || strings.Contains(targetAuthor, authorLower) {
						matchesAuthor = true
						break
					}
				}

				if matchesAuthor {
					hash := c.Hash.String()
					if _, exists := uniqueCommits[hash]; !exists {
						uniqueCommits[hash] = Commit{
							Hash:      hash,
							Branch:    ref.Name().Short(),
							Author:    c.Author.Name,
							Message:   c.Message,
							Timestamp: c.Author.When,
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating references: %w", err)
	}

	commits := make([]Commit, 0, len(uniqueCommits))
	for _, commit := range uniqueCommits {
		commits = append(commits, commit)
	}

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.After(commits[j].Timestamp)
	})

	return commits, nil
}
