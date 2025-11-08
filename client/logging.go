package client

import (
	"fmt"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/sirupsen/logrus"
)

// ChatLogHook is a logrus hook that sends client log messages to chat
type ChatLogHook struct {
	client *Client
}

func (h *ChatLogHook) Levels() []logrus.Level {
	// Send all log levels to chat (including Debug)
	return []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel}
}

func (h *ChatLogHook) Fire(entry *logrus.Entry) error {
	// Format the message with level prefix
	levelPrefix := ""
	switch entry.Level {
	case logrus.DebugLevel:
		levelPrefix = "DEBUG"
	case logrus.InfoLevel:
		levelPrefix = "INFO"
	case logrus.WarnLevel:
		levelPrefix = "WARN"
	case logrus.ErrorLevel:
		levelPrefix = "ERROR"
	}

	msg := entry.Message
	if len(entry.Data) > 0 {
		// Add fields to message
		for k, v := range entry.Data {
			msg += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	chatMsg := &message.ChatMessage{
		Author:  levelPrefix,
		Message: msg,
	}
	h.client.handleChatMessage(chatMsg)

	return nil
}
