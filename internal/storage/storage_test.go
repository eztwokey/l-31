package storage

import (
	"testing"
)

func TestNotifyKey(t *testing.T) {
	key := notifyKey("abc-123")
	expected := "notify:abc-123"

	if key != expected {
		t.Errorf("expected %q, got %q", expected, key)
	}
}

func TestNotifyKey_Empty(t *testing.T) {
	key := notifyKey("")
	expected := "notify:"

	if key != expected {
		t.Errorf("expected %q, got %q", expected, key)
	}
}
