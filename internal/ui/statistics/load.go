package statistics

import (
	"fmt"
	"project-void/internal/jira"
	commitstable "project-void/internal/ui/statistics/commits-table"
	jiratable "project-void/internal/ui/statistics/jira-table"
	slacktable "project-void/internal/ui/statistics/slack-table"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func loadCommitsCmd(repoSource string, since time.Time) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if repoSource == "" {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommits(repoSource, since)
		if err != nil {
			return LoadErrorMsg{Error: err.Error()}
		}

		return LoadedMsg{CommitsTable: commitsTable}
	})
}

func loadCommitsByAuthorsCmd(repoSource string, since time.Time, authorNames []string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if repoSource == "" {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommitsByAuthors(repoSource, since, authorNames)
		if err != nil {
			return LoadErrorMsg{Error: err.Error()}
		}

		return LoadedMsg{CommitsTable: commitsTable}
	})
}

func loadCommitsByBranchesCmd(repoSource string, since time.Time, branchNames []string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if repoSource == "" {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommitsByBranches(repoSource, since, branchNames)
		if err != nil {
			return LoadErrorMsg{Error: err.Error()}
		}

		return LoadedMsg{CommitsTable: commitsTable}
	})
}

func loadCommitsByAuthorsAndBranchesCmd(repoSource string, since time.Time, authorNames []string, branchNames []string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if repoSource == "" {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommitsByAuthorsAndBranches(repoSource, since, authorNames, branchNames)
		if err != nil {
			return LoadErrorMsg{Error: err.Error()}
		}

		return LoadedMsg{CommitsTable: commitsTable}
	})
}

func loadJiraCmd(since time.Time) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		var jiraTable jiratable.Model = jiratable.InitialModel()
		jiraTable.StartLoading()

		config, err := jira.LoadConfig()
		if err != nil {
			return JiraLoadErrorMsg{Error: fmt.Sprintf("Failed to load JIRA config: %v", err)}
		}

		client := jira.NewClientFromConfig(config)

		err = jiraTable.LoadIssues(client, since, config)
		if err != nil {
			return JiraLoadErrorMsg{Error: fmt.Sprintf("Failed to load JIRA issues: %v", err)}
		}
		return JiraLoadedMsg{JiraTable: jiraTable}
	})
}

func loadSlackCmd() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		var slackTable slacktable.Model = slacktable.InitialModel()
		slackTable.StartLoading()
		slackTable.SetPlaceholder()
		return SlackLoadedMsg{SlackTable: slackTable}
	})
}
