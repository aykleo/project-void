package jira

import (
	"fmt"
	"net/http"
	"os"
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
	config := &Config{
		BaseURL:        os.Getenv("JIRA_BASE_URL"),
		Username:       os.Getenv("JIRA_USERNAME"),
		ApiToken:       os.Getenv("JIRA_API_TOKEN"),
		FilterByUser:   false,
		UserFilterType: "participant",
	}

	if config.BaseURL == "" {
		return nil, fmt.Errorf("JIRA_BASE_URL is required")
	}
	if config.Username == "" {
		return nil, fmt.Errorf("JIRA_USERNAME is required")
	}
	if config.ApiToken == "" {
		return nil, fmt.Errorf("JIRA_API_TOKEN is required")
	}

	projectKeysStr := os.Getenv("JIRA_PROJECT_KEY")
	if projectKeysStr != "" {
		keys := strings.Split(projectKeysStr, ",")
		config.ProjectKeys = make([]string, len(keys))
		for i, key := range keys {
			config.ProjectKeys[i] = strings.TrimSpace(key)
		}
	}

	filterByUserStr := os.Getenv("JIRA_FILTER_BY_USER")
	if filterByUserStr != "" {
		config.FilterByUser = strings.ToLower(filterByUserStr) == "true"
	}

	if config.FilterByUser {
		userFilterType := os.Getenv("JIRA_USER_FILTER_TYPE")
		if userFilterType != "" {
			config.UserFilterType = userFilterType
		}
	}

	validFilterTypes := map[string]bool{
		"assignee":    true,
		"reporter":    true,
		"participant": true,
		"all":         true,
	}

	if !validFilterTypes[config.UserFilterType] {
		return nil, fmt.Errorf("invalid user filter type: %s. Valid options are: assignee, reporter, participant, all", config.UserFilterType)
	}

	return config, nil
}

func NewClientFromConfig(config *Config) *JiraClient {
	return &JiraClient{
		BaseURL:  config.BaseURL,
		Username: config.Username,
		ApiToken: config.ApiToken,
		Client:   &http.Client{Timeout: 30 * time.Second},
	}
}
