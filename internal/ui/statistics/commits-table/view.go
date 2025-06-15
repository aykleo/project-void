package commitstable

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.loadingState == LoadingInProgress {
		var loadingText string
		if m.spinner != nil {
			loadingText = fmt.Sprintf("%s Loading commits...", m.spinner.View())
		} else {
			loadingText = "Loading commits..."
		}
		progressView := m.progress.View()

		content := lipgloss.JoinVertical(lipgloss.Center,
			loadingText,
			progressView,
		)

		return content
	}

	if m.loadingState == LoadingError {
		errorText := fmt.Sprintf("Error loading commits: %s", m.loadError)
		content := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errorText)

		return content
	}

	tableView := baseStyle.Render(m.table.View())

	if m.borderFocused {
		tableView = focusedStyle.Render(m.table.View())
	}

	if m.width > 0 {
		return lipgloss.NewStyle().Width(m.width).Render(tableView)
	}

	return tableView
}
