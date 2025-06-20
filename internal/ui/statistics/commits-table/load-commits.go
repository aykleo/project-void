package commitstable

import (
	"fmt"
	"project-void/internal/git"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

func (m *Model) LoadCommits(repoPath string, since time.Time) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	commits, err := git.GetCommitsSince(repoPath, since)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load commits: %w", err)
	}

	rows := make([]table.Row, len(commits))
	for i, commit := range commits {
		repoDisplay := commit.RepoName
		if len(repoDisplay) > 13 {
			repoDisplay = repoDisplay[:10] + "..."
		}

		shortBranch := commit.Branch
		if len(shortBranch) > 10 {
			shortBranch = shortBranch[:10]
		}

		dateStr := commit.Timestamp.Format("2006-01-02")

		message := strings.ReplaceAll(commit.Message, "\n", " ")
		message = strings.TrimSpace(message)
		if len(message) > 80 {
			message = message[:77] + "..."
		}

		rows[i] = table.Row{
			repoDisplay,
			shortBranch,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}

func (m *Model) LoadCommitsByAuthors(repoPath string, since time.Time, authorNames []string) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	commits, err := git.GetCommitsSinceByAuthors(repoPath, since, authorNames)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load commits by authors: %w", err)
	}

	rows := make([]table.Row, len(commits))
	for i, commit := range commits {
		repoDisplay := commit.RepoName
		if len(repoDisplay) > 13 {
			repoDisplay = repoDisplay[:10] + "..."
		}

		shortBranch := commit.Branch
		if len(shortBranch) > 10 {
			shortBranch = shortBranch[:10]
		}

		dateStr := commit.Timestamp.Format("2006-01-02")

		message := strings.ReplaceAll(commit.Message, "\n", " ")
		message = strings.TrimSpace(message)
		if len(message) > 80 {
			message = message[:77] + "..."
		}

		rows[i] = table.Row{
			repoDisplay,
			shortBranch,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}

func (m *Model) LoadCommitsByBranches(repoPath string, since time.Time, branchNames []string) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	commits, err := git.GetCommitsSinceByBranches(repoPath, since, branchNames)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load commits by branches: %w", err)
	}

	rows := make([]table.Row, len(commits))
	for i, commit := range commits {
		repoDisplay := commit.RepoName
		if len(repoDisplay) > 13 {
			repoDisplay = repoDisplay[:10] + "..."
		}

		shortBranch := commit.Branch
		if len(shortBranch) > 10 {
			shortBranch = shortBranch[:10]
		}

		dateStr := commit.Timestamp.Format("2006-01-02")

		message := strings.ReplaceAll(commit.Message, "\n", " ")
		message = strings.TrimSpace(message)
		if len(message) > 80 {
			message = message[:77] + "..."
		}

		rows[i] = table.Row{
			repoDisplay,
			shortBranch,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}

func (m *Model) LoadCommitsByAuthorsAndBranches(repoPath string, since time.Time, authorNames []string, branchNames []string) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	commits, err := git.GetCommitsSinceByAuthorsAndBranches(repoPath, since, authorNames, branchNames)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load commits by authors and branches: %w", err)
	}

	rows := make([]table.Row, len(commits))
	for i, commit := range commits {
		repoDisplay := commit.RepoName
		if len(repoDisplay) > 13 {
			repoDisplay = repoDisplay[:10] + "..."
		}

		shortBranch := commit.Branch
		if len(shortBranch) > 10 {
			shortBranch = shortBranch[:10]
		}

		dateStr := commit.Timestamp.Format("2006-01-02")

		message := strings.ReplaceAll(commit.Message, "\n", " ")
		message = strings.TrimSpace(message)
		if len(message) > 80 {
			message = message[:77] + "..."
		}

		rows[i] = table.Row{
			repoDisplay,
			shortBranch,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}

func (m *Model) LoadCommitsFromMultipleRepos(repoPaths []string, since time.Time) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	commits, err := git.GetCommitsSinceFromMultipleRepos(repoPaths, since)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load commits from multiple repos: %w", err)
	}

	rows := make([]table.Row, len(commits))
	for i, commit := range commits {
		repoDisplay := commit.RepoName
		if len(repoDisplay) > 13 {
			repoDisplay = repoDisplay[:10] + "..."
		}

		shortBranch := commit.Branch
		if len(shortBranch) > 10 {
			shortBranch = shortBranch[:10]
		}

		dateStr := commit.Timestamp.Format("2006-01-02")

		message := strings.ReplaceAll(commit.Message, "\n", " ")
		message = strings.TrimSpace(message)
		if len(message) > 80 {
			message = message[:77] + "..."
		}

		rows[i] = table.Row{
			repoDisplay,
			shortBranch,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}

func (m *Model) LoadCommitsByAuthorsFromMultipleRepos(repoPaths []string, since time.Time, authorNames []string) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	commits, err := git.GetCommitsSinceByAuthorsFromMultipleRepos(repoPaths, since, authorNames)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load commits by authors from multiple repos: %w", err)
	}

	rows := make([]table.Row, len(commits))
	for i, commit := range commits {
		repoDisplay := commit.RepoName
		if len(repoDisplay) > 13 {
			repoDisplay = repoDisplay[:10] + "..."
		}

		shortBranch := commit.Branch
		if len(shortBranch) > 10 {
			shortBranch = shortBranch[:10]
		}

		dateStr := commit.Timestamp.Format("2006-01-02")

		message := strings.ReplaceAll(commit.Message, "\n", " ")
		message = strings.TrimSpace(message)
		if len(message) > 80 {
			message = message[:77] + "..."
		}

		rows[i] = table.Row{
			repoDisplay,
			shortBranch,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}

func (m *Model) LoadCommitsByBranchesFromMultipleRepos(repoPaths []string, since time.Time, branchNames []string) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	commits, err := git.GetCommitsSinceByBranchesFromMultipleRepos(repoPaths, since, branchNames)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load commits by branches from multiple repos: %w", err)
	}

	rows := make([]table.Row, len(commits))
	for i, commit := range commits {
		repoDisplay := commit.RepoName
		if len(repoDisplay) > 13 {
			repoDisplay = repoDisplay[:10] + "..."
		}

		shortBranch := commit.Branch
		if len(shortBranch) > 10 {
			shortBranch = shortBranch[:10]
		}

		dateStr := commit.Timestamp.Format("2006-01-02")

		message := strings.ReplaceAll(commit.Message, "\n", " ")
		message = strings.TrimSpace(message)
		if len(message) > 80 {
			message = message[:77] + "..."
		}

		rows[i] = table.Row{
			repoDisplay,
			shortBranch,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}

func (m *Model) LoadCommitsByAuthorsAndBranchesFromMultipleRepos(repoPaths []string, since time.Time, authorNames []string, branchNames []string) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	commits, err := git.GetCommitsSinceByAuthorsAndBranchesFromMultipleRepos(repoPaths, since, authorNames, branchNames)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load commits by authors and branches from multiple repos: %w", err)
	}

	rows := make([]table.Row, len(commits))
	for i, commit := range commits {
		repoDisplay := commit.RepoName
		if len(repoDisplay) > 13 {
			repoDisplay = repoDisplay[:10] + "..."
		}

		shortBranch := commit.Branch
		if len(shortBranch) > 10 {
			shortBranch = shortBranch[:10]
		}

		dateStr := commit.Timestamp.Format("2006-01-02")

		message := strings.ReplaceAll(commit.Message, "\n", " ")
		message = strings.TrimSpace(message)
		if len(message) > 80 {
			message = message[:77] + "..."
		}

		rows[i] = table.Row{
			repoDisplay,
			shortBranch,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}
