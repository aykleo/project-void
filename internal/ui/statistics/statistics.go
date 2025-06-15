package statistics

import (
	"project-void/internal/config"
	"project-void/internal/ui/common"
	commitstable "project-void/internal/ui/statistics/commits-table"
	jiratable "project-void/internal/ui/statistics/jira-table"
	slacktable "project-void/internal/ui/statistics/slack-table"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	commitsTable       commitstable.Model
	jiraTable          jiratable.Model
	slackTable         slacktable.Model
	commandHandler     common.StatisticsCommandHandler
	selectedFolder     string
	selectedRepoSource string
	selectedDate       time.Time
	isDev              bool
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
	slackSpinner       spinner.Model
	commitsLoading     bool
	jiraLoading        bool
	slackLoading       bool
}

func InitialModel(selectedFolder string, selectedDate time.Time, isDev bool) Model {
	commitsTable := commitstable.InitialModel()
	jiraTable := jiratable.InitialModel()
	slackTable := slacktable.InitialModel()

	commitsTable.StartLoading()
	jiraTable.StartLoading()
	slackTable.StartLoading()

	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	commitsSpinner := spinner.New()
	commitsSpinner.Style = spinnerStyle
	commitsSpinner.Spinner = spinner.Dot

	jiraSpinner := spinner.New()
	jiraSpinner.Style = spinnerStyle
	jiraSpinner.Spinner = spinner.Dot

	slackSpinner := spinner.New()
	slackSpinner.Style = spinnerStyle
	slackSpinner.Spinner = spinner.Dot

	repoSource := selectedFolder
	if selectedFolder == "" {
		if gitConfig, err := config.LoadUserConfig(); err == nil && gitConfig.Git.RepoURL != "" {
			repoSource = gitConfig.Git.RepoURL
		}
	}

	actualIsDev := repoSource != ""

	if actualIsDev {
		commitsTable.Focus()
		commitsTable.SetFocusedStyle()
		jiraTable.Blur()
		jiraTable.SetBlurredStyle()
		slackTable.Blur()
		slackTable.SetBlurredStyle()
	} else {
		jiraTable.Focus()
		jiraTable.SetFocusedStyle()
		slackTable.Blur()
		slackTable.SetBlurredStyle()
	}

	commitsTable.SetSpinner(&commitsSpinner)
	jiraTable.SetSpinner(&jiraSpinner)
	slackTable.SetSpinner(&slackSpinner)

	return Model{
		commitsTable:       commitsTable,
		jiraTable:          jiraTable,
		slackTable:         slackTable,
		commandHandler:     common.NewStatisticsCommandHandler("Enter a command (e.g., git repo <url>, git a <author>, void help)...", repoSource, actualIsDev),
		selectedFolder:     selectedFolder,
		selectedRepoSource: repoSource,
		selectedDate:       selectedDate,
		isDev:              actualIsDev,
		focusedTable:       0,
		commitsSpinner:     commitsSpinner,
		jiraSpinner:        jiraSpinner,
		slackSpinner:       slackSpinner,
		commitsLoading:     true,
		jiraLoading:        true,
		slackLoading:       true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.commitsTable.Init(),
		m.jiraTable.Init(),
		m.slackTable.Init(),
		m.commandHandler.Init(),
		m.commitsSpinner.Tick,
		m.jiraSpinner.Tick,
		m.slackSpinner.Tick,
		loadCommitsCmd(m.selectedRepoSource, m.selectedDate),
		loadJiraCmd(m.selectedDate),
		loadSlackCmd(),
	)
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

type SlackLoadedMsg struct {
	SlackTable slacktable.Model
}

type SlackLoadErrorMsg struct {
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
