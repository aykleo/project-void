package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"
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

	registry.RegisterCommand("void help", "Show all available commands", "help")
	registry.RegisterCommand("void h", "Show all available commands", "help")
	registry.RegisterCommand("void start", "Go to the statistics screen", "start")
	registry.RegisterCommand("void st", "Go to the statistics screen", "start")
	registry.RegisterCommand("void reset", "Go back to the starting screen", "reset")
	registry.RegisterCommand("void quit", "Exit the application", "quit")
	registry.RegisterCommand("void q", "Exit the application", "quit")

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

func (r *Registry) ValidateCommand(input string) (Command, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return Command{}, fmt.Errorf("command cannot be empty")
	}

	if strings.HasPrefix(input, "git ") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			return Command{}, fmt.Errorf("git command requires a subcommand. Use 'void help git' to see available git commands")
		}

		subCommand := parts[1]

		if subCommand == "status" || subCommand == "stat" {
			return Command{
				Name:        "git status",
				Description: "Show current Git repository configuration",
				Action:      "git_status",
			}, nil
		}

		if subCommand == "repo" || subCommand == "repository" || subCommand == "r" {
			if len(parts) == 2 {
				return Command{
					Name:        "git repo",
					Description: "Show current Git repository configuration",
					Action:      "git_status",
				}, nil
			}

			if parts[2] == "clear" || parts[2] == "reset" {
				return Command{
					Name:        "git repo clear",
					Description: "Clear all Git repositories",
					Action:      "git_clear_repo",
				}, nil
			}

			if parts[2] == "remove" || parts[2] == "rm" {
				if len(parts) < 4 {
					return Command{}, fmt.Errorf("git repo remove command requires a repository URL. Usage: git repo remove <url>")
				}
				value := strings.Join(parts[3:], " ")
				return Command{
					Name:        fmt.Sprintf("git repo remove %s", value),
					Description: "Remove a specific Git repository",
					Action:      "git_remove_repo",
				}, nil
			}

			if parts[2] == "list" || parts[2] == "ls" {
				return Command{
					Name:        "git repo list",
					Description: "List all configured Git repositories",
					Action:      "git_list_repos",
				}, nil
			}

			value := strings.Join(parts[2:], " ")
			return Command{
				Name:        fmt.Sprintf("git repo %s", value),
				Description: "Add Git repository URL or path",
				Action:      "git_set_repo",
			}, nil
		}

		if subCommand == "token" || subCommand == "t" {
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

		if subCommand == "author" || subCommand == "a" {
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

		if subCommand == "branch" || subCommand == "b" {
			branchNames := ""
			if len(parts) > 2 {
				branchNames = strings.Join(parts[2:], " ")
				branchNames = strings.TrimSpace(branchNames)
			}

			if branchNames == "" {
				return Command{
					Name:        "git b",
					Description: "Clear branch filter and show all commits",
					Action:      "clear_branch_filter",
				}, nil
			}

			return Command{
				Name:        "git b " + branchNames,
				Description: "Filter commits by branch name(s)",
				Action:      "filter_by_branch",
			}, nil
		}

		return Command{}, fmt.Errorf("unknown git subcommand: %s\nAvailable: status, repo, token, author, branch\nFor Git help, use: void help git", subCommand)
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

	if strings.HasPrefix(input, "git b ") || input == "git b" {
		branchNames := ""
		if strings.HasPrefix(input, "git b ") {
			branchNames = strings.TrimPrefix(input, "git b ")
			branchNames = strings.TrimSpace(branchNames)
		}

		if branchNames == "" {
			return Command{
				Name:        "git b",
				Description: "Clear branch filter and show all commits",
				Action:      "clear_branch_filter",
			}, nil
		}

		return Command{
			Name:        "git b " + branchNames,
			Description: "Filter commits by branch name(s)",
			Action:      "filter_by_branch",
		}, nil
	}

	if strings.HasPrefix(input, "jira ") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			return Command{}, fmt.Errorf("jira command requires a subcommand. Use 'help' to see available commands")
		}

		subCommand := parts[1]

		if subCommand == "status" || subCommand == "stat" {
			return Command{
				Name:        "jira status",
				Description: "Show current JIRA configuration status",
				Action:      "jira_status",
			}, nil
		}

		if subCommand == "f" || subCommand == "filter" {
			return Command{
				Name:        "jira f",
				Description: "Enable user filtering (show only your issues)",
				Action:      "jira_filter_on",
			}, nil
		}

		if subCommand == "nof" || subCommand == "nofilter" {
			return Command{
				Name:        "jira nof",
				Description: "Disable user filtering (show all issues)",
				Action:      "jira_filter_off",
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
		case "project", "projects", "p":
			return Command{
				Name:        fmt.Sprintf("jira project %s", value),
				Description: "Set JIRA project key(s)",
				Action:      "jira_set_project",
			}, nil
		default:
			return Command{}, fmt.Errorf("unknown jira subcommand: %s\nAvailable: url, user, token, project, status, f, nof", subCommand)
		}
	}

	if strings.HasPrefix(input, "void ") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			return Command{}, fmt.Errorf("void command requires a subcommand. Use 'void help' to see available commands")
		}

		subCommand := parts[1]

		if (subCommand == "help" || subCommand == "h") && len(parts) == 3 && (parts[2] == "git" || parts[2] == "g") {
			return Command{
				Name:        "void help git",
				Description: "Show Git help",
				Action:      "git_help",
			}, nil
		}

		if (subCommand == "help" || subCommand == "h") && len(parts) == 3 && (parts[2] == "jira" || parts[2] == "j") {
			return Command{
				Name:        "void help jira",
				Description: "Show JIRA help",
				Action:      "jira_help",
			}, nil
		}

		if subCommand == "set-date" || subCommand == "sd" {
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

		cmd, exists := r.commands[input]
		if exists {
			return cmd, nil
		}

		return Command{}, fmt.Errorf("unknown void subcommand: %s\nAvailable: help, start, reset, quit, set-date, help git, help jira", subCommand)
	}

	cleanName := strings.TrimPrefix(input, "./")

	cmd, exists := r.commands[cleanName]
	if !exists {
		return Command{}, fmt.Errorf("unknown command: %s\nType 'void help' to see available commands\nOr use 'void help git' to see Git commands\nOr use 'void help jira' to see JIRA commands", cleanName)
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

func GetGitHelpText() string {
	return GlobalRegistry.GetGitHelpText()
}

func GetJiraHelpText() string {
	return GlobalRegistry.GetJiraHelpText()
}

func ValidateCommand(input string) (Command, error) {
	return GlobalRegistry.ValidateCommand(input)
}
