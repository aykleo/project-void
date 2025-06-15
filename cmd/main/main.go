package main

import (
	"fmt"
	"os"
	"time"

	"project-void/internal/git"
	"project-void/internal/ui/statistics"
	"project-void/internal/ui/welcome"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type AppState int

const (
	WelcomeState AppState = iota
	StatisticsState
)

type MainModel struct {
	state        AppState
	welcomeModel welcome.Model
	statsModel   statistics.Model
	width        int
	height       int
}

func InitialMainModel() MainModel {
	return MainModel{
		state:        WelcomeState,
		welcomeModel: welcome.InitialModel(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.welcomeModel.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.state == WelcomeState {
				return m, tea.Quit
			}
			m.state = WelcomeState
			if m.width > 0 && m.height > 0 {
				windowSizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
				updatedWelcome, cmd := m.welcomeModel.Update(windowSizeMsg)
				m.welcomeModel = updatedWelcome.(welcome.Model)
				return m, cmd
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		switch m.state {
		case WelcomeState:
			updatedWelcome, cmd := m.welcomeModel.Update(msg)
			m.welcomeModel = updatedWelcome.(welcome.Model)
			return m, cmd
		case StatisticsState:
			updatedStats, cmd := m.statsModel.Update(msg)
			m.statsModel = updatedStats.(statistics.Model)
			return m, cmd
		}
	}

	switch m.state {
	case WelcomeState:
		updatedWelcome, cmd := m.welcomeModel.Update(msg)
		m.welcomeModel = updatedWelcome.(welcome.Model)

		if m.welcomeModel.HasCommand() {
			command := m.welcomeModel.GetCommand()
			m.welcomeModel.ResetCommand()

			switch command {
			case "start":
				m.state = StatisticsState

				var selectedDate time.Time
				if welcomeDate := m.welcomeModel.GetSelectedDate(); welcomeDate != nil {
					selectedDate = *welcomeDate
				} else {
					selectedDate = time.Now()
				}

				repoSource := git.GetConfiguredRepoSource()

				m.statsModel = statistics.InitialModel(repoSource, selectedDate, false)

				initCmd := m.statsModel.Init()

				var cmds []tea.Cmd
				cmds = append(cmds, tea.EnterAltScreen, initCmd)
				if m.width > 0 && m.height > 0 {
					windowSizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
					updatedStats, sizeCmd := m.statsModel.Update(windowSizeMsg)
					m.statsModel = updatedStats.(statistics.Model)
					if sizeCmd != nil {
						cmds = append(cmds, sizeCmd)
					}
				}

				return m, tea.Batch(cmds...)
			case "reset":
				return m, nil
			case "quit":
				return m, tea.Quit
			default:
				return m, nil
			}
		}

		return m, cmd
	case StatisticsState:
		updatedStats, cmd := m.statsModel.Update(msg)
		m.statsModel = updatedStats.(statistics.Model)

		if m.statsModel.HasCommand() {
			command := m.statsModel.GetCommand()
			m.statsModel.ResetCommand()

			switch command {
			case "start":
				return m, nil
			case "reset":
				m.state = WelcomeState
				if m.width > 0 && m.height > 0 {
					windowSizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
					updatedWelcome, cmd := m.welcomeModel.Update(windowSizeMsg)
					m.welcomeModel = updatedWelcome.(welcome.Model)
					return m, cmd
				}
				return m, nil
			default:
				return m, nil
			}
		}

		return m, cmd
	}

	return m, nil
}

func (m MainModel) View() string {
	switch m.state {
	case WelcomeState:
		return m.welcomeModel.View()
	case StatisticsState:
		return m.statsModel.View()
	}
	return ""
}

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	m := InitialMainModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
