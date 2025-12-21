package cmd_svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"github.com/stretchr/testify/mock"
)

type RoomSvcMock struct {
	mock.Mock
}

func (m *RoomSvcMock) ListRooms(ctx *atylabmongo.MongoCtxSvc) ([]model.Room, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Room), args.Error(1)
}
