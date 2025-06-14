package styles

import (
	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	HighlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#9f11d4"}
	DocStyle       = lipgloss.NewStyle().Padding(1, 4, 1, 4)
	WelcomeStyle   = lipgloss.NewStyle().Foreground(HighlightColor).Bold(true).Align(lipgloss.Center).MarginBottom(1)
	QuitStyle      = lipgloss.NewStyle().Blink(true).Italic(true).Align(lipgloss.Center).MarginBottom(2)
)
