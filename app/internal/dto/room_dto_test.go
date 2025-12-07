package dto

import (
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetRoomInfo(t *testing.T) {
	dto := NewRoomDtoStruct()

	room := model.Room{
		ID:        primitive.NewObjectID(),
		Name:      "Test Room",
		OwnerID:   "owner-uuid",
		IsPrivate: true,
		Members:   []string{"member-uuid-1", "member-uuid-2"},
		CreatedAt: time.Now(),
	}

	userId := "member-uuid-1"
	response := dto.GetRoomInfo(room, userId)

	assert.Equal(t, room.ID.Hex(), response.ID)
	assert.Equal(t, room.Name, response.Name)
	assert.Equal(t, room.OwnerID, response.OwnerID)
	assert.Equal(t, room.IsPrivate, response.IsPrivate)
	assert.True(t, response.IsMember)
	assert.False(t, response.IsOwner)
	assert.Equal(t, len(room.Members), response.MemberCount)
	assert.Equal(t, room.CreatedAt.String(), response.CreatedAt)
}

func TestResponseRoomList(t *testing.T) {
	dto := NewRoomDtoStruct()

	rooms := []model.Room{
		{
			ID:        primitive.NewObjectID(),
			Name:      "Room 1",
			OwnerID:   "owner-uuid-1",
			IsPrivate: false,
			Members:   []string{"member-uuid-1"},
			CreatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "Room 2",
			OwnerID:   "owner-uuid-2",
			IsPrivate: true,
			Members:   []string{"member-uuid-2", "member-uuid-3"},
			CreatedAt: time.Now(),
		},
	}

	userId := "member-uuid-2"
	responses := dto.ResponseRoomList(rooms, userId)

	assert.Len(t, responses, 2)

	assert.Equal(t, rooms[0].ID.Hex(), responses[0].ID)
	assert.Equal(t, rooms[0].Name, responses[0].Name)
	assert.Equal(t, rooms[0].OwnerID, responses[0].OwnerID)
	assert.Equal(t, rooms[0].IsPrivate, responses[0].IsPrivate)
	assert.False(t, responses[0].IsMember)
	assert.False(t, responses[0].IsOwner)
	assert.Equal(t, len(rooms[0].Members), responses[0].MemberCount)
	assert.Equal(t, rooms[0].CreatedAt.String(), responses[0].CreatedAt)

	assert.Equal(t, rooms[1].ID.Hex(), responses[1].ID)
	assert.Equal(t, rooms[1].Name, responses[1].Name)
	assert.Equal(t, rooms[1].OwnerID, responses[1].OwnerID)
	assert.Equal(t, rooms[1].IsPrivate, responses[1].IsPrivate)
	assert.True(t, responses[1].IsMember)
	assert.False(t, responses[1].IsOwner)
	assert.Equal(t, len(rooms[1].Members), responses[1].MemberCount)
	assert.Equal(t, rooms[1].CreatedAt.String(), responses[1].CreatedAt)
}
