package svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/stretchr/testify/mock"
)

type RoomSvcMock struct {
	mock.Mock
}

func (m *RoomSvcMock) GetRoomList(uuid string, target string) ([]model.Room, error) {
	args := m.Called(uuid, target)
	return args.Get(0).([]model.Room), args.Error(1)
}

func (m *RoomSvcMock) CreateRoom(room model.Room) (string, error) {
	args := m.Called(room)
	return args.String(0), args.Error(1)
}

func (m *RoomSvcMock) GetRoomByID(roomID string) (model.Room, error) {
	args := m.Called(roomID)
	return args.Get(0).(model.Room), args.Error(1)
}

func (m *RoomSvcMock) JoinRoom(roomID string, uuid string) error {
	args := m.Called(roomID, uuid)
	return args.Error(0)
}
