package service

import (
	"slices"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
)

type RoomSvcInterface interface {
	GetRoom(roomId string, ctx *atylabmongo.MongoCtxSvc) (model.Room, error)
	IsMember(room model.Room, uuid string) bool
	IsOwner(room model.Room, uuid string) bool
}

type RoomSvc struct {
	mongoRoomSvc mongo_svc.RoomSvcInterface
}

func NewRoomSvc(
	mongoRoomSvc mongo_svc.RoomSvcInterface,
) RoomSvcInterface {
	return &RoomSvc{
		mongoRoomSvc: mongoRoomSvc,
	}
}

func (s *RoomSvc) GetRoom(roomId string, ctx *atylabmongo.MongoCtxSvc) (model.Room, error) {
	return s.mongoRoomSvc.GetRoomByID(roomId, ctx)
}

func (s *RoomSvc) IsMember(room model.Room, uuid string) bool {
	return slices.Contains(room.Members, uuid)
}

func (s *RoomSvc) IsOwner(room model.Room, uuid string) bool {
	return room.OwnerID == uuid
}
