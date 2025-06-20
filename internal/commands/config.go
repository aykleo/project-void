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
		if strings.HasPrefix(value, "remove ") || strings.HasPrefix(value, "rm ") {
			urlPart := strings.TrimPrefix(value, "remove ")
			urlPart = strings.TrimPrefix(urlPart, "rm ")
			urlPart = strings.TrimSpace(urlPart)
			if urlPart != "" {
				return "repo", urlPart
			}
			return "", ""
		}

		specialCommands := []string{"list", "ls", "clear", "reset"}
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
