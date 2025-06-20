package commands

import "strings"

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

	if key == "repo" {
		specialCommands := []string{"list", "ls", "clear", "reset", "remove", "rm"}
		for _, special := range specialCommands {
			if strings.HasPrefix(value, special) {
				return "", ""
			}
		}
	}

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
