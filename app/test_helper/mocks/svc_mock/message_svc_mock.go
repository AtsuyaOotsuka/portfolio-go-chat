package svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/stretchr/testify/mock"
)

type MessageSvcMock struct {
	mock.Mock
}

func (m *MessageSvcMock) SendMessage(message model.Message) (string, error) {
	args := m.Called(message)
	return args.String(0), args.Error(1)
}
func (m *MessageSvcMock) GetMessageList(roomID string) ([]model.Message, error) {
	args := m.Called(roomID)
	return args.Get(0).([]model.Message), args.Error(1)
}

func (m *MessageSvcMock) ReadMessages(messageIds []string, roomId string, userId string) error {
	args := m.Called(messageIds, roomId, userId)
	return args.Error(0)
}

func (m *MessageSvcMock) IsSender(messageID string, roomID string, userID string) error {
	args := m.Called(messageID, roomID, userID)
	return args.Error(0)
}

func (m *MessageSvcMock) DeleteMessage(messageID string, roomID string) error {
	args := m.Called(messageID, roomID)
	return args.Error(0)
}
