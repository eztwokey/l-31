package models

import "time"

const (
	StatusScheduled NotifyStatus = "scheduled"
	StatusCanceled  NotifyStatus = "canceled"
	StatusSent      NotifyStatus = "sent"
	StatusFailed    NotifyStatus = "failed"
)

type NotifyStatus string

const (
	ChannelTelegram NotifyChannel = "telegram"
	ChannelLog      NotifyChannel = "log"
)

type NotifyChannel string

type Notification struct {
	ID          string        `json:"id"`
	Message     string        `json:"message"`
	Channel     NotifyChannel `json:"channel"`
	Recipient   string        `json:"recipient"`
	ScheduledAt time.Time     `json:"scheduled_at"`
	Status      NotifyStatus  `json:"status"`
	RetryCount  int           `json:"retry_count"`
	MaxRetries  int           `json:"max_retries"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type CreateNotifyRequest struct {
	Message     string `json:"message"`
	Channel     string `json:"channel"`
	Recipient   string `json:"recipient"`
	DelaySec    int64  `json:"delay_sec"`
	ScheduledAt string `json:"scheduled_at"`
}
