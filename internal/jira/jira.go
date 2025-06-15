package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Issue struct {
	Key        string
	Status     string
	Assignee   string
	Summary    string
	Updated    time.Time
	Created    time.Time
	IssueType  string
	Priority   string
	UserAction string
	ActionDate time.Time
}

type JiraClient struct {
	BaseURL  string
	Username string
	ApiToken string
	Client   *http.Client
}

func NewJiraClient(baseURL, username, apiToken string) *JiraClient {
	return &JiraClient{
		BaseURL:  strings.TrimSuffix(baseURL, "/"),
		Username: username,
		ApiToken: apiToken,
		Client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *JiraClient) GetIssuesSince(since time.Time, config *Config) ([]Issue, error) {
	sinceStr := since.Format("2006-01-02")

	jql := fmt.Sprintf("updated >= '%s'", sinceStr)

	if len(config.ProjectKeys) > 0 {
		if len(config.ProjectKeys) == 1 {
			jql = fmt.Sprintf("project = '%s' AND %s", config.ProjectKeys[0], jql)
		} else {
			quotedKeys := make([]string, len(config.ProjectKeys))
			for i, key := range config.ProjectKeys {
				quotedKeys[i] = fmt.Sprintf("'%s'", key)
			}
			projectList := strings.Join(quotedKeys, ", ")
			jql = fmt.Sprintf("project IN (%s) AND %s", projectList, jql)
		}
	}

	if config.FilterByUser {
		username := config.Username

		var userConditions []string

		if strings.Contains(username, "@") {
			userConditions = []string{
				fmt.Sprintf("assignee = '%s'", username),
				fmt.Sprintf("reporter = '%s'", username),
			}

			localPart := strings.Split(username, "@")[0]
			userConditions = append(userConditions,
				fmt.Sprintf("assignee = '%s'", localPart),
				fmt.Sprintf("reporter = '%s'", localPart),
			)
		} else {
			userConditions = []string{
				fmt.Sprintf("assignee = '%s'", username),
				fmt.Sprintf("reporter = '%s'", username),
			}
		}

		userFilter := "(" + strings.Join(userConditions, " OR ") + ")"
		jql = fmt.Sprintf("(%s) AND %s", userFilter, jql)
	}

	jql += " ORDER BY updated DESC"

	allIssues := []Issue{}
	startAt := 0
	maxResults := 100

	for {
		url := fmt.Sprintf("%s/rest/api/2/search?jql=%s&startAt=%d&maxResults=%d",
			c.BaseURL, url.QueryEscape(jql), startAt, maxResults)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.SetBasicAuth(c.Username, c.ApiToken)
		req.Header.Set("Accept", "application/json")

		resp, err := c.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}

		var searchResult struct {
			Issues []struct {
				Key    string `json:"key"`
				Fields struct {
					Status struct {
						Name string `json:"name"`
					} `json:"status"`
					Assignee *struct {
						DisplayName  string `json:"displayName"`
						EmailAddress string `json:"emailAddress"`
					} `json:"assignee"`
					Summary   string `json:"summary"`
					Updated   string `json:"updated"`
					Created   string `json:"created"`
					IssueType struct {
						Name string `json:"name"`
					} `json:"issuetype"`
					Priority *struct {
						Name string `json:"name"`
					} `json:"priority"`
					Reporter *struct {
						DisplayName  string `json:"displayName"`
						EmailAddress string `json:"emailAddress"`
					} `json:"reporter"`
				} `json:"fields"`
			} `json:"issues"`
			Total      int `json:"total"`
			StartAt    int `json:"startAt"`
			MaxResults int `json:"maxResults"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		for _, jiraIssue := range searchResult.Issues {
			assignee := ""
			if jiraIssue.Fields.Assignee != nil {
				assignee = jiraIssue.Fields.Assignee.DisplayName
			}

			reporter := ""
			if jiraIssue.Fields.Reporter != nil {
				reporter = jiraIssue.Fields.Reporter.DisplayName
			}

			priority := ""
			if jiraIssue.Fields.Priority != nil {
				priority = jiraIssue.Fields.Priority.Name
			}

			userAction := c.determineUserAction(config, assignee, reporter)

			updatedTime, _ := time.Parse("2006-01-02T15:04:05.000-0700", jiraIssue.Fields.Updated)
			createdTime, _ := time.Parse("2006-01-02T15:04:05.000-0700", jiraIssue.Fields.Created)

			issue := Issue{
				Key:        jiraIssue.Key,
				Status:     jiraIssue.Fields.Status.Name,
				Assignee:   assignee,
				Summary:    jiraIssue.Fields.Summary,
				Updated:    updatedTime,
				Created:    createdTime,
				IssueType:  jiraIssue.Fields.IssueType.Name,
				Priority:   priority,
				UserAction: userAction,
				ActionDate: updatedTime,
			}

			allIssues = append(allIssues, issue)
		}

		if len(searchResult.Issues) < maxResults {
			break
		}
		startAt += maxResults
	}

	return allIssues, nil
}

func (c *JiraClient) determineUserAction(config *Config, assignee, reporter string) string {
	if !config.FilterByUser {
		return "All Issues"
	}

	currentUserEmail := config.Username
	var currentUserName string

	if strings.Contains(currentUserEmail, "@") {
		parts := strings.Split(currentUserEmail, "@")
		if len(parts) > 0 {
			namePart := parts[0]
			namePart = strings.ReplaceAll(namePart, ".", " ")
			currentUserName = namePart
		}
	}

	actions := []string{}

	if c.isUserMatch(assignee, currentUserEmail, currentUserName) {
		actions = append(actions, "Assigned")
	}

	if c.isUserMatch(reporter, currentUserEmail, currentUserName) {
		actions = append(actions, "Created")
	}

	if len(actions) == 0 {
		return "Participated"
	}

	return strings.Join(actions, ", ")
}

func (c *JiraClient) isUserMatch(nameToCheck, userEmail, userName string) bool {
	if nameToCheck == "" {
		return false
	}

	nameToCheckLower := strings.ToLower(nameToCheck)
	userEmailLower := strings.ToLower(userEmail)
	userNameLower := strings.ToLower(userName)

	if strings.Contains(nameToCheckLower, userEmailLower) || strings.Contains(userEmailLower, nameToCheckLower) {
		return true
	}

	if userName != "" && (strings.Contains(nameToCheckLower, userNameLower) || strings.Contains(userNameLower, nameToCheckLower)) {
		return true
	}

	if strings.Contains(userEmail, ".") {
		emailParts := strings.Split(strings.Split(userEmail, "@")[0], ".")
		for _, part := range emailParts {
			if len(part) > 2 && strings.Contains(nameToCheckLower, strings.ToLower(part)) {
				return true
			}
		}
	}

	return false
}
