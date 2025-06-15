package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (h CommandHandler) Update(msg tea.Msg) (CommandHandler, tea.Cmd, *CommandResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > 60 {
			h.textInput.Width = 50
		} else {
			h.textInput.Width = msg.Width - 10
		}
		return h, nil, nil

	case tea.KeyMsg:
		if !h.enabled {
			return h, nil, nil
		}

		if h.showingHelp {
			h.showingHelp = false
			h.commandError = ""
			return h, nil, nil
		}

		if h.showingGitHelp {
			h.showingGitHelp = false
			h.commandError = ""
			return h, nil, nil
		}

		if h.showingCommand {
			switch msg.Type {
			case tea.KeyEnter:
				return h.processCommand()
			case tea.KeyCtrlC:
				return h, tea.Quit, &CommandResult{ShouldQuit: true}
			case tea.KeyEsc:
				h.showingCommand = false
				h.commandError = ""
				h.textInput.SetValue("")
				return h, nil, nil
			}

			if msg.String() == "'" {
				h.showingCommand = false
				h.commandError = ""
				h.textInput.SetValue("")
				return h, nil, nil
			}

			var cmd tea.Cmd
			h.textInput, cmd = h.textInput.Update(msg)
			return h, cmd, nil
		}

		if msg.String() == "c" {
			h.showingCommand = true
			h.textInput.Focus()
			return h, nil, nil
		}

		switch msg.Type {
		case tea.KeyEnter:
			return h.processCommand()
		case tea.KeyCtrlC, tea.KeyEsc:
			return h, tea.Quit, &CommandResult{ShouldQuit: true}
		}
	}

	if h.enabled && !h.showingHelp && !h.showingCommand && !h.showingGitHelp {
		var cmd tea.Cmd
		h.textInput, cmd = h.textInput.Update(msg)
		return h, cmd, nil
	}

	return h, nil, nil
}

func (h StatisticsCommandHandler) Update(msg tea.Msg) (StatisticsCommandHandler, tea.Cmd, *CommandResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > 60 {
			h.textInput.Width = 50
		} else {
			h.textInput.Width = msg.Width - 10
		}
		return h, nil, nil

	case tea.KeyMsg:
		if !h.enabled {
			return h, nil, nil
		}

		if h.showingHelp {
			h.showingHelp = false
			h.commandError = ""
			return h, nil, nil
		}

		if h.showingGitHelp {
			h.showingGitHelp = false
			h.commandError = ""
			return h, nil, nil
		}

		if h.showingCommand {
			switch msg.Type {
			case tea.KeyEnter:
				return h.processCommand()
			case tea.KeyCtrlC:
				return h, tea.Quit, &CommandResult{ShouldQuit: true}
			case tea.KeyEsc:
				h.showingCommand = false
				h.commandError = ""
				h.textInput.SetValue("")
				return h, nil, nil
			}

			if msg.String() == "'" {
				h.showingCommand = false
				h.commandError = ""
				h.textInput.SetValue("")
				return h, nil, nil
			}

			var cmd tea.Cmd
			h.textInput, cmd = h.textInput.Update(msg)
			return h, cmd, nil
		}

		if msg.String() == "c" {
			h.showingCommand = true
			h.textInput.Focus()
			return h, nil, nil
		}

		switch msg.Type {
		case tea.KeyEnter:
			return h.processCommand()
		case tea.KeyCtrlC, tea.KeyEsc:
			return h, tea.Quit, &CommandResult{ShouldQuit: true}
		}
	}

	if h.enabled && !h.showingHelp && !h.showingCommand && !h.showingGitHelp {
		var cmd tea.Cmd
		h.textInput, cmd = h.textInput.Update(msg)
		return h, cmd, nil
	}

	return h, nil, nil
}
