package git

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"project-void/internal/helpers"

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
	RepoName  string
	RepoType  string
}

type GitProvider interface {
	GetCommitsSince(repoURL string, since time.Time) ([]Commit, error)
	GetCommitsSinceByAuthors(repoURL string, since time.Time, authorNames []string) ([]Commit, error)
	GetCommitsSinceByBranches(repoURL string, since time.Time, branchNames []string) ([]Commit, error)
	GetCommitsSinceByAuthorsAndBranches(repoURL string, since time.Time, authorNames []string, branchNames []string) ([]Commit, error)
}

func GetCommitsSince(repoPathOrURL string, since time.Time) ([]Commit, error) {
	var commits []Commit
	var err error

	if isRemoteURL(repoPathOrURL) {
		provider, err := detectProvider(repoPathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to detect Git provider: %w", err)
		}
		commits, err = provider.GetCommitsSince(repoPathOrURL, since)
	} else {
		commits, err = getCommitsSinceLocal(repoPathOrURL, since)
	}

	if err != nil {
		return nil, err
	}

	repoName := getRepoDisplayName(repoPathOrURL)
	repoType := "local"
	if isRemoteURL(repoPathOrURL) {
		repoType = "remote"
	}

	for i := range commits {
		commits[i].RepoName = repoName
		commits[i].RepoType = repoType
	}

	return commits, nil
}

func GetCommitsSinceByAuthors(repoPathOrURL string, since time.Time, authorNames []string) ([]Commit, error) {
	var commits []Commit
	var err error

	if isRemoteURL(repoPathOrURL) {
		provider, err := detectProvider(repoPathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to detect Git provider: %w", err)
		}
		commits, err = provider.GetCommitsSinceByAuthors(repoPathOrURL, since, authorNames)
	} else {
		commits, err = getCommitsSinceByAuthorsLocal(repoPathOrURL, since, authorNames)
	}

	if err != nil {
		return nil, err
	}

	repoName := getRepoDisplayName(repoPathOrURL)
	repoType := "local"
	if isRemoteURL(repoPathOrURL) {
		repoType = "remote"
	}

	for i := range commits {
		commits[i].RepoName = repoName
		commits[i].RepoType = repoType
	}

	return commits, nil
}

func GetCommitsSinceByAuthorsFromMultipleRepos(repoPathsOrURLs []string, since time.Time, authorNames []string) ([]Commit, error) {
	if len(repoPathsOrURLs) == 0 {
		return []Commit{}, nil
	}

	allCommits := make(map[string]Commit)
	var errors []string

	for _, repo := range repoPathsOrURLs {
		commits, err := GetCommitsSinceByAuthors(repo, since, authorNames)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to load commits from %s: %v", repo, err))
			continue
		}

		for _, commit := range commits {
			key := repo + ":" + commit.Hash
			allCommits[key] = commit
		}
	}

	result := make([]Commit, 0, len(allCommits))
	for _, commit := range allCommits {
		result = append(result, commit)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	if len(errors) > 0 && len(result) == 0 {
		return nil, fmt.Errorf("all repositories failed to load commits: %s", strings.Join(errors, "; "))
	}

	return result, nil
}

func GetCommitsSinceByBranches(repoPathOrURL string, since time.Time, branchNames []string) ([]Commit, error) {
	var commits []Commit
	var err error

	if isRemoteURL(repoPathOrURL) {
		provider, err := detectProvider(repoPathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to detect Git provider: %w", err)
		}
		commits, err = provider.GetCommitsSinceByBranches(repoPathOrURL, since, branchNames)
	} else {
		commits, err = getCommitsSinceByBranchesLocal(repoPathOrURL, since, branchNames)
	}

	if err != nil {
		return nil, err
	}

	repoName := getRepoDisplayName(repoPathOrURL)
	repoType := "local"
	if isRemoteURL(repoPathOrURL) {
		repoType = "remote"
	}

	for i := range commits {
		commits[i].RepoName = repoName
		commits[i].RepoType = repoType
	}

	return commits, nil
}

func GetCommitsSinceByBranchesFromMultipleRepos(repoPathsOrURLs []string, since time.Time, branchNames []string) ([]Commit, error) {
	if len(repoPathsOrURLs) == 0 {
		return []Commit{}, nil
	}

	allCommits := make(map[string]Commit)
	var errors []string

	for _, repo := range repoPathsOrURLs {
		commits, err := GetCommitsSinceByBranches(repo, since, branchNames)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to load commits from %s: %v", repo, err))
			continue
		}

		for _, commit := range commits {
			key := repo + ":" + commit.Hash
			allCommits[key] = commit
		}
	}

	result := make([]Commit, 0, len(allCommits))
	for _, commit := range allCommits {
		result = append(result, commit)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	if len(errors) > 0 && len(result) == 0 {
		return nil, fmt.Errorf("all repositories failed to load commits: %s", strings.Join(errors, "; "))
	}

	return result, nil
}

func GetCommitsSinceByAuthorsAndBranches(repoPathOrURL string, since time.Time, authorNames []string, branchNames []string) ([]Commit, error) {
	var commits []Commit
	var err error

	if isRemoteURL(repoPathOrURL) {
		provider, err := detectProvider(repoPathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to detect Git provider: %w", err)
		}
		commits, err = provider.GetCommitsSinceByAuthorsAndBranches(repoPathOrURL, since, authorNames, branchNames)
	} else {
		commits, err = getCommitsSinceByAuthorsAndBranchesLocal(repoPathOrURL, since, authorNames, branchNames)
	}

	if err != nil {
		return nil, err
	}

	repoName := getRepoDisplayName(repoPathOrURL)
	repoType := "local"
	if isRemoteURL(repoPathOrURL) {
		repoType = "remote"
	}

	for i := range commits {
		commits[i].RepoName = repoName
		commits[i].RepoType = repoType
	}

	return commits, nil
}

func GetCommitsSinceByAuthorsAndBranchesFromMultipleRepos(repoPathsOrURLs []string, since time.Time, authorNames []string, branchNames []string) ([]Commit, error) {
	if len(repoPathsOrURLs) == 0 {
		return []Commit{}, nil
	}

	allCommits := make(map[string]Commit)
	var errors []string

	for _, repo := range repoPathsOrURLs {
		commits, err := GetCommitsSinceByAuthorsAndBranches(repo, since, authorNames, branchNames)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to load commits from %s: %v", repo, err))
			continue
		}

		for _, commit := range commits {
			key := repo + ":" + commit.Hash
			allCommits[key] = commit
		}
	}

	result := make([]Commit, 0, len(allCommits))
	for _, commit := range allCommits {
		result = append(result, commit)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	if len(errors) > 0 && len(result) == 0 {
		return nil, fmt.Errorf("all repositories failed to load commits: %s", strings.Join(errors, "; "))
	}

	return result, nil
}

func getCommitsSinceLocal(repoPath string, since time.Time) ([]Commit, error) {
	since = helpers.ToMidnight(since)

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
	since = helpers.ToMidnight(since)

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

func getCommitsSinceByBranchesLocal(repoPath string, since time.Time, branchNames []string) ([]Commit, error) {
	since = helpers.ToMidnight(since)

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	refs, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed to get references: %w", err)
	}

	uniqueCommits := make(map[string]Commit)

	var lowerBranchNames []string
	for _, name := range branchNames {
		lowerBranchNames = append(lowerBranchNames, strings.ToLower(name))
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() == plumbing.HEAD {
			return nil
		}

		branchLower := strings.ToLower(ref.Name().Short())
		matchesBranch := false
		for _, targetBranch := range lowerBranchNames {
			if strings.Contains(branchLower, targetBranch) {
				matchesBranch = true
				break
			}
		}

		if !matchesBranch {
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

func getCommitsSinceByAuthorsAndBranchesLocal(repoPath string, since time.Time, authorNames []string, branchNames []string) ([]Commit, error) {
	since = helpers.ToMidnight(since)

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

	var lowerBranchNames []string
	for _, name := range branchNames {
		lowerBranchNames = append(lowerBranchNames, strings.ToLower(name))
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() == plumbing.HEAD {
			return nil
		}

		branchLower := strings.ToLower(ref.Name().Short())
		matchesBranch := false
		for _, targetBranch := range lowerBranchNames {
			if strings.Contains(branchLower, targetBranch) {
				matchesBranch = true
				break
			}
		}

		if !matchesBranch {
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

func GetCommitsSinceFromMultipleRepos(repoPathsOrURLs []string, since time.Time) ([]Commit, error) {
	if len(repoPathsOrURLs) == 0 {
		return []Commit{}, nil
	}

	allCommits := make(map[string]Commit)
	var errors []string

	for _, repo := range repoPathsOrURLs {
		commits, err := GetCommitsSince(repo, since)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to load commits from %s: %v", repo, err))
			continue
		}

		for _, commit := range commits {
			key := repo + ":" + commit.Hash
			allCommits[key] = commit
		}
	}

	result := make([]Commit, 0, len(allCommits))
	for _, commit := range allCommits {
		result = append(result, commit)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	if len(errors) > 0 && len(result) == 0 {
		return nil, fmt.Errorf("all repositories failed to load commits: %s", strings.Join(errors, "; "))
	}

	return result, nil
}

func getRepoDisplayName(repoPathOrURL string) string {
	if isRemoteURL(repoPathOrURL) {
		parts := strings.Split(repoPathOrURL, "/")
		if len(parts) > 0 {
			name := parts[len(parts)-1]
			if strings.HasSuffix(name, ".git") {
				name = name[:len(name)-4]
			}
			return name + " (R)"
		}
		return repoPathOrURL
	}

	parts := strings.Split(strings.ReplaceAll(repoPathOrURL, "\\", "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return repoPathOrURL
}
