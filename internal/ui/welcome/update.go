package welcome

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	updatedHandler, cmd, result := m.commandHandler.Update(msg)
	m.commandHandler = updatedHandler

	if result != nil {
		if result.ShouldQuit {
			return m, tea.Quit
		}

		if result.ShouldNavigate {
			m.command = result.Action
			m.submitted = true
			return m, cmd
		}

		if result.Success && result.Action == "void_set_date" {
			if dateData, ok := result.Data["date"].(time.Time); ok {
				m.selectedDate = &dateData
			}
			return m, cmd
		}

		if result.Success {
			return m, cmd
		}
	}

	return m, cmd
}
