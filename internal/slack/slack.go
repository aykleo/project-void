package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Message struct {
	ID         string
	Channel    string
	User       string
	UserID     string
	Text       string
	Timestamp  time.Time
	UserAction string
}

type SlackClient struct {
	Token        string
	Client       *http.Client
	channelCache map[string]string
	userCache    map[string]string
}

type SlackConversationsHistoryResponse struct {
	OK       bool `json:"ok"`
	Messages []struct {
		Type string `json:"type"`
		User string `json:"user"`
		Text string `json:"text"`
		Ts   string `json:"ts"`
	} `json:"messages"`
	HasMore          bool `json:"has_more"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
	Error string `json:"error,omitempty"`
}

type SlackUsersInfoResponse struct {
	OK   bool `json:"ok"`
	User struct {
		Name        string `json:"name"`
		RealName    string `json:"real_name"`
		DisplayName string `json:"display_name"`
	} `json:"user"`
}

type SlackConversationsListResponse struct {
	OK       bool `json:"ok"`
	Channels []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channels"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
	Error string `json:"error,omitempty"`
}

func NewSlackClient(token string) *SlackClient {
	return &SlackClient{
		Token:        token,
		Client:       &http.Client{Timeout: 30 * time.Second},
		channelCache: make(map[string]string),
		userCache:    make(map[string]string),
	}
}

func (c *SlackClient) GetMessagesSince(since time.Time, config *Config) ([]Message, error) {
	var allMessages []Message

	conversations, err := c.getAllConversations()
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	for i, conversation := range conversations {
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		messages, err := c.getChannelMessages(conversation, since, config)
		if err != nil {
			continue
		}

		allMessages = append(allMessages, messages...)
	}

	return allMessages, nil
}

func (c *SlackClient) GetChannelName(channelID string) (string, error) {
	if name, exists := c.channelCache[channelID]; exists {
		return name, nil
	}

	requestURL := "https://slack.com/api/conversations.info"

	params := url.Values{}
	params.Add("channel", channelID)

	req, err := http.NewRequest("GET", requestURL+"?"+params.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var infoResponse struct {
		OK      bool `json:"ok"`
		Channel struct {
			Name   string `json:"name"`
			IsIM   bool   `json:"is_im"`
			IsMpim bool   `json:"is_mpim"`
			User   string `json:"user"`
		} `json:"channel"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&infoResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !infoResponse.OK {
		return "", fmt.Errorf("Slack API returned error")
	}

	var displayName string

	if infoResponse.Channel.IsIM {
		userName, err := c.getUserName(infoResponse.Channel.User)
		if err != nil {
			displayName = "DM:" + infoResponse.Channel.User
		} else {
			displayName = "DM:" + userName
		}
	} else if infoResponse.Channel.IsMpim {
		displayName = "Group:" + infoResponse.Channel.Name
	} else {
		displayName = "#" + infoResponse.Channel.Name
	}

	c.channelCache[channelID] = displayName
	return displayName, nil
}

func (c *SlackClient) getChannelMessages(channelID string, since time.Time, config *Config) ([]Message, error) {
	requestURL := "https://slack.com/api/conversations.history"

	oldest := fmt.Sprintf("%.6f", float64(since.UnixNano())/1e9)

	var allMessages []Message
	cursor := ""

	for {
		params := url.Values{}
		params.Add("channel", channelID)
		params.Add("oldest", oldest)
		params.Add("limit", "200")
		if cursor != "" {
			params.Add("cursor", cursor)
		}

		req, err := http.NewRequest("GET", requestURL+"?"+params.Encode(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.Token)

		resp, err := c.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
		}

		var historyResponse SlackConversationsHistoryResponse
		if err := json.NewDecoder(resp.Body).Decode(&historyResponse); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		if !historyResponse.OK {
			return nil, fmt.Errorf("Slack API returned error: %s", historyResponse.Error)
		}

		for _, slackMessage := range historyResponse.Messages {
			if slackMessage.Type != "message" || slackMessage.Text == "" {
				continue
			}

			timestamp, err := c.parseSlackTimestamp(slackMessage.Ts)
			if err != nil {
				continue
			}

			if timestamp.Before(since) {
				continue
			}

			if config.FilterByUser && config.UserID != "" && slackMessage.User != config.UserID {
				continue
			}

			userName, err := c.getUserName(slackMessage.User)
			if err != nil {
				userName = slackMessage.User
			}

			allMessages = append(allMessages, Message{
				ID:         slackMessage.Ts,
				Channel:    channelID,
				User:       userName,
				UserID:     slackMessage.User,
				Text:       slackMessage.Text,
				Timestamp:  timestamp,
				UserAction: "Sent",
			})
		}

		if !historyResponse.HasMore || historyResponse.ResponseMetadata.NextCursor == "" {
			break
		}

		cursor = historyResponse.ResponseMetadata.NextCursor
		time.Sleep(400 * time.Millisecond)
	}

	return allMessages, nil
}

func (c *SlackClient) parseSlackTimestamp(ts string) (time.Time, error) {
	parts := strings.Split(ts, ".")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid timestamp format: %s", ts)
	}

	seconds := parts[0]
	var unixTime int64
	_, err := fmt.Sscanf(seconds, "%d", &unixTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	timestamp := time.Unix(unixTime, 0)
	return timestamp, nil
}

func (c *SlackClient) getUserName(userID string) (string, error) {
	if name, ok := c.userCache[userID]; ok {
		return name, nil
	}

	requestURL := "https://slack.com/api/users.info"

	params := url.Values{}
	params.Add("user", userID)

	req, err := http.NewRequest("GET", requestURL+"?"+params.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var userResponse SlackUsersInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !userResponse.OK {
		return "", fmt.Errorf("Slack API returned error")
	}

	var bestName string
	if userResponse.User.DisplayName != "" {
		bestName = userResponse.User.DisplayName
	} else if userResponse.User.RealName != "" {
		bestName = userResponse.User.RealName
	} else {
		bestName = userResponse.User.Name
	}

	c.userCache[userID] = bestName
	return bestName, nil
}

func (c *SlackClient) TestConnection() error {
	requestURL := "https://slack.com/api/auth.test"

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed with status %d", resp.StatusCode)
	}

	var authResponse struct {
		OK bool `json:"ok"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !authResponse.OK {
		return fmt.Errorf("authentication failed")
	}

	return nil
}

func (c *SlackClient) getAllConversations() ([]string, error) {
	requestURL := "https://slack.com/api/conversations.list"

	var memberConversations []string
	var otherConversations []string
	cursor := ""

	for {
		params := url.Values{}
		params.Add("types", "public_channel,private_channel,mpim,im")
		params.Add("limit", "200")
		params.Add("exclude_archived", "true")
		if cursor != "" {
			params.Add("cursor", cursor)
		}

		req, err := http.NewRequest("GET", requestURL+"?"+params.Encode(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.Token)

		resp, err := c.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
		}

		var listResponse struct {
			OK       bool `json:"ok"`
			Channels []struct {
				ID       string `json:"id"`
				IsMember bool   `json:"is_member"`
				IsIM     bool   `json:"is_im"`
			} `json:"channels"`
			ResponseMetadata struct {
				NextCursor string `json:"next_cursor"`
			} `json:"response_metadata"`
			Error string `json:"error,omitempty"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		if !listResponse.OK {
			return nil, fmt.Errorf("Slack API returned error: %s", listResponse.Error)
		}

		for _, channel := range listResponse.Channels {
			if channel.IsMember || channel.IsIM {
				memberConversations = append(memberConversations, channel.ID)
			} else {
				otherConversations = append(otherConversations, channel.ID)
			}
		}

		if listResponse.ResponseMetadata.NextCursor == "" {
			break
		}
		cursor = listResponse.ResponseMetadata.NextCursor
		time.Sleep(400 * time.Millisecond)
	}

	conversations := append(memberConversations, otherConversations...)
	return conversations, nil
}
