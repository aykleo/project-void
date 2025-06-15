package welcome

import (
	"project-void/internal/ui/common"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	commandHandler common.CommandHandler
	width          int
	height         int
	command        string
	submitted      bool
	selectedDate   *time.Time
}

func InitialModel() Model {
	return Model{
		commandHandler: common.NewCommandHandler("Enter a command (e.g., git repo <url>, void help)..."),
		submitted:      false,
		selectedDate:   nil,
	}
}

func (m Model) Init() tea.Cmd {
	return m.commandHandler.Init()
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

func (m Model) GetSelectedDate() *time.Time {
	return m.selectedDate
}
