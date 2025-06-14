package home

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	Tabs       []string
	TabContent []string
	ActiveTab  int
	Width      int
}

func InitialModel(tabs []string, tabContent []string) model {
	return model{
		Tabs:       tabs,
		TabContent: tabContent,
		ActiveTab:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		return m, nil
	}
	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 4, 1, 4)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	welcomeStyle      = lipgloss.NewStyle().Foreground(highlightColor).Bold(true).Align(lipgloss.Center).MarginBottom(2)
)

func (m model) View() string {
	doc := strings.Builder{}

	horizontalPadding := 4
	contentWidth := m.Width - (horizontalPadding * 2)

	welcomeMessage := "Welcome to Project Void"
	doc.WriteString(welcomeStyle.Width(contentWidth).Render(welcomeMessage))
	doc.WriteString("\n\n")

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.ActiveTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	tabsWidth := lipgloss.Width(row)

	doc.WriteString(lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(row))
	doc.WriteString("\n")

	doc.WriteString(lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(
		windowStyle.Width(tabsWidth - windowStyle.GetHorizontalFrameSize()).Render(m.TabContent[m.ActiveTab])))

	return docStyle.Width(m.Width).Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
