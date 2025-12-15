package service

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/svc_mock/mongo_svc_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRoom(t *testing.T) {
	room := model.Room{}
	mongo_svc_mock := new(mongo_svc_mock.RoomSvcMock)

	mongo_svc_mock.On("GetRoomByID", "roomId", mock.Anything).Return(room, nil)

	roomSvc := NewRoomSvc(mongo_svc_mock)

	got, err := roomSvc.GetRoom("roomId", nil)
	assert.NoError(t, err)
	assert.Equal(t, room, got)
}

func TestIsMember(t *testing.T) {
	mongo_svc_mock := new(mongo_svc_mock.RoomSvcMock)

	roomSvc := NewRoomSvc(mongo_svc_mock)

	room := model.Room{
		Members: []string{"uuid1", "uuid2"},
	}

	assert.True(t, roomSvc.IsMember(room, "uuid1"))
	assert.False(t, roomSvc.IsMember(room, "uuid3"))
}

func TestIsOwner(t *testing.T) {
	mongo_svc_mock := new(mongo_svc_mock.RoomSvcMock)

	roomSvc := NewRoomSvc(mongo_svc_mock)

	room := model.Room{
		OwnerID: "ownerUuid",
	}

	assert.True(t, roomSvc.IsOwner(room, "ownerUuid"))
	assert.False(t, roomSvc.IsOwner(room, "otherUuid"))
}
