package slack

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	Token        string
	FilterByUser bool
	UserID       string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Token:        os.Getenv("SLACK_TOKEN"),
		FilterByUser: true,
		UserID:       os.Getenv("SLACK_USER_ID"),
	}

	if config.Token == "" {
		return nil, fmt.Errorf("SLACK_TOKEN is required")
	}

	filterByUserStr := os.Getenv("SLACK_FILTER_BY_USER")
	if filterByUserStr != "" {
		config.FilterByUser = strings.ToLower(filterByUserStr) == "true"
	}

	return config, nil
}

func NewClientFromConfig(config *Config) *SlackClient {
	return &SlackClient{
		Token:        config.Token,
		Client:       &http.Client{Timeout: 30 * time.Second},
		channelCache: make(map[string]string),
	}
}
