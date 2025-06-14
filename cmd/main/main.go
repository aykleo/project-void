package main

import (
	"fmt"
	"os"

	"project-void/internal/ui/home"
	"project-void/internal/ui/statistics"

	tea "github.com/charmbracelet/bubbletea"
)

type AppState int

const (
	HomeState AppState = iota
	StatisticsState
)

type MainModel struct {
	state      AppState
	homeModel  home.Model
	statsModel statistics.Model
	width      int
	height     int
}

func InitialMainModel() MainModel {
	return MainModel{
		state:     HomeState,
		homeModel: home.InitialModel(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.homeModel.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		}
	case home.ProceedMsg:
		if m.homeModel.ShouldProceed() {
			m.state = StatisticsState
			selectedFolder := m.homeModel.GetSelectedFolder()
			selectedDate := m.homeModel.GetSelectedDate()
			isDev := m.homeModel.IsDevMode()

			if selectedDate != nil {
				m.statsModel = statistics.InitialModel(selectedFolder, *selectedDate, isDev)
				var cmds []tea.Cmd
				cmds = append(cmds, m.statsModel.Init())
				if m.width > 0 && m.height > 0 {
					windowSizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
					updatedStats, cmd := m.statsModel.Update(windowSizeMsg)
					m.statsModel = updatedStats.(statistics.Model)
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
				}
				return m, tea.Batch(cmds...)
			}
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		switch m.state {
		case HomeState:
			updatedHome, cmd := m.homeModel.Update(msg)
			m.homeModel = updatedHome.(home.Model)
			return m, cmd
		case StatisticsState:
			updatedStats, cmd := m.statsModel.Update(msg)
			m.statsModel = updatedStats.(statistics.Model)
			return m, cmd
		}
	}

	switch m.state {
	case HomeState:
		updatedHome, cmd := m.homeModel.Update(msg)
		m.homeModel = updatedHome.(home.Model)
		return m, cmd
	case StatisticsState:
		updatedStats, cmd := m.statsModel.Update(msg)
		m.statsModel = updatedStats.(statistics.Model)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) View() string {
	switch m.state {
	case HomeState:
		return m.homeModel.View()
	case StatisticsState:
		return m.statsModel.View()
	}
	return ""
}

func main() {
	m := InitialMainModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
