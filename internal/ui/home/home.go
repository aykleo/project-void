package home

import (
	"fmt"
	"project-void/internal/config"
	"project-void/internal/ui/common"
	datepicker "project-void/internal/ui/home/date-picker"
	folderpicker "project-void/internal/ui/home/folder-picker"
	"project-void/internal/ui/home/tabs"
	"project-void/internal/ui/styles"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProceedMsg struct{}

type Model struct {
	datePicker        datepicker.Model
	folderPicker      folderpicker.Model
	tabs              tabs.Model
	commandHandler    common.CommandHandler
	selectedDate      *time.Time
	selectedFolder    string
	width             int
	height            int
	devModeSelected   bool
	isDev             bool
	needsConfirmation bool
	shouldProceed     bool
}

func InitialModel() Model {
	tabNames := []string{}
	tabContent := []string{}
	return Model{
		datePicker:     datepicker.InitialModel(),
		folderPicker:   folderpicker.InitialModel(),
		tabs:           tabs.InitialModel(tabNames, tabContent),
		commandHandler: common.NewCommandHandler("Enter a command (e.g., dev or help)..."),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.datePicker.Init(), m.folderPicker.Init(), m.commandHandler.Init())
}

func (m Model) UpdateWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	return m, nil
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

		updatedFP, fpCmd := m.folderPicker.Update(contentSizeMsg)
		m.folderPicker = updatedFP.(folderpicker.Model)

		updatedDP, dpCmd := m.datePicker.Update(contentSizeMsg)
		m.datePicker = updatedDP.(datepicker.Model)

		return m, tea.Batch(tabCmd, fpCmd, dpCmd)
	case tea.KeyMsg:
		if !m.devModeSelected {
			updatedHandler, cmd, result := m.commandHandler.Update(msg)
			m.commandHandler = updatedHandler

			if result != nil {
				if result.ShouldQuit {
					return m, tea.Quit
				}

				if result.Action == "dev" {
					m.isDev = true
					m.devModeSelected = true
					m.commandHandler.ClearMessages()
					return m, cmd
				} else if result.Action == "nodev" {
					m.isDev = false
					m.devModeSelected = true

					username := "Unknown User"
					if userConfig, err := config.LoadUserConfig(); err == nil && userConfig.Jira.Username != "" {
						username = userConfig.Jira.Username
					}

					tabNames := []string{"Jira cards", "Slack messages"}
					tabContent := []string{
						fmt.Sprintf("Jira cards for %s", username),
						"Not yet implemented",
					}
					m.tabs = tabs.InitialModel(tabNames, tabContent)

					if m.width > 0 {
						contentSizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
						updatedTabs, _ := m.tabs.Update(contentSizeMsg)
						m.tabs = updatedTabs.(tabs.Model)
					}
					m.commandHandler.ClearMessages()
					return m, cmd
				}

				if result.Action == "reset" {
					m.commandHandler.SetError("Use Esc or q to go back to welcome screen")
					return m, cmd
				}
			}

			return m, cmd
		}

		if m.needsConfirmation {
			key := strings.ToLower(msg.String())
			if key == "y" || key == "yes" {
				m.shouldProceed = true
				return m, tea.Batch(tea.Cmd(func() tea.Msg { return ProceedMsg{} }))
			} else if key == "n" || key == "no" {
				return InitialModel(), tea.Cmd(func() tea.Msg { return tea.WindowSizeMsg{Width: m.width, Height: m.height} })
			}

			updatedTabs, cmd := m.tabs.Update(msg)
			m.tabs = updatedTabs.(tabs.Model)
			return m, cmd
		}

		if m.isDev && m.selectedFolder == "" {
			updatedFolderPicker, cmd := m.folderPicker.Update(msg)
			m.folderPicker = updatedFolderPicker.(folderpicker.Model)

			if m.folderPicker.GetSelectedFolder() != "" {
				m.selectedFolder = m.folderPicker.GetSelectedFolder()

				username := "Unknown User"
				if userConfig, err := config.LoadUserConfig(); err == nil && userConfig.Jira.Username != "" {
					username = userConfig.Jira.Username
				}

				tabNames := []string{"Git commits", "Jira cards", "Slack messages"}
				tabContent := []string{
					fmt.Sprintf("Git commits for %s", m.selectedFolder),
					fmt.Sprintf("Jira cards for %s", username),
					"Not yet implemented",
				}
				m.tabs = tabs.InitialModel(tabNames, tabContent)

				if m.width > 0 {
					contentSizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
					updatedTabs, _ := m.tabs.Update(contentSizeMsg)
					m.tabs = updatedTabs.(tabs.Model)
				}
			}

			return m, cmd
		}

		if m.selectedDate == nil {
			updatedDatePicker, cmd := m.datePicker.Update(msg)
			m.datePicker = updatedDatePicker.(datepicker.Model)

			if m.datePicker.IsDateSelected() {
				selectedDate := m.datePicker.GetSelectedDate()
				m.selectedDate = &selectedDate
				m.needsConfirmation = true
			}

			return m, cmd
		} else if !m.needsConfirmation {
			m.needsConfirmation = true
			return m, nil
		} else {
			updatedTabs, cmd := m.tabs.Update(msg)
			m.tabs = updatedTabs.(tabs.Model)
			return m, cmd
		}
	default:
		if !m.devModeSelected {
			return m, nil
		}

		if m.needsConfirmation {
			return m, nil
		}

		if m.isDev && m.selectedFolder == "" {
			updatedFolderPicker, cmd := m.folderPicker.Update(msg)
			m.folderPicker = updatedFolderPicker.(folderpicker.Model)

			if m.folderPicker.GetSelectedFolder() != "" {
				m.selectedFolder = m.folderPicker.GetSelectedFolder()

				username := "Unknown User"
				if userConfig, err := config.LoadUserConfig(); err == nil && userConfig.Jira.Username != "" {
					username = userConfig.Jira.Username
				}

				tabNames := []string{"Git commits", "Jira cards", "Slack messages"}
				tabContent := []string{
					fmt.Sprintf("Repo: %s", m.selectedFolder),
					fmt.Sprintf("Jira cards for %s", username),
					"Not yet implemented",
				}
				m.tabs = tabs.InitialModel(tabNames, tabContent)

				if m.width > 0 {
					contentSizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
					updatedTabs, _ := m.tabs.Update(contentSizeMsg)
					m.tabs = updatedTabs.(tabs.Model)
				}
			}

			return m, cmd
		}

		if m.selectedDate == nil {
			updatedDatePicker, cmd := m.datePicker.Update(msg)
			m.datePicker = updatedDatePicker.(datepicker.Model)

			if m.datePicker.IsDateSelected() {
				selectedDate := m.datePicker.GetSelectedDate()
				m.selectedDate = &selectedDate
				m.needsConfirmation = true
			}

			return m, cmd
		} else if !m.needsConfirmation {
			m.needsConfirmation = true
			return m, nil
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

	var content string

	if !m.devModeSelected {
		if helpView := m.commandHandler.RenderHelp(m.width, m.height); helpView != "" {
			return helpView
		}
	}

	if !m.devModeSelected {
		helpText := "Type 'dev' to activate dev mode, 'nodev' to continue without dev mode, or 'help' for all commands"
		inputSection := m.commandHandler.RenderCommandPrompt(helpText)
		content = inputSection
		return styles.DocStyle.Width(m.width).Height(m.height).Align(lipgloss.Center, lipgloss.Center).Render(content)
	}

	if m.isDev && m.selectedFolder == "" {
		folderPicker := styles.FolderPickerStyle.Width(contentWidth).Render(m.folderPicker.View())
		content = folderPicker
		return styles.DocStyle.Width(m.width).Height(m.height).Align(lipgloss.Center, lipgloss.Center).Render(content)
	}

	var folderInfo string
	if m.isDev && m.selectedFolder != "" {
		folderMessage := fmt.Sprintf("Selected folder: %s", m.selectedFolder)
		folderInfo = styles.WelcomeStyle.Width(contentWidth).Render(folderMessage)
	}

	if m.selectedDate == nil {
		datePickerPrompt := "← →: navigate days, ↑ ↓: navigate weeks, enter/space: select date"
		prompt := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Align(lipgloss.Center).MarginBottom(2).Width(contentWidth).Render(datePickerPrompt)

		currentDate := m.datePicker.GetSelectedDate()
		initialDateMessage := fmt.Sprint(currentDate.Format("January 2, 2006"))
		dateInfo := styles.WelcomeStyle.Width(contentWidth).Render(initialDateMessage)

		if m.isDev {
			content = folderInfo + "\n" + dateInfo + "\n" + prompt
		} else {
			content = dateInfo + "\n" + prompt
		}
		return styles.DocStyle.Width(m.width).Height(m.height).Align(lipgloss.Center, lipgloss.Center).Render(content)
	}

	selectedDateMessage := fmt.Sprint(m.selectedDate.Format("January 2, 2006"))
	dateInfo := styles.WelcomeStyle.Width(contentWidth).Render(selectedDateMessage)

	contentSizeMsg := tea.WindowSizeMsg{Width: contentWidth, Height: m.height}
	updatedTabs, _ := m.tabs.Update(contentSizeMsg)
	m.tabs = updatedTabs.(tabs.Model)

	tabsView := m.tabs.View()

	var confirmationPrompt string
	if m.needsConfirmation {
		confirmationPrompt = styles.NeutralStyle.Width(contentWidth).Render("Proceed with these options? (y/n)")
	}

	if m.isDev {
		if m.needsConfirmation {
			content = dateInfo + "\n\n" + tabsView + "\n\n" + confirmationPrompt
		} else {
			content = dateInfo + "\n\n" + tabsView
		}
	} else {
		if m.needsConfirmation {
			content = dateInfo + "\n\n" + tabsView + "\n\n" + confirmationPrompt
		} else {
			content = dateInfo + "\n\n" + tabsView
		}
	}

	return styles.DocStyle.Width(m.width).Height(m.height).Align(lipgloss.Center, lipgloss.Center).Render(content)
}

func (m Model) ShouldProceed() bool {
	return m.shouldProceed
}

func (m Model) GetSelectedFolder() string {
	return m.selectedFolder
}

func (m Model) GetSelectedDate() *time.Time {
	return m.selectedDate
}

func (m Model) IsDevMode() bool {
	return m.isDev
}
