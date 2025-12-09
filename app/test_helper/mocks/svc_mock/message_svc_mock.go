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
