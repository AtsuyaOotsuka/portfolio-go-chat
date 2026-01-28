package model

import (
	"testing"
	"time"
)

func TestMessageModel(t *testing.T) {
	timeNow := time.Now()

	message := Message{
		RoomID:        "room123",
		Sender:        "456",
		Message:       "Hello, World!",
		CreatedAt:     timeNow,
		IsReadUserIds: []string{"789", "101"},
	}

	if message.RoomID != "room123" {
		t.Errorf("Expected RoomID to be 'room123', got %s", message.RoomID)
	}

	if message.Sender != "456" {
		t.Errorf("Expected Sender to be 'user456', got %s", message.Sender)
	}

	if message.Message != "Hello, World!" {
		t.Errorf("Expected Message to be 'Hello, World!', got %s", message.Message)
	}

	if !message.CreatedAt.Equal(timeNow) {
		t.Errorf("Expected CreatedAt to be %v, got %v", timeNow, message.CreatedAt)
	}

	if len(message.IsReadUserIds) != 2 {
		t.Errorf("Expected 2 IsReadUserIds, got %d", len(message.IsReadUserIds))
	}
}
