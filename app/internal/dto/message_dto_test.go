package dto

import (
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetMessageInfo(t *testing.T) {

	dto := NewMessageDtoStruct()

	messageIsRead := model.Message{
		ID:            primitive.NewObjectID(),
		RoomID:        "room-uuid",
		Sender:        "sender-uuid",
		Message:       "Hello, World!",
		CreatedAt:     time.Now(),
		IsReadUserIds: []string{"reader-uuid-1", "reader-uuid-2"},
	}

	messageIsNotRead := model.Message{
		ID:            primitive.NewObjectID(),
		RoomID:        "room-uuid",
		Sender:        "sender-uuid",
		Message:       "Hello, World!",
		CreatedAt:     time.Now(),
		IsReadUserIds: []string{"reader-uuid-2", "reader-uuid-3"},
	}

	userId := "reader-uuid-1"

	response := dto.GetMessageInfo(messageIsRead, userId)

	assert.Equal(t, messageIsRead.ID.Hex(), response.ID)
	assert.Equal(t, messageIsRead.RoomID, response.RoomID)
	assert.Equal(t, messageIsRead.Sender, response.Sender)
	assert.Equal(t, messageIsRead.Message, response.Message)
	assert.Equal(t, messageIsRead.CreatedAt.String(), response.CreatedAt)
	assert.True(t, response.IsRead)

	response = dto.GetMessageInfo(messageIsNotRead, userId)

	assert.Equal(t, messageIsNotRead.ID.Hex(), response.ID)
	assert.Equal(t, messageIsNotRead.RoomID, response.RoomID)
	assert.Equal(t, messageIsNotRead.Sender, response.Sender)
	assert.Equal(t, messageIsNotRead.Message, response.Message)
	assert.Equal(t, messageIsNotRead.CreatedAt.String(), response.CreatedAt)
	assert.False(t, response.IsRead)
}

func TestResponseMessageList(t *testing.T) {
	dto := NewMessageDtoStruct()

	messages := []model.Message{
		{
			ID:            primitive.NewObjectID(),
			RoomID:        "room-uuid-1",
			Sender:        "sender-uuid-1",
			Message:       "Hello, World 1!",
			CreatedAt:     time.Now(),
			IsReadUserIds: []string{"reader-uuid-1"},
		},
		{
			ID:            primitive.NewObjectID(),
			RoomID:        "room-uuid-2",
			Sender:        "sender-uuid-2",
			Message:       "Hello, World 2!",
			CreatedAt:     time.Now(),
			IsReadUserIds: []string{"reader-uuid-2"},
		},
	}

	userId := "reader-uuid-1"
	responses := dto.ResponseMessageList(messages, userId)

	assert.Equal(t, len(messages), len(responses))

	assert.Equal(t, messages[0].ID.Hex(), responses[0].ID)
	assert.Equal(t, messages[0].RoomID, responses[0].RoomID)
	assert.Equal(t, messages[0].Sender, responses[0].Sender)
	assert.Equal(t, messages[0].Message, responses[0].Message)
	assert.Equal(t, messages[0].CreatedAt.String(), responses[0].CreatedAt)
	assert.True(t, responses[0].IsRead)

	assert.Equal(t, messages[1].ID.Hex(), responses[1].ID)
	assert.Equal(t, messages[1].RoomID, responses[1].RoomID)
	assert.Equal(t, messages[1].Sender, responses[1].Sender)
	assert.Equal(t, messages[1].Message, responses[1].Message)
	assert.Equal(t, messages[1].CreatedAt.String(), responses[1].CreatedAt)
	assert.False(t, responses[1].IsRead)
}
