package svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabapi"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
	"github.com/stretchr/testify/mock"
)

type RoomSvcMock struct {
	mock.Mock
}

func (m *RoomSvcMock) IsMember(room model.Room, uuid string) bool {
	args := m.Called(room, uuid)
	return args.Bool(0)
}

func (m *RoomSvcMock) IsOwner(room model.Room, uuid string) bool {
	args := m.Called(room, uuid)
	return args.Bool(0)
}

func (m *RoomSvcMock) GetRoom(roomId string, ctx *atylabmongo.MongoCtxSvc) (model.Room, error) {
	args := m.Called(roomId, ctx)
	return args.Get(0).(model.Room), args.Error(1)
}

func (m *RoomSvcMock) GetMemberInfos(room model.Room, ctx *atylabapi.ApiCtxSvc) ([]model.RoomMember, error) {
	args := m.Called(room, ctx)
	return args.Get(0).([]model.RoomMember), args.Error(1)
}
