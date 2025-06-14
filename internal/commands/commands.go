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

	return help.String()
}

func (r *Registry) ValidateCommand(input string) (Command, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return Command{}, fmt.Errorf("command cannot be empty")
	}

	cleanName := strings.TrimPrefix(input, "./")

	cmd, exists := r.commands[cleanName]
	if !exists {
		return Command{}, fmt.Errorf("unknown command: %s\nType 'help' to see available commands", cleanName)
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
