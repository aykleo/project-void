package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

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

	help.WriteString(sectionHeaderStyle.Render("\nGit Configuration Commands:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			commandStyle.Render("git status"),
			descStyle.Render("Show current Git repository configuration"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("git repo"),
			argStyle.Render("<url-or-path>"),
			descStyle.Render("Set Git repository URL or local path"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("git token"),
			argStyle.Render("<github-token>"),
			descStyle.Render("Set GitHub API token (increases rate limit from 60 to 5,000/hour)"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nGit Analysis Commands (available in development mode):") + "\n")
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

	help.WriteString(sectionHeaderStyle.Render("\nNavigation Commands:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			commandStyle.Render("start"),
			descStyle.Render("Start analyzing your data (uses current date by default)"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			commandStyle.Render("reset"),
			descStyle.Render("Return to welcome screen"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			commandStyle.Render("help"),
			descStyle.Render("Show this help message"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			commandStyle.Render("quit"),
			descStyle.Render("Exit the application"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nDate Commands:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("void sd"),
			argStyle.Render("<YYYY-MM-DD>"),
			descStyle.Render("Set analysis date (e.g., void sd 2025-06-01)"),
		),
	)

	return help.String()
}

func (r *Registry) ValidateCommand(input string) (Command, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return Command{}, fmt.Errorf("command cannot be empty")
	}

	if strings.HasPrefix(input, "git ") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			return Command{}, fmt.Errorf("git command requires a subcommand. Use 'help' to see available commands")
		}

		subCommand := parts[1]

		if subCommand == "status" {
			return Command{
				Name:        "git status",
				Description: "Show current Git repository configuration",
				Action:      "git_status",
			}, nil
		}

		if subCommand == "repo" || subCommand == "repository" {
			if len(parts) < 3 {
				return Command{}, fmt.Errorf("git %s command requires a value. Usage: git %s <url-or-path>", subCommand, subCommand)
			}

			value := strings.Join(parts[2:], " ")
			return Command{
				Name:        fmt.Sprintf("git repo %s", value),
				Description: "Set Git repository URL or path",
				Action:      "git_set_repo",
			}, nil
		}

		if subCommand == "token" {
			if len(parts) < 3 {
				return Command{}, fmt.Errorf("git token command requires a value. Usage: git token <github-api-token>")
			}

			value := strings.Join(parts[2:], " ")
			return Command{
				Name:        fmt.Sprintf("git token %s", value),
				Description: "Set GitHub API token for increased rate limits",
				Action:      "git_set_token",
			}, nil
		}

		if subCommand == "a" {
			authorNames := ""
			if len(parts) > 2 {
				authorNames = strings.Join(parts[2:], " ")
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

		return Command{}, fmt.Errorf("unknown git subcommand: %s\nAvailable: status, repo, token, a", subCommand)
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

	if strings.HasPrefix(input, "void ") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			return Command{}, fmt.Errorf("void command requires a subcommand. Use 'help' to see available commands")
		}

		subCommand := parts[1]

		if subCommand == "sd" {
			if len(parts) < 3 {
				return Command{}, fmt.Errorf("void sd command requires a date. Usage: void sd <YYYY-MM-DD>")
			}

			dateStr := parts[2]
			_, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return Command{}, fmt.Errorf("invalid date format. Use YYYY-MM-DD format (e.g., void sd 2025-06-01)")
			}

			return Command{
				Name:        fmt.Sprintf("void sd %s", dateStr),
				Description: "Set analysis date",
				Action:      "void_set_date",
			}, nil
		}

		return Command{}, fmt.Errorf("unknown void subcommand: %s\nAvailable: sd", subCommand)
	}

	cleanName := strings.TrimPrefix(input, "./")

	cmd, exists := r.commands[cleanName]
	if !exists {
		return Command{}, fmt.Errorf("unknown command: %s\nType 'help' to see available commands\nOr use 'git a <name>' to filter commits by author\nOr use 'git repo <url>' to set repository\nOr use 'jira <setting> <value>' to configure JIRA", cleanName)
	}

	return cmd, nil
}

func GetGitConfigValue(commandName string) (string, string) {
	if !strings.HasPrefix(commandName, "git ") {
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

	authorPart := strings.TrimPrefix(commandName, "git a ")
	authorPart = strings.TrimSpace(authorPart)

	if authorPart == "" {
		return nil
	}

	authors := strings.Split(authorPart, ",")
	var cleanAuthors []string
	for _, author := range authors {
		cleanAuthor := strings.TrimSpace(author)
		if cleanAuthor != "" {
			cleanAuthors = append(cleanAuthors, cleanAuthor)
		}
	}

	return cleanAuthors
}

func GetDateFromCommand(commandName string) (time.Time, error) {
	if !strings.HasPrefix(commandName, "void sd ") {
		return time.Time{}, fmt.Errorf("not a void sd command")
	}

	datePart := strings.TrimPrefix(commandName, "void sd ")
	datePart = strings.TrimSpace(datePart)

	if datePart == "" {
		return time.Time{}, fmt.Errorf("no date provided")
	}

	parsedDate, err := time.Parse("2006-01-02", datePart)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	return parsedDate, nil
}
