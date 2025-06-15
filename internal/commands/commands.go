package commands

import (
	"fmt"
	"sort"
	"strings"

	styles "project-void/internal/ui/styles"

	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	sectionHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230"))
	commandStyle       = lipgloss.NewStyle().Foreground(styles.HighlightColor)
	argStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	descStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
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
	registry.RegisterCommand("dev", "Activate development mode to select a local git repository and see commits", "dev")
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
	help.WriteString(sectionHeaderStyle.Render("Available Commands:") + "\n")

	commands := r.GetAllCommands()
	for _, cmd := range commands {
		help.WriteString(fmt.Sprintf("  %s - %s\n",
			commandStyle.Render(cmd.Name),
			descStyle.Render(cmd.Description),
		))
	}

	help.WriteString(sectionHeaderStyle.Render("\nJIRA Configuration Commands:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			commandStyle.Render("jira status"),
			descStyle.Render("Show current JIRA configuration"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("jira url"),
			argStyle.Render("<url>"),
			descStyle.Render("Set JIRA base URL"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("jira user"),
			argStyle.Render("<username>"),
			descStyle.Render("Set JIRA username"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("jira token"),
			argStyle.Render("<token>"),
			descStyle.Render("Set JIRA API token"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("jira project"),
			argStyle.Render("<key>"),
			descStyle.Render("Set JIRA project key"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("jira project"),
			argStyle.Render("<key1,key2>"),
			descStyle.Render("Set multiple JIRA project keys"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nGit Commands (available in development mode):") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			commandStyle.Render("git a"),
			descStyle.Render("Clear author filter and show all commits"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("git a"),
			argStyle.Render("<name>"),
			descStyle.Render("Filter commits by author name"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("git a"),
			argStyle.Render("<name1,name2>"),
			descStyle.Render("Filter commits by multiple author names"),
		),
	)

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

	if strings.HasPrefix(input, "jira ") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			return Command{}, fmt.Errorf("jira command requires a subcommand. Use 'help' to see available commands")
		}

		subCommand := parts[1]

		if subCommand == "status" {
			return Command{
				Name:        "jira status",
				Description: "Show current JIRA configuration status",
				Action:      "jira_status",
			}, nil
		}

		if len(parts) < 3 {
			return Command{}, fmt.Errorf("jira %s command requires a value. Usage: jira %s <value>", subCommand, subCommand)
		}

		value := strings.Join(parts[2:], " ")

		switch subCommand {
		case "url":
			return Command{
				Name:        fmt.Sprintf("jira url %s", value),
				Description: "Set JIRA base URL",
				Action:      "jira_set_url",
			}, nil
		case "user", "username":
			return Command{
				Name:        fmt.Sprintf("jira user %s", value),
				Description: "Set JIRA username",
				Action:      "jira_set_user",
			}, nil
		case "token":
			return Command{
				Name:        fmt.Sprintf("jira token %s", value),
				Description: "Set JIRA API token",
				Action:      "jira_set_token",
			}, nil
		case "project", "projects":
			return Command{
				Name:        fmt.Sprintf("jira project %s", value),
				Description: "Set JIRA project key(s)",
				Action:      "jira_set_project",
			}, nil
		default:
			return Command{}, fmt.Errorf("unknown jira subcommand: %s\nAvailable: url, user, token, project, status", subCommand)
		}
	}

	cleanName := strings.TrimPrefix(input, "./")

	cmd, exists := r.commands[cleanName]
	if !exists {
		return Command{}, fmt.Errorf("unknown command: %s\nType 'help' to see available commands\nOr use 'git a <name>' to filter commits by author\nOr use 'jira <setting> <value>' to configure JIRA", cleanName)
	}

	return cmd, nil
}

func GetJiraConfigValue(commandName string) (string, string) {
	if !strings.HasPrefix(commandName, "jira ") {
		return "", ""
	}

	parts := strings.Fields(commandName)
	if len(parts) < 3 {
		return "", ""
	}

	key := parts[1]
	value := strings.Join(parts[2:], " ")
	return key, value
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
