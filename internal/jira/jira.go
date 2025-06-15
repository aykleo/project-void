package jira

import (
	"encoding/json"
	"fmt"
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

type JiraSearchResponse struct {
	Issues []struct {
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
			Status  struct {
				Name string `json:"name"`
			} `json:"status"`
			Assignee struct {
				DisplayName string `json:"displayName"`
			} `json:"assignee"`
			Reporter struct {
				DisplayName string `json:"displayName"`
			} `json:"reporter"`
			Updated   string `json:"updated"`
			Created   string `json:"created"`
			IssueType struct {
				Name string `json:"name"`
			} `json:"issuetype"`
			Priority struct {
				Name string `json:"name"`
			} `json:"priority"`
		} `json:"fields"`
	} `json:"issues"`
	Total int `json:"total"`
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
		switch config.UserFilterType {
		case "assignee":
			jql = fmt.Sprintf("(%s) AND assignee = currentUser()", jql)
		case "reporter":
			jql = fmt.Sprintf("(%s) AND reporter = currentUser()", jql)
		case "participant":
			jql = fmt.Sprintf("(%s) AND (assignee = currentUser() OR reporter = currentUser())", jql)
		case "all":
			jql = fmt.Sprintf("(%s) AND (assignee = currentUser() OR reporter = currentUser())", jql)
		}
	}

	jql += " ORDER BY updated DESC"

	params := url.Values{}
	params.Add("jql", jql)
	params.Add("fields", "summary,status,assignee,reporter,updated,created,issuetype,priority")
	params.Add("maxResults", "100")

	requestURL := fmt.Sprintf("%s/rest/api/2/search?%s", c.BaseURL, params.Encode())

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.ApiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var searchResponse JiraSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	issues := make([]Issue, len(searchResponse.Issues))
	for i, jiraIssue := range searchResponse.Issues {
		updated, _ := time.Parse("2006-01-02T15:04:05.000-0700", jiraIssue.Fields.Updated)
		created, _ := time.Parse("2006-01-02T15:04:05.000-0700", jiraIssue.Fields.Created)

		assignee := "Unassigned"
		if jiraIssue.Fields.Assignee.DisplayName != "" {
			assignee = jiraIssue.Fields.Assignee.DisplayName
		}

		reporter := "Unknown"
		if jiraIssue.Fields.Reporter.DisplayName != "" {
			reporter = jiraIssue.Fields.Reporter.DisplayName
		}

		userAction := c.determineUserAction(config, assignee, reporter, jiraIssue.Key)

		actionDate := updated
		if strings.Contains(userAction, "Created") {
			actionDate = created
		}

		issues[i] = Issue{
			Key:        jiraIssue.Key,
			Status:     jiraIssue.Fields.Status.Name,
			Assignee:   assignee,
			Summary:    jiraIssue.Fields.Summary,
			Updated:    updated,
			Created:    created,
			IssueType:  jiraIssue.Fields.IssueType.Name,
			Priority:   jiraIssue.Fields.Priority.Name,
			UserAction: userAction,
			ActionDate: actionDate,
		}
	}

	return issues, nil
}

func (c *JiraClient) determineUserAction(config *Config, assignee, reporter, issueKey string) string {
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
		switch config.UserFilterType {
		case "assignee":
			return "Assigned"
		case "reporter":
			return "Created"
		case "participant":
			return "Participated"
		case "all":
			return "Involved"
		default:
			return "Related"
		}
	}

	if len(actions) == 1 {
		return actions[0]
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
