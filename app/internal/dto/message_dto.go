package dto

import "github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"

type MessageDtoInterface interface {
	GetMessageInfo(message model.Message, userId string) MessageResponse
	ResponseMessageList(messages []model.Message, uuid string) []MessageResponse
}

type MessageDtoStruct struct{}

func NewMessageDtoStruct() *MessageDtoStruct {
	return &MessageDtoStruct{}
}

type MessageResponse struct {
	ID        string `json:"ID"`
	RoomID    string `json:"RoomID"`
	Sender    string `json:"Sender"`
	Message   string `json:"Message"`
	CreatedAt string `json:"CreatedAt"`
	IsRead    bool   `json:"IsRead"`
}

func (d *MessageDtoStruct) GetMessageInfo(message model.Message, userId string) MessageResponse {
	isRead := false
	for _, id := range message.IsReadUserIds {
		if id == userId {
			isRead = true
			break
		}
	}

	return MessageResponse{
		ID:        message.ID.Hex(),
		RoomID:    message.RoomID,
		Sender:    message.Sender,
		Message:   message.Message,
		CreatedAt: message.CreatedAt.String(),
		IsRead:    isRead,
	}
}

func (d *MessageDtoStruct) ResponseMessageList(messages []model.Message, uuid string) []MessageResponse {
	var responses []MessageResponse
	for _, msg := range messages {
		responses = append(responses, d.GetMessageInfo(msg, uuid))
	}
	return responses
}
