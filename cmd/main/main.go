package main

import (
	"fmt"
	"os"

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
