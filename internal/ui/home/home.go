package home

import (
	"project-void/internal/ui/home/tabs"
	"project-void/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tabs   tabs.Model
	width  int
	height int
}

func InitialModel() Model {
	tabNames := []string{"Git commits", "Jira cards", "Slack messages"}
	tabContent := []string{
		"This is the content of the Git commits tab",
		"This is the content of the Jira cards tab",
		"This is the content of the Slack messages tab",
	}
	return Model{
		tabs: tabs.InitialModel(tabNames, tabContent),
	}
}

func (m Model) Init() tea.Cmd {
	return m.tabs.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		horizontalPadding := 4
		contentWidth := msg.Width - (horizontalPadding * 2)

		contentSizeMsg := tea.WindowSizeMsg{Width: contentWidth, Height: msg.Height}
		updatedTabs, cmd := m.tabs.Update(contentSizeMsg)
		m.tabs = updatedTabs.(tabs.Model)
		return m, cmd
	default:

		updatedTabs, cmd := m.tabs.Update(msg)
		m.tabs = updatedTabs.(tabs.Model)
		return m, cmd
	}
}

func (m Model) View() string {
	horizontalPadding := 4
	contentWidth := m.width - (horizontalPadding * 2)

	welcomeMessage := "Welcome to Project Void"
	welcome := styles.WelcomeStyle.Width(contentWidth).Render(welcomeMessage)

	quitMessage := "Q or Esc to quit"
	quit := styles.QuitStyle.Width(contentWidth).Render(quitMessage)

	tabsView := m.tabs.View()

	content := welcome + "\n" + quit + "\n\n" + tabsView

	return styles.DocStyle.Width(m.width).Render(content)
}
