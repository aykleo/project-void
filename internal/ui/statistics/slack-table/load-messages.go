package slacktable

import (
	"fmt"
	"project-void/internal/slack"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

func (m *Model) LoadMessages(slackClient *slack.SlackClient, since time.Time, config *slack.Config) error {
	m.loadingState = LoadingInProgress
	m.progress.SetPercent(0.0)

	m.progress.SetPercent(0.3)
	messages, err := slackClient.GetMessagesSince(since, config)
	if err != nil {
		m.loadingState = LoadingError
		m.loadError = err.Error()
		return fmt.Errorf("failed to load messages: %w", err)
	}

	m.progress.SetPercent(0.6)
	var userMessages []slack.Message
	for _, message := range messages {
		if config.FilterByUser && config.UserID != "" {
			if message.UserID == config.UserID {
				userMessages = append(userMessages, message)
			}
		} else {
			userMessages = append(userMessages, message)
		}
	}

	rows := make([]table.Row, len(userMessages))
	for i, message := range userMessages {
		if len(userMessages) > 0 {
			progress := 0.6 + (0.3 * float64(i) / float64(len(userMessages)))
			m.progress.SetPercent(progress)
		}

		fromName := message.User
		if len(fromName) > 11 {
			fromName = fromName[:8] + "..."
		}

		channelName, err := slackClient.GetChannelName(message.Channel)
		if err != nil {
			channelName = message.Channel
		}

		toName := channelName
		if len(toName) > 15 {
			toName = toName[:12] + "..."
		}

		timestamp := message.Timestamp.Format("15:04")

		text := strings.ReplaceAll(message.Text, "\n", " ")
		text = strings.TrimSpace(text)
		maxTextLength := 60
		if len(text) > maxTextLength {
			text = text[:maxTextLength-3] + "..."
		}

		rows[i] = table.Row{
			fromName,
			toName,
			timestamp,
			text,
		}
	}

	m.progress.SetPercent(0.9)
	m.table.SetRows(rows)
	m.progress.SetPercent(1.0)
	m.loadingState = LoadingComplete
	return nil
}
