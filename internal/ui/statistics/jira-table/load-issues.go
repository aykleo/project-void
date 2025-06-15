package jiratable

import (
	"fmt"
	"project-void/internal/jira"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

func (m *Model) LoadIssues(jiraClient *jira.JiraClient, since time.Time, config *jira.Config) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	m.progress.SetPercent(0.2)
	issues, err := jiraClient.GetIssuesSince(since, config)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load issues: %w", err)
	}

	rows := make([]table.Row, len(issues))
	for i, issue := range issues {
		if len(issues) > 0 {
			progress := 0.2 + (0.7 * float64(i) / float64(len(issues)))
			m.progress.SetPercent(progress)
		}

		action := issue.UserAction
		if len(action) > 11 {
			action = action[:8] + "..."
		}

		status := issue.Status
		if len(status) > 14 {
			status = status[:11] + "..."
		}

		actionDate := issue.ActionDate.Format("2006-01-02")

		summary := strings.ReplaceAll(issue.Summary, "\n", " ")
		summary = strings.TrimSpace(summary)
		maxSummaryLength := 80
		if len(summary) > maxSummaryLength {
			summary = summary[:maxSummaryLength-3] + "..."
		}

		rows[i] = table.Row{
			issue.Key,
			status,
			action,
			actionDate,
			summary,
		}
	}

	m.progress.SetPercent(0.9)
	m.table.SetRows(rows)
	m.progress.SetPercent(1.0)
	m.loadingState = LoadingComplete
	return nil
}
