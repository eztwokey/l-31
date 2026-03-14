package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/eztwokey/l3-serv/internal/models"
	"github.com/eztwokey/l3-serv/internal/storage"
)

var (
	ErrBadRequest = errors.New("bad request")
)

func (l *Logic) CreateNotify(ctx context.Context, req models.CreateNotifyRequest) (models.Notification, error) {
	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		return models.Notification{}, ErrBadRequest
	}

	var scheduled time.Time
	if req.ScheduledAt != "" {
		t, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			return models.Notification{}, ErrBadRequest
		}
		scheduled = t
	} else {
		if req.DelaySec < 0 {
			return models.Notification{}, ErrBadRequest
		}
		scheduled = time.Now().Add(time.Duration(req.DelaySec) * time.Second)
	}

	channel := models.NotifyChannel(req.Channel)
	if channel == "" {
		channel = models.ChannelLog
	}

	now := time.Now()
	n := models.Notification{
		ID:          uuid.NewString(),
		Message:     msg,
		Channel:     channel,
		Recipient:   req.Recipient,
		ScheduledAt: scheduled,
		Status:      models.StatusScheduled,
		RetryCount:  0,
		MaxRetries:  5,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	out, err := l.store.Create(ctx, n)
	if err != nil {
		l.logger.Error("notify.create failed", "err", err)
		return models.Notification{}, err
	}

	body, err := json.Marshal(out)
	if err != nil {
		l.logger.Error("notify.create marshal failed", "id", out.ID, "err", err)
		return models.Notification{}, err
	}

	if err := l.publisher.Publish(ctx, body, "notify"); err != nil {
		l.logger.Error("notify.create publish failed", "id", out.ID, "err", err)
	}

	l.logger.Info("notify created", "id", out.ID, "scheduled_at", out.ScheduledAt, "status", out.Status)
	return out, nil
}

func (l *Logic) GetNotify(ctx context.Context, id string) (models.Notification, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return models.Notification{}, ErrBadRequest
	}

	n, err := l.store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.Notification{}, storage.ErrNotFound
		}
		l.logger.Error("notify.get failed", "id", id, "err", err)
		return models.Notification{}, err
	}
	return n, nil
}

func (l *Logic) CancelNotify(ctx context.Context, id string) (models.Notification, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return models.Notification{}, ErrBadRequest
	}

	n, err := l.store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.Notification{}, storage.ErrNotFound
		}
		l.logger.Error("notify.cancel get failed", "id", id, "err", err)
		return models.Notification{}, err
	}

	if n.Status == models.StatusSent {
		return models.Notification{}, ErrBadRequest
	}

	n.Status = models.StatusCanceled
	n.UpdatedAt = time.Now()

	n, err = l.store.Update(ctx, n)
	if err != nil {
		l.logger.Error("notify.cancel update failed", "id", id, "err", err)
		return models.Notification{}, err
	}

	l.logger.Warn("notify canceled", "id", id)
	return n, nil
}
