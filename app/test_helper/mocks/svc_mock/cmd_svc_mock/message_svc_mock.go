package cmd_svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
	"github.com/stretchr/testify/mock"
)

type MessageSvcMock struct {
	mock.Mock
}

func (m *MessageSvcMock) GetMessageList(roomID string, ctx *atylabmongo.MongoCtxSvc) ([]model.Message, error) {
	args := m.Called(roomID, ctx)
	return args.Get(0).([]model.Message), args.Error(1)
}

func (m *MessageSvcMock) ContainsForbiddenWords(message string) bool {
	args := m.Called(message)
	return args.Bool(0)
}
