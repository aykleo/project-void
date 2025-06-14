package home

import (
	"fmt"
	"os"
	"project-void/internal/commands"
	datepicker "project-void/internal/ui/home/date-picker"
	folderpicker "project-void/internal/ui/home/folder-picker"
	"project-void/internal/ui/home/tabs"
	"project-void/internal/ui/styles"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProceedMsg struct{}

type Model struct {
	datePicker        datepicker.Model
	folderPicker      folderpicker.Model
	tabs              tabs.Model
	textInput         textinput.Model
	selectedDate      *time.Time
	selectedFolder    string
	width             int
	height            int
	devModeSelected   bool
	isDev             bool
	needsConfirmation bool
	shouldProceed     bool
	commandError      string
	showingHelp       bool
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter a command (e.g., dev or help)..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	tabNames := []string{}
	tabContent := []string{}
	return Model{
		datePicker:   datepicker.InitialModel(),
		folderPicker: folderpicker.InitialModel(),
		tabs:         tabs.InitialModel(tabNames, tabContent),
		textInput:    ti,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.datePicker.Init(), m.folderPicker.Init(), textinput.Blink)
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

		if msg.Width > 60 {
			m.textInput.Width = 50
		} else {
			m.textInput.Width = msg.Width - 10
		}

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
		if m.showingHelp {
			m.showingHelp = false
			m.commandError = ""
			return m, nil
		}

		if !m.devModeSelected {
			switch msg.Type {
			case tea.KeyEnter:
				inputValue := m.textInput.Value()
				validatedCmd, err := commands.ValidateCommand(inputValue)
				if err != nil {
					m.commandError = err.Error()
					m.textInput.SetValue("")
					return m, nil
				}

				if validatedCmd.Action == "help" {
					m.showingHelp = true
					m.commandError = ""
					m.textInput.SetValue("")
					return m, nil
				}

				if validatedCmd.Action == "reset" {
					m.commandError = "Use Esc or q to go back to welcome screen"
					m.textInput.SetValue("")
					return m, nil
				}

				if validatedCmd.Action == "quit" {
					return m, tea.Quit
				}

				if validatedCmd.Action == "dev" {
					m.isDev = true
					m.devModeSelected = true
					m.commandError = ""
					return m, nil
				} else if validatedCmd.Action == "nodev" {
					m.isDev = false
					m.devModeSelected = true
					username := os.Getenv("JIRA_USERNAME")
					tabNames := []string{"Jira cards", "Slack messages"}
					tabContent := []string{
						fmt.Sprintf("Jira cards for %s", username),
						"This is the content of the Slack messages tab",
					}
					m.tabs = tabs.InitialModel(tabNames, tabContent)

					if m.width > 0 {
						contentSizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
						updatedTabs, _ := m.tabs.Update(contentSizeMsg)
						m.tabs = updatedTabs.(tabs.Model)
					}
					m.commandError = ""
					return m, nil
				}

				m.commandError = fmt.Sprintf("Unknown command action: %s", validatedCmd.Action)
				m.textInput.SetValue("")
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			}

			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
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
				username := os.Getenv("JIRA_USERNAME")
				tabNames := []string{"Git commits", "Jira cards", "Slack messages"}
				tabContent := []string{
					fmt.Sprintf("Git commits for %s", m.selectedFolder),
					fmt.Sprintf("Jira cards for %s", username),
					"This is the content of the Slack messages tab",
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
				username := os.Getenv("JIRA_USERNAME")
				tabNames := []string{"Git commits", "Jira cards", "Slack messages"}
				tabContent := []string{
					fmt.Sprintf("Repo: %s", m.selectedFolder),
					fmt.Sprintf("Jira cards for %s", username),
					"This is the content of the Slack messages tab",
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

	if m.showingHelp {
		helpText := commands.GetHelpText()
		helpContent := fmt.Sprintf("%s\n\nPress any key to return to command prompt", helpText)

		centerStyle := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1, 2)

		return centerStyle.Render(styles.NeutralStyle.Render(helpContent))
	}

	if !m.devModeSelected {
		var inputSection string
		if m.commandError != "" {
			errorText := fmt.Sprintf("Error: %s\n\n", m.commandError)
			inputSection = fmt.Sprintf("%sCommand: %s\n\n%s\n\n%s",
				lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errorText),
				lipgloss.NewStyle().Align(lipgloss.Center).Render(m.textInput.View()),
				styles.NeutralStyle.Render("Try 'devmode' to activate dev mode, 'nodev' to continue without dev mode, or 'help' for all commands"),
				styles.QuitStyle.Render("Press Ctrl+C or Esc to quit"))
		} else {
			helpText := "Type 'dev' to activate dev mode, 'nodev' to continue without dev mode, or 'help' for all commands"
			inputSection = fmt.Sprintf("Command: %s\n\n%s\n\n%s",
				lipgloss.NewStyle().Align(lipgloss.Center).Render(m.textInput.View()),
				styles.NeutralStyle.Render(helpText),
				styles.QuitStyle.Render("Press Ctrl+C or Esc to quit"))
		}

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
		datePickerPrompt := "Please select a date to continue: ← → for days, ↑ ↓ for months, Enter/Space to select"
		prompt := styles.NeutralStyle.Width(contentWidth).Render(datePickerPrompt)

		currentDate := m.datePicker.GetSelectedDate()
		initialDateMessage := fmt.Sprint(currentDate.Format("January 2, 2006"))
		dateInfo := styles.WelcomeStyle.Width(contentWidth).Render(initialDateMessage)

		if m.isDev {
			content = folderInfo + "\n\n" + prompt + "\n" + dateInfo
		} else {
			content = prompt + "\n" + dateInfo
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
