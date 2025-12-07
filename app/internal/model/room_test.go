package model

import (
	"testing"
	"time"
)

func TestRoomModel(t *testing.T) {
	timeNow := time.Now()

	room := Room{
		Name:      "Test Room",
		OwnerID:   "owner123",
		CreatedAt: timeNow,
		Members:   []string{"member1", "member2"},
		IsPrivate: false,
	}

	if room.Name != "Test Room" {
		t.Errorf("Expected room name to be 'Test Room', got %s", room.Name)
	}

	if room.OwnerID != "owner123" {
		t.Errorf("Expected owner ID to be 'owner123', got %s", room.OwnerID)
	}

	if !room.CreatedAt.Equal(timeNow) {
		t.Errorf("Expected CreatedAt to be %v, got %v", timeNow, room.CreatedAt)
	}

	if len(room.Members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(room.Members))
	}

	if room.IsPrivate {
		t.Errorf("Expected room to be public, but it is private")
	}

}
