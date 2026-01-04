package mongo_svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
	"github.com/stretchr/testify/mock"
)

type RoomSvcMock struct {
	mock.Mock
}

func (m *RoomSvcMock) GetRoomList(uuid string, target string, ctx *atylabmongo.MongoCtxSvc) ([]model.Room, error) {
	args := m.Called(uuid, target, ctx)
	return args.Get(0).([]model.Room), args.Error(1)
}

func (m *RoomSvcMock) CreateRoom(room model.Room, ctx *atylabmongo.MongoCtxSvc) (string, error) {
	args := m.Called(room, ctx)
	return args.String(0), args.Error(1)
}

func (m *RoomSvcMock) GetRoomByID(roomID string, ctx *atylabmongo.MongoCtxSvc) (model.Room, error) {
	args := m.Called(roomID, ctx)
	return args.Get(0).(model.Room), args.Error(1)
}

func (m *RoomSvcMock) JoinRoom(roomID string, uuid string, ctx *atylabmongo.MongoCtxSvc) error {
	args := m.Called(roomID, uuid, ctx)
	return args.Error(0)
}

func (m *RoomSvcMock) IsJoinedRoom(roomID string, uuid string, ctx *atylabmongo.MongoCtxSvc) error {
	args := m.Called(roomID, uuid, ctx)
	return args.Error(0)
}

func (m *RoomSvcMock) IsRoomOwner(roomID string, uuid string, ctx *atylabmongo.MongoCtxSvc) error {
	args := m.Called(roomID, uuid, ctx)
	return args.Error(0)
}

func (m *RoomSvcMock) LeaveRoom(roomID string, uuid string, ctx *atylabmongo.MongoCtxSvc) error {
	args := m.Called(roomID, uuid, ctx)
	return args.Error(0)
}

func (m *RoomSvcMock) DeleteRoom(roomID string, ctx *atylabmongo.MongoCtxSvc) error {
	args := m.Called(roomID, ctx)
	return args.Error(0)
}
