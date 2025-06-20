package statistics

import (
	"fmt"
	"project-void/internal/jira"
	commitstable "project-void/internal/ui/statistics/commits-table"
	jiratable "project-void/internal/ui/statistics/jira-table"
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

func loadCommitsFromMultipleReposCmd(repoSources []string, since time.Time) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if len(repoSources) == 0 {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommitsFromMultipleRepos(repoSources, since)
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

func loadCommitsByAuthorsFromMultipleReposCmd(repoSources []string, since time.Time, authorNames []string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if len(repoSources) == 0 {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommitsByAuthorsFromMultipleRepos(repoSources, since, authorNames)
		if err != nil {
			return LoadErrorMsg{Error: err.Error()}
		}

		return LoadedMsg{CommitsTable: commitsTable}
	})
}

func loadCommitsByBranchesFromMultipleReposCmd(repoSources []string, since time.Time, branchNames []string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if len(repoSources) == 0 {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommitsByBranchesFromMultipleRepos(repoSources, since, branchNames)
		if err != nil {
			return LoadErrorMsg{Error: err.Error()}
		}

		return LoadedMsg{CommitsTable: commitsTable}
	})
}

func loadCommitsByAuthorsAndBranchesFromMultipleReposCmd(repoSources []string, since time.Time, authorNames []string, branchNames []string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if len(repoSources) == 0 {
			emptyTable := commitstable.InitialModel()
			emptyTable.StartLoading()
			return LoadedMsg{CommitsTable: emptyTable}
		}

		var commitsTable commitstable.Model = commitstable.InitialModel()
		commitsTable.StartLoading()
		err := commitsTable.LoadCommitsByAuthorsAndBranchesFromMultipleRepos(repoSources, since, authorNames, branchNames)
		if err != nil {
			return LoadErrorMsg{Error: err.Error()}
		}

		return LoadedMsg{CommitsTable: commitsTable}
	})
}

func loadJiraCmd(jiraSource string, since time.Time) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if jiraSource == "" {
			emptyTable := jiratable.InitialModel()
			emptyTable.StartLoading()
			return JiraLoadedMsg{JiraTable: emptyTable}
		}

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
