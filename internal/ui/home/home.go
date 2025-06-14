package home

import (
	"fmt"
	datepicker "project-void/internal/ui/home/date-picker"
	"project-void/internal/ui/home/tabs"
	"project-void/internal/ui/styles"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	datePicker   datepicker.Model
	tabs         tabs.Model
	selectedDate *time.Time
	width        int
	height       int
}

func InitialModel() Model {
	tabNames := []string{"Git commits", "Jira cards", "Slack messages"}
	tabContent := []string{
		"This is the content of the Git commits tab",
		"This is the content of the Jira cards tab",
		"This is the content of the Slack messages tab",
	}
	return Model{
		datePicker: datepicker.InitialModel(),
		tabs:       tabs.InitialModel(tabNames, tabContent),
	}
}

func (m Model) Init() tea.Cmd {
	return m.datePicker.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		horizontalPadding := 4
		contentWidth := msg.Width - (horizontalPadding * 2)

		contentSizeMsg := tea.WindowSizeMsg{Width: contentWidth, Height: msg.Height}

		updatedTabs, tabCmd := m.tabs.Update(contentSizeMsg)
		m.tabs = updatedTabs.(tabs.Model)

		if m.selectedDate != nil {
			return m, tabCmd
		}
		return m, nil
	default:

		if m.selectedDate == nil {
			updatedDatePicker, cmd := m.datePicker.Update(msg)
			m.datePicker = updatedDatePicker.(datepicker.Model)

			if m.datePicker.IsDateSelected() {
				selectedDate := m.datePicker.GetSelectedDate()
				m.selectedDate = &selectedDate
			}

			return m, cmd
		} else {
			updatedTabs, cmd := m.tabs.Update(msg)
			m.tabs = updatedTabs.(tabs.Model)
			return m, cmd
		}
	}
}

func (m Model) View() string {
	horizontalPadding := 4
	contentWidth := m.width - (horizontalPadding * 2)

	welcomeMessage := "Welcome to Project Void"
	welcome := styles.WelcomeStyle.Width(contentWidth).Render(welcomeMessage)

	quitMessage := "Q or Esc to quit"
	quit := styles.QuitStyle.Width(contentWidth).Render(quitMessage)

	var content string

	if m.selectedDate == nil {
		datePickerPrompt := "Please select a date to continue: ← → for days, ↑ ↓ for months, Enter/Space to select"
		prompt := styles.NeutralStyle.Width(contentWidth).Render(datePickerPrompt)

		currentDate := m.datePicker.GetSelectedDate()
		initialDateMessage := fmt.Sprint(currentDate.Format("January 2, 2006"))
		dateInfo := styles.WelcomeStyle.Width(contentWidth).Render(initialDateMessage)

		content = welcome + "\n" + quit + "\n\n" + prompt + "\n" + dateInfo
	} else {
		selectedDateMessage := fmt.Sprint(m.selectedDate.Format("January 2, 2006"))
		dateInfo := styles.WelcomeStyle.Width(contentWidth).Render(selectedDateMessage)

		tabsView := m.tabs.View()

		content = welcome + "\n" + quit + "\n\n" + dateInfo + "\n\n" + tabsView
	}

	return styles.DocStyle.Width(m.width).Render(content)
}
