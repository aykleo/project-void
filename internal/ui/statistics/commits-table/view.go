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

		if m.width > 0 {
			content = lipgloss.NewStyle().
				Width(m.width).
				Align(lipgloss.Center).
				Render(content)
		}

		return content
	}

	if m.loadingState == LoadingError {
		errorText := fmt.Sprintf("Error loading commits: %s", m.loadError)

		if m.width > 0 {
			content := lipgloss.NewStyle().
				Width(m.width).
				Align(lipgloss.Center).
				Foreground(lipgloss.Color("196")).
				Render(errorText)
			return content
		}

		content := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errorText)
		return content
	}

	var tableView string
	if m.borderFocused {
		style := focusedStyle
		if m.width > 0 {
			style = style.Width(m.width)
		}
		tableView = style.Render(m.table.View())
	} else {
		style := baseStyle
		if m.width > 0 {
			style = style.Width(m.width)
		}
		tableView = style.Render(m.table.View())
	}

	return tableView
}
