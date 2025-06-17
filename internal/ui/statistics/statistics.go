package statistics

import (
	"project-void/internal/config"
	"project-void/internal/ui/common"
	commitstable "project-void/internal/ui/statistics/commits-table"
	jiratable "project-void/internal/ui/statistics/jira-table"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	commitsTable       commitstable.Model
	jiraTable          jiratable.Model
	commandHandler     common.StatisticsCommandHandler
	selectedFolder     string
	selectedRepoSource string
	selectedJiraSource string
	selectedDate       time.Time
	hasGit             bool
	hasJira            bool
	width              int
	height             int
	loaded             bool
	loadError          string
	focusedTable       int
	command            string
	submitted          bool
	authorFilter       []string
	branchFilter       []string
	commitsSpinner     spinner.Model
	jiraSpinner        spinner.Model
	commitsLoading     bool
	jiraLoading        bool
	noConfigMessage    string
}

func InitialModel(selectedFolder string, selectedDate time.Time, hasGit, hasJira bool) Model {
	commitsTable := commitstable.InitialModel()
	jiraTable := jiratable.InitialModel()

	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	commitsSpinner := spinner.New()
	commitsSpinner.Style = spinnerStyle
	commitsSpinner.Spinner = spinner.Dot

	jiraSpinner := spinner.New()
	jiraSpinner.Style = spinnerStyle
	jiraSpinner.Spinner = spinner.Dot

	repoSource := selectedFolder
	if selectedFolder == "" {
		if gitConfig, err := config.LoadUserConfig(); err == nil && gitConfig.Git.RepoURL != "" {
			repoSource = gitConfig.Git.RepoURL
		}
	}

	jiraSource := ""
	if jiraConfig, err := config.LoadUserConfig(); err == nil && jiraConfig.Jira.BaseURL != "" {
		jiraSource = jiraConfig.Jira.BaseURL
	}

	actualHasGit := hasGit
	actualHasJira := hasJira

	if actualHasGit {
		commitsTable.StartLoading()
	}
	if actualHasJira {
		jiraTable.StartLoading()
	}
	noConfigMessage := "Please use 'void help', 'void help git', or 'void help jira' to configure your repositories"

	if actualHasGit && actualHasJira {
		commitsTable.Focus()
		commitsTable.SetFocusedStyle()
		jiraTable.Blur()
		jiraTable.SetBlurredStyle()
	} else if !actualHasGit && actualHasJira {
		jiraTable.Focus()
		jiraTable.SetFocusedStyle()
	} else if actualHasGit && !actualHasJira {
		commitsTable.Focus()
		commitsTable.SetFocusedStyle()
	}

	commitsTable.SetSpinner(&commitsSpinner)
	jiraTable.SetSpinner(&jiraSpinner)

	return Model{
		commitsTable:       commitsTable,
		jiraTable:          jiraTable,
		commandHandler:     common.NewStatisticsCommandHandler("Enter a command (e.g., git repo <url>, git a <author>, void help)...", repoSource, actualHasGit, actualHasJira),
		selectedFolder:     selectedFolder,
		selectedRepoSource: repoSource,
		selectedJiraSource: jiraSource,
		selectedDate:       selectedDate,
		hasGit:             actualHasGit,
		hasJira:            actualHasJira,
		noConfigMessage:    noConfigMessage,
		focusedTable:       0,
		commitsSpinner:     commitsSpinner,
		jiraSpinner:        jiraSpinner,
		commitsLoading:     actualHasGit,
		jiraLoading:        actualHasJira,
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.commandHandler.Init())

	if m.hasGit {
		cmds = append(cmds, m.commitsTable.Init(), m.commitsSpinner.Tick, loadCommitsCmd(m.selectedRepoSource, m.selectedDate))
	}

	if m.hasJira {
		cmds = append(cmds, m.jiraTable.Init(), m.jiraSpinner.Tick, loadJiraCmd(m.selectedJiraSource, m.selectedDate))
	}

	return tea.Batch(cmds...)
}

type LoadedMsg struct {
	CommitsTable commitstable.Model
}

type LoadErrorMsg struct {
	Error string
}

type JiraLoadedMsg struct {
	JiraTable jiratable.Model
}

type JiraLoadErrorMsg struct {
	Error string
}

func (m Model) GetCommand() string {
	return m.command
}

func (m Model) HasCommand() bool {
	return m.submitted
}

func (m *Model) ResetCommand() {
	m.submitted = false
	m.command = ""
	m.commandHandler.ClearMessages()
}
