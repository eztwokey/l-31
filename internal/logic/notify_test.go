package logic

import (
	"testing"

	"github.com/eztwokey/l3-serv/internal/models"
)

func TestCreateNotify_EmptyMessage(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateNotify(t.Context(), models.CreateNotifyRequest{
		Message: "",
	})
	if err == nil {
		t.Fatal("expected error for empty message, got nil")
	}
}

func TestCreateNotify_WhitespaceMessage(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateNotify(t.Context(), models.CreateNotifyRequest{
		Message: "   ",
	})
	if err == nil {
		t.Fatal("expected error for whitespace message, got nil")
	}
}

func TestCreateNotify_NegativeDelay(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateNotify(t.Context(), models.CreateNotifyRequest{
		Message:  "hello",
		DelaySec: -5,
	})
	if err == nil {
		t.Fatal("expected error for negative delay, got nil")
	}
}

func TestCreateNotify_InvalidScheduledAt(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateNotify(t.Context(), models.CreateNotifyRequest{
		Message:     "hello",
		ScheduledAt: "not-a-date",
	})
	if err == nil {
		t.Fatal("expected error for invalid scheduled_at, got nil")
	}
}
