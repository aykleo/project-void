package commands

import (
	"fmt"
	"sort"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Action      string
}

type Registry struct {
	commands map[string]Command
}

func NewRegistry() *Registry {
	registry := &Registry{
		commands: make(map[string]Command),
	}

	registry.RegisterCommand("help", "Show all available commands", "help")
	registry.RegisterCommand("start", "Go to the home page to start working", "start")
	registry.RegisterCommand("reset", "Go back to the welcome screen", "reset")
	registry.RegisterCommand("dev", "Activate development mode to select a local git repository", "dev")
	registry.RegisterCommand("nodev", "Continue without development mode", "nodev")
	registry.RegisterCommand("quit", "Exit the application", "quit")
	registry.RegisterCommand("exit", "Exit the application", "quit")

	return registry
}

func (r *Registry) RegisterCommand(name, description, action string) {
	r.commands[name] = Command{
		Name:        name,
		Description: description,
		Action:      action,
	}
}

func (r *Registry) GetCommand(name string) (Command, bool) {
	cleanName := strings.TrimPrefix(name, "./")
	cmd, exists := r.commands[cleanName]
	return cmd, exists
}

func (r *Registry) GetAllCommands() []Command {
	var commands []Command
	for _, cmd := range r.commands {
		commands = append(commands, cmd)
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands
}

func (r *Registry) GetHelpText() string {
	var help strings.Builder
	help.WriteString("Available Commands:\n\n")

	commands := r.GetAllCommands()
	for _, cmd := range commands {
		help.WriteString(fmt.Sprintf("  %s - %s\n", cmd.Name, cmd.Description))
	}

	help.WriteString("\nGit Commands (available in development mode):\n")
	help.WriteString("  git a - Clear author filter and show all commits\n")
	help.WriteString("  git a <name> - Filter commits by author name\n")
	help.WriteString("  git a <name1,name2> - Filter commits by multiple author names\n")
	help.WriteString("  Examples:\n")
	help.WriteString("    git a - Show all commits (clear filter)\n")
	help.WriteString("    git a john - Show commits by authors containing 'john'\n")
	help.WriteString("    git a john,alice - Show commits by authors containing 'john' or 'alice'\n")

	return help.String()
}

func (r *Registry) ValidateCommand(input string) (Command, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return Command{}, fmt.Errorf("command cannot be empty")
	}

	if strings.HasPrefix(input, "git a ") || input == "git a" {
		authorNames := ""
		if strings.HasPrefix(input, "git a ") {
			authorNames = strings.TrimPrefix(input, "git a ")
			authorNames = strings.TrimSpace(authorNames)
		}

		if authorNames == "" {
			return Command{
				Name:        "git a",
				Description: "Clear author filter and show all commits",
				Action:      "clear_author_filter",
			}, nil
		}

		return Command{
			Name:        "git a " + authorNames,
			Description: "Filter commits by author name(s)",
			Action:      "filter_by_author",
		}, nil
	}

	cleanName := strings.TrimPrefix(input, "./")

	cmd, exists := r.commands[cleanName]
	if !exists {
		return Command{}, fmt.Errorf("unknown command: %s\nType 'help' to see available commands\nOr use 'git a <name>' to filter commits by author", cleanName)
	}

	return cmd, nil
}

var GlobalRegistry = NewRegistry()

func RegisterCommand(name, description, action string) {
	GlobalRegistry.RegisterCommand(name, description, action)
}

func GetCommand(name string) (Command, bool) {
	return GlobalRegistry.GetCommand(name)
}

func GetAllCommands() []Command {
	return GlobalRegistry.GetAllCommands()
}

func GetHelpText() string {
	return GlobalRegistry.GetHelpText()
}

func ValidateCommand(input string) (Command, error) {
	return GlobalRegistry.ValidateCommand(input)
}

func GetAuthorNamesFromCommand(commandName string) []string {
	if !strings.HasPrefix(commandName, "git a ") {
		return nil
	}

	authorNames := strings.TrimPrefix(commandName, "git a ")
	authorNames = strings.TrimSpace(authorNames)

	if authorNames == "" {
		return nil
	}

	names := strings.Split(authorNames, ",")
	var cleanNames []string
	for _, name := range names {
		cleanName := strings.TrimSpace(name)
		if cleanName != "" {
			cleanNames = append(cleanNames, cleanName)
		}
	}

	return cleanNames
}
