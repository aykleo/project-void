package datepicker

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

type Model struct {
	DatePicker datepicker.Model
}

func InitialModel() Model {
	dp := datepicker.New(time.Now())

	return Model{
		DatePicker: dp,
	}
}

func (m Model) Init() tea.Cmd {
	return m.DatePicker.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			m.DatePicker.SelectDate()
			return m, nil
		}
	}

	dp, cmd := m.DatePicker.Update(msg)
	m.DatePicker = dp

	return m, cmd
}

func (m Model) View() string {
	return m.DatePicker.View()
}

func (m Model) GetSelectedDate() time.Time {
	return m.DatePicker.Time
}

func (m Model) IsDateSelected() bool {
	return m.DatePicker.Selected
}
