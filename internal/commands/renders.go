package commands

import (
	"fmt"
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

func (r *Registry) GetHelpText() string {
	var help strings.Builder

	help.WriteString(sectionHeaderStyle.Render("\nGeneral Commands:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("void"),
			sectionHeaderStyle.Render("start"),
			descStyle.Render("Start analyzing your data (uses current date by default)"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("void"),
			sectionHeaderStyle.Render("reset"),
			descStyle.Render("Return to welcome screen"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("void"),
			sectionHeaderStyle.Render("help"),
			descStyle.Render("Show this help message"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("void"),
			sectionHeaderStyle.Render("help git"),
			descStyle.Render("Show Git help and setup instructions"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("void"),
			sectionHeaderStyle.Render("help jira"),
			descStyle.Render("Show JIRA help and setup instructions"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("void"),
			sectionHeaderStyle.Render("quit"),
			descStyle.Render("Exit the application"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("void"),
			sectionHeaderStyle.Render("set-date"),
			argStyle.Render("<YYYY-MM-DD>"),
			descStyle.Render("Set analysis date (e.g., void sd 2025-06-01)"),
		),
	)

	return help.String()
}

func (r *Registry) GetGitHelpText() string {
	var help strings.Builder

	help.WriteString(sectionHeaderStyle.Render("\nGit Configuration Commands:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("git"),
			sectionHeaderStyle.Render("status"),
			descStyle.Render("Show current Git repository configuration and token status"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("git"),
			sectionHeaderStyle.Render("repo list"),
			descStyle.Render("List all configured Git repositories"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("git"),
			sectionHeaderStyle.Render("repo clear"),
			descStyle.Render("Clear all Git repository configurations"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("git"),
			sectionHeaderStyle.Render("repo"),
			argStyle.Render("<url-or-path>"),
			descStyle.Render("Add Git repository URL or local path"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("git"),
			sectionHeaderStyle.Render("repo remove"),
			argStyle.Render("<url-or-path>"),
			descStyle.Render("Remove specific Git repository from configuration"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("git"),
			sectionHeaderStyle.Render("token"),
			argStyle.Render("<github-token>"),
			descStyle.Render("Set GitHub API token (needed for remote repositories)"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nRepository Types:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("Local"),
			descStyle.Render("/path/to/repo, C:\\path\\to\\repo, ~/workspace/repo"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("Remote"),
			descStyle.Render("https://github.com/user/repo, git@github.com:user/repo.git"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nMulti-Repository Support:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s\n",
			descStyle.Render("• Add multiple repositories to track commits from all of them"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s\n",
			descStyle.Render("• Commits are aggregated and deduplicated across repositories"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s\n",
			descStyle.Render("• Repository names are shown in the branch column for identification"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nGitHub API Token Setup:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("1. Visit"),
			descStyle.Render("https://github.com/settings/tokens"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("2. Click"),
			descStyle.Render("Generate new token (classic)"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("3. Set note"),
			descStyle.Render("Project Void Analytics"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("4. Set expiration"),
			descStyle.Render("90 days or No expiration"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("5. Select scopes"),
			descStyle.Render("See permissions below"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("6. Generate and copy"),
			descStyle.Render("Token will only be shown once"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nRequired Token Permissions:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			argStyle.Render("public_repo"),
			descStyle.Render("Access public repositories"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			argStyle.Render("repo"),
			descStyle.Render("Full repository access (required for private repos)"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nGit Analysis Commands:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("git"),
			sectionHeaderStyle.Render("author"),
			descStyle.Render("Clear author filter and show all commits"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("git"),
			sectionHeaderStyle.Render("author"),
			argStyle.Render("<name>"),
			descStyle.Render("Filter commits by author name. Comma-separated for multiple authors."),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("git"),
			sectionHeaderStyle.Render("branch"),
			argStyle.Render("<name>"),
			descStyle.Render("Filter commits by branch name. Comma-separated for multiple branches."),
		),
	)

	return help.String()
}

func (r *Registry) GetJiraHelpText() string {
	var help strings.Builder

	help.WriteString(sectionHeaderStyle.Render("\nJIRA Configuration Commands:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("jira"),
			sectionHeaderStyle.Render("status"),
			descStyle.Render("Show current JIRA configuration and status"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("jira"),
			sectionHeaderStyle.Render("url"),
			argStyle.Render("<url>"),
			descStyle.Render("Set JIRA base URL (e.g., https://your-domain.atlassian.net)"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("jira"),
			sectionHeaderStyle.Render("user"),
			argStyle.Render("<username>"),
			descStyle.Render("Set JIRA username (usually your email address)"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("jira"),
			sectionHeaderStyle.Render("token"),
			argStyle.Render("<token>"),
			descStyle.Render("Set JIRA API token (see setup instructions below)"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s %s - %s\n",
			commandStyle.Render("jira"),
			sectionHeaderStyle.Render("project"),
			argStyle.Render("<key>"),
			descStyle.Render("Set JIRA project key(s). Comma-separated for multiple projects (e.g., TIP,SP)"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("jira"),
			sectionHeaderStyle.Render("f"),
			descStyle.Render("Enable user filtering (show only your issues)"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s %s - %s\n",
			commandStyle.Render("jira"),
			sectionHeaderStyle.Render("nof"),
			descStyle.Render("Disable user filtering (show all issues)"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nJIRA API Token Setup:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("1. Visit"),
			descStyle.Render("https://id.atlassian.com/manage-profile/security/api-tokens"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("2. Click"),
			descStyle.Render("Create API token"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("3. Set note"),
			descStyle.Render("Project Void"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("4. Copy and save"),
			descStyle.Render("Token will only be shown once"),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			descStyle.Render("5. Set in app"),
			descStyle.Render("Use: jira token <your-token>"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nProject Keys:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s - %s\n",
			argStyle.Render("<key>"),
			descStyle.Render("Project key (e.g., TIP). For multiple, separate with commas: TIP,SP"),
		),
	)

	help.WriteString(sectionHeaderStyle.Render("\nTroubleshooting:") + "\n")
	help.WriteString(
		fmt.Sprintf("  %s\n",
			descStyle.Render("If you see incomplete configuration, set all required fields: url, user, token, and at least one project key."),
		),
	)
	help.WriteString(
		fmt.Sprintf("  %s\n",
			descStyle.Render("You can check your current configuration with: jira status"),
		),
	)

	return help.String()
}
