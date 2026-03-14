package models

import (
	"encoding/json"
	"testing"
)

func TestCreateNotifyRequest_JSON(t *testing.T) {
	input := `{"message":"hello","channel":"telegram","recipient":"123","delay_sec":10}`

	var req CreateNotifyRequest
	if err := json.Unmarshal([]byte(input), &req); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if req.Message != "hello" {
		t.Errorf("expected message 'hello', got %q", req.Message)
	}
	if req.Channel != "telegram" {
		t.Errorf("expected channel 'telegram', got %q", req.Channel)
	}
	if req.Recipient != "123" {
		t.Errorf("expected recipient '123', got %q", req.Recipient)
	}
	if req.DelaySec != 10 {
		t.Errorf("expected delay_sec 10, got %d", req.DelaySec)
	}
}

func TestNotification_JSON(t *testing.T) {
	n := Notification{
		ID:      "abc-123",
		Message: "test",
		Channel: ChannelTelegram,
		Status:  StatusScheduled,
	}

	data, err := json.Marshal(n)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if result["channel"] != "telegram" {
		t.Errorf("expected channel 'telegram', got %v", result["channel"])
	}
	if result["status"] != "scheduled" {
		t.Errorf("expected status 'scheduled', got %v", result["status"])
	}
}
