package commitstable

import (
	"math/rand"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.progress.Width = msg.Width - 20
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}

		tableHeight := m.height - 4
		if tableHeight < 6 {
			tableHeight = 6
		}
		m.table.SetHeight(tableHeight)

		if m.width > 0 {
			columns := getCommitTableColumns(m.width)
			m.table.SetColumns(columns)
		}

	case LoadCommitsProgressMsg:
		if m.loadingState == LoadingInProgress {
			cmd = m.progress.SetPercent(msg.Percent)
			return m, cmd
		}

	case tickMsg:
		if m.loadingState == LoadingInProgress {
			currentPercent := m.progress.Percent()

			if currentPercent < 0.95 {
				maxIncrement := 0.1 * (1.0 - currentPercent)
				increment := rand.Float64() * maxIncrement

				newPercent := currentPercent + increment
				if newPercent > 0.95 {
					newPercent = 0.95
				}

				return m, tea.Batch(tickCmd(), m.progress.SetPercent(newPercent))
			}

			return m, tickCmd()
		}

	case LoadingCompleteMsg:
		if m.loadingState == LoadingInProgress {
			m.loadingState = LoadingComplete
			return m, m.progress.SetPercent(1.0)
		}

	case progress.FrameMsg:
		if m.loadingState == LoadingInProgress || m.progress.Percent() < 1.0 {
			progressModel, cmd := m.progress.Update(msg)
			m.progress = progressModel.(progress.Model)
			return m, cmd
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
