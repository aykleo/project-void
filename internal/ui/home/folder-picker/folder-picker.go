package folderpicker

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	filepicker     filepicker.Model
	selectedFolder string
	quitting       bool
	err            error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m Model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func InitialModel() Model {
	fp := filepicker.New()

	fp.DirAllowed = true
	fp.FileAllowed = true

	if wd, err := os.Getwd(); err == nil {
		fp.CurrentDirectory = wd
	} else if homeDir, err := os.UserHomeDir(); err == nil {
		fp.CurrentDirectory = homeDir
	} else {
		fp.CurrentDirectory = "."
	}

	m := Model{
		filepicker: fp,
	}
	return m
}

func (m *Model) RefreshDirectory() tea.Cmd {
	return m.filepicker.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case " ":
			m.selectedFolder = m.filepicker.CurrentDirectory
			return m, nil
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			m.selectedFolder = path
			return m, nil
		} else {

			m.err = errors.New("please select a directory, not a file")
			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not accessible")
		m.selectedFolder = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFolder == "" {
		s.WriteString("Pick a folder\n")
		s.WriteString("\n  Current: " + m.filepicker.Styles.Directory.Render(m.filepicker.CurrentDirectory))
		navigationHelp := "← →: navigate • enter: select folder • space: select current folder"
		s.WriteString("\n\n  " + lipgloss.NewStyle().Align(lipgloss.Center).Foreground(lipgloss.Color("8")).Render(navigationHelp))

	} else {
		s.WriteString("Selected folder: " + m.filepicker.Styles.Selected.Render(m.selectedFolder))
	}

	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

func (m Model) GetSelectedFolder() string {
	return m.selectedFolder
}
