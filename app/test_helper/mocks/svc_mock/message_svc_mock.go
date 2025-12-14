package svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"github.com/stretchr/testify/mock"
)

type MessageSvcMock struct {
	mock.Mock
}

func (m *MessageSvcMock) SendMessage(message model.Message, ctx *atylabmongo.MongoCtxSvc) (string, error) {
	args := m.Called(message, ctx)
	return args.String(0), args.Error(1)
}
func (m *MessageSvcMock) GetMessageList(roomID string, ctx *atylabmongo.MongoCtxSvc) ([]model.Message, error) {
	args := m.Called(roomID, ctx)
	return args.Get(0).([]model.Message), args.Error(1)
}

func (m *MessageSvcMock) ReadMessages(messageIds []string, roomId string, userId string, ctx *atylabmongo.MongoCtxSvc) error {
	args := m.Called(messageIds, roomId, userId, ctx)
	return args.Error(0)
}

func (m *MessageSvcMock) IsSender(messageID string, roomID string, userID string, ctx *atylabmongo.MongoCtxSvc) error {
	args := m.Called(messageID, roomID, userID, ctx)
	return args.Error(0)
}

func (m *MessageSvcMock) DeleteMessage(messageID string, roomID string, ctx *atylabmongo.MongoCtxSvc) error {
	args := m.Called(messageID, roomID, ctx)
	return args.Error(0)
}
