package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type UserConfig struct {
	Jira JiraConfig `json:"jira"`
}

type JiraConfig struct {
	BaseURL        string   `json:"base_url"`
	Username       string   `json:"username"`
	ApiToken       string   `json:"api_token"`
	ProjectKeys    []string `json:"project_keys"`
	FilterByUser   bool     `json:"filter_by_user"`
	UserFilterType string   `json:"user_filter_type"`
}

const configFileName = ".project-void-config.json"

func GetConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

func LoadUserConfig() (*UserConfig, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &UserConfig{
			Jira: JiraConfig{
				FilterByUser:   false,
				UserFilterType: "participant",
			},
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config UserConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func SaveUserConfig(config *UserConfig) error {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func SetJiraConfig(key, value string) error {
	config, err := LoadUserConfig()
	if err != nil {
		return err
	}

	switch key {
	case "url", "baseurl", "base_url":
		config.Jira.BaseURL = value
	case "user", "username":
		config.Jira.Username = value
	case "token", "apitoken", "api_token":
		config.Jira.ApiToken = value
	case "project", "projects", "project_keys":
		if value == "" {
			config.Jira.ProjectKeys = nil
		} else {
			keys := strings.Split(value, ",")
			var cleanKeys []string
			for _, key := range keys {
				cleanKey := strings.TrimSpace(key)
				if cleanKey != "" {
					cleanKeys = append(cleanKeys, cleanKey)
				}
			}
			config.Jira.ProjectKeys = cleanKeys
		}
	default:
		return fmt.Errorf("unknown JIRA config key: %s", key)
	}

	return SaveUserConfig(config)
}

func GetJiraConfigStatus() (string, error) {
	config, err := LoadUserConfig()
	if err != nil {
		return "", err
	}

	var status strings.Builder
	status.WriteString("Current JIRA Configuration:\n\n")

	if config.Jira.BaseURL != "" {
		status.WriteString(fmt.Sprintf("  URL: %s\n", config.Jira.BaseURL))
	} else {
		status.WriteString("  URL: (not set)\n")
	}

	if config.Jira.Username != "" {
		status.WriteString(fmt.Sprintf("  Username: %s\n", config.Jira.Username))
	} else {
		status.WriteString("  Username: (not set)\n")
	}

	if config.Jira.ApiToken != "" {
		masked := config.Jira.ApiToken
		if len(masked) > 8 {
			masked = masked[:4] + "..." + masked[len(masked)-4:]
		}
		status.WriteString(fmt.Sprintf("  API Token: %s\n", masked))
	} else {
		status.WriteString("  API Token: (not set)\n")
	}

	if len(config.Jira.ProjectKeys) > 0 {
		status.WriteString(fmt.Sprintf("  Projects: %s\n", strings.Join(config.Jira.ProjectKeys, ", ")))
	} else {
		status.WriteString("  Projects: (not set)\n")
	}

	if config.Jira.BaseURL != "" && config.Jira.Username != "" && config.Jira.ApiToken != "" {
		status.WriteString("\n✓ Configuration is complete and ready to use!")
	} else {
		status.WriteString("\n⚠ Configuration is incomplete. Missing:")
		if config.Jira.BaseURL == "" {
			status.WriteString("\n  - URL (use: jira url <your-jira-url>)")
		}
		if config.Jira.Username == "" {
			status.WriteString("\n  - Username (use: jira user <your-username>)")
		}
		if config.Jira.ApiToken == "" {
			status.WriteString("\n  - API Token (use: jira token <your-token>)")
		}
	}

	return status.String(), nil
}
