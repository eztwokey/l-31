package sender

import (
	"context"

	"github.com/wb-go/wbf/logger"
)

type LogSender struct {
	logger logger.Logger
}

func NewLog(logger logger.Logger) *LogSender {
	return &LogSender{logger: logger}
}

func (l *LogSender) Send(ctx context.Context, recipient, message string) error {
	l.logger.Info("LOG SENDER: notification sent",
		"recipient", recipient,
		"message", message,
	)
	return nil
}
