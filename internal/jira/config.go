package jira

import (
	"fmt"
	"net/http"
	"os"
	"project-void/internal/config"
	"strings"
	"time"
)

type Config struct {
	BaseURL        string
	Username       string
	ApiToken       string
	ProjectKeys    []string
	FilterByUser   bool
	UserFilterType string
}

func LoadConfig() (*Config, error) {

	userConfig, err := config.LoadUserConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load user config: %w", err)
	}

	jiraConfig := &Config{
		BaseURL:        userConfig.Jira.BaseURL,
		Username:       userConfig.Jira.Username,
		ApiToken:       userConfig.Jira.ApiToken,
		ProjectKeys:    userConfig.Jira.ProjectKeys,
		FilterByUser:   userConfig.Jira.FilterByUser,
		UserFilterType: userConfig.Jira.UserFilterType,
	}

	if jiraConfig.BaseURL == "" {
		jiraConfig.BaseURL = os.Getenv("JIRA_BASE_URL")
	}
	if jiraConfig.Username == "" {
		jiraConfig.Username = os.Getenv("JIRA_USERNAME")
	}
	if jiraConfig.ApiToken == "" {
		jiraConfig.ApiToken = os.Getenv("JIRA_API_TOKEN")
	}
	if len(jiraConfig.ProjectKeys) == 0 {
		projectKeysStr := os.Getenv("JIRA_PROJECT_KEY")
		if projectKeysStr != "" {
			keys := strings.Split(projectKeysStr, ",")
			jiraConfig.ProjectKeys = make([]string, len(keys))
			for i, key := range keys {
				jiraConfig.ProjectKeys[i] = strings.TrimSpace(key)
			}
		}
	}

	filterByUserStr := os.Getenv("JIRA_FILTER_BY_USER")
	if filterByUserStr != "" {
		jiraConfig.FilterByUser = strings.ToLower(filterByUserStr) == "true"
	}

	if jiraConfig.FilterByUser {
		userFilterType := os.Getenv("JIRA_USER_FILTER_TYPE")
		if userFilterType != "" {
			jiraConfig.UserFilterType = userFilterType
		}
	}

	if jiraConfig.UserFilterType == "" {
		jiraConfig.UserFilterType = "participant"
	}

	if jiraConfig.BaseURL == "" {
		return nil, fmt.Errorf("JIRA_BASE_URL is required (set with: jira url <your-jira-url>)")
	}
	if jiraConfig.Username == "" {
		return nil, fmt.Errorf("JIRA_USERNAME is required (set with: jira user <your-username>)")
	}
	if jiraConfig.ApiToken == "" {
		return nil, fmt.Errorf("JIRA_API_TOKEN is required (set with: jira token <your-token>)")
	}

	validFilterTypes := map[string]bool{
		"assignee":    true,
		"reporter":    true,
		"participant": true,
		"all":         true,
	}

	if !validFilterTypes[jiraConfig.UserFilterType] {
		return nil, fmt.Errorf("invalid user filter type: %s. Valid options are: assignee, reporter, participant, all", jiraConfig.UserFilterType)
	}

	return jiraConfig, nil
}

func NewClientFromConfig(config *Config) *JiraClient {
	return &JiraClient{
		BaseURL:  config.BaseURL,
		Username: config.Username,
		ApiToken: config.ApiToken,
		Client:   &http.Client{Timeout: 30 * time.Second},
	}
}
