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
			shortBranch,
			commit.Author,
			dateStr,
			message,
		}
	}

	m.table.SetRows(rows)
	return nil
}
