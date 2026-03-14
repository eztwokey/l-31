package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/logger"

	"github.com/eztwokey/l3-serv/internal/models"
	"github.com/eztwokey/l3-serv/internal/sender"
	"github.com/eztwokey/l3-serv/internal/storage"
)

type Worker struct {
	store   *storage.Storage
	senders map[models.NotifyChannel]sender.Sender
	logger  logger.Logger
}

func New(store *storage.Storage, senders map[models.NotifyChannel]sender.Sender, logger logger.Logger) *Worker {
	return &Worker{
		store:   store,
		senders: senders,
		logger:  logger,
	}
}

func (w *Worker) Handle(ctx context.Context, msg amqp091.Delivery) error {
	var n models.Notification
	if err := json.Unmarshal(msg.Body, &n); err != nil {
		w.logger.Error("worker: unmarshal failed, skipping message", "err", err)
		return nil
	}

	delay := time.Until(n.ScheduledAt)
	if delay > 0 {
		w.logger.Info("worker: waiting for scheduled time",
			"id", n.ID,
			"delay", delay.String(),
		)

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}

	fresh, err := w.store.Get(ctx, n.ID)
	if err != nil {
		w.logger.Error("worker: failed to get notification from Redis",
			"id", n.ID,
			"err", err,
		)
		return fmt.Errorf("get from redis: %w", err)
	}

	if fresh.Status == models.StatusCanceled {
		w.logger.Info("worker: notification was canceled, skipping", "id", n.ID)
		return nil
	}

	if fresh.Status == models.StatusSent {
		w.logger.Info("worker: notification already sent, skipping", "id", n.ID)
		return nil

	}

	if fresh.RetryCount >= fresh.MaxRetries {
		fresh.Status = models.StatusFailed
		fresh.UpdatedAt = time.Now()

		if _, err := w.store.Update(ctx, fresh); err != nil {
			w.logger.Error("worker: failed to update retry count", "id", fresh.ID, "err", err)
		}

		w.logger.Error("worker: max retries reached, marking as failed", "id", n.ID)
		return nil
	}

	s, ok := w.senders[fresh.Channel]
	if !ok {
		s, ok = w.senders[models.ChannelLog]
		if !ok {
			w.logger.Error("worker: no sender for channel", "channel", fresh.Channel, "id", n.ID)
			return fmt.Errorf("no sender for channel %s", fresh.Channel)
		}
	}

	if err := s.Send(ctx, fresh.Recipient, fresh.Message); err != nil {
		w.logger.Error("worker: send failed",
			"id", n.ID,
			"channel", fresh.Channel,
			"retry", fresh.RetryCount,
			"err", err,
		)

		fresh.RetryCount++
		fresh.UpdatedAt = time.Now()
		if _, err := w.store.Update(ctx, fresh); err != nil {
			w.logger.Error("worker: failed to update retry count", "id", fresh.ID, "err", err)
		}

		return fmt.Errorf("send failed: %w", err)
	}

	fresh.Status = models.StatusSent
	fresh.UpdatedAt = time.Now()

	if _, err := w.store.Update(ctx, fresh); err != nil {
		w.logger.Error("worker: failed to update status to sent", "id", n.ID, "err", err)
	}

	w.logger.Info("worker: notification sent successfully",
		"id", n.ID,
		"channel", fresh.Channel,
		"recipient", fresh.Recipient,
	)

	return nil
}
