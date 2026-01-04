package service

import (
	"encoding/json"
	"slices"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabapi"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
)

type RoomSvcInterface interface {
	GetRoom(roomId string, ctx *atylabmongo.MongoCtxSvc) (model.Room, error)
	IsMember(room model.Room, uuid string) bool
	IsOwner(room model.Room, uuid string) bool
	GetMemberInfos(room model.Room, ctx *atylabapi.ApiCtxSvc) ([]model.RoomMember, error)
}

type RoomSvc struct {
	mongoRoomSvc mongo_svc.RoomSvcInterface
	api          atylabapi.ApiPostInterface
}

func NewRoomSvc(
	mongoRoomSvc mongo_svc.RoomSvcInterface,
	api atylabapi.ApiPostInterface,
) RoomSvcInterface {
	return &RoomSvc{
		mongoRoomSvc: mongoRoomSvc,
		api:          api,
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

func (s *RoomSvc) GetMemberInfos(room model.Room, ctx *atylabapi.ApiCtxSvc) ([]model.RoomMember, error) {
	// redisにキャッシュがあればそちらを返す

	// なければapiをコールして取得
	rawJSON, err := s.callApiToGetMemberInfos(room, ctx)
	if err != nil {
		return nil, err
	}

	// 取得したらredisにキャッシュする

	// 最後に取得した情報を返す

	var members []model.RoomMember
	for _, v := range room.Members {
		member := model.RoomMember{
			Uuid: v,
		}
		var apiMembers []map[string]string
		if err := json.Unmarshal(rawJSON, &apiMembers); err != nil {
			return nil, err
		}

		for _, apiMember := range apiMembers {
			if apiMember["uuid"] == v {
				member.Name = apiMember["username"]
				member.Email = apiMember["email"]
				break
			}
		}
		members = append(members, member)
	}

	return members, nil
}

func (s *RoomSvc) callApiToGetMemberInfos(room model.Room, ctx *atylabapi.ApiCtxSvc) ([]byte, error) {
	uuids := map[string][]string{
		"uuids": room.Members,
	}

	rawJSON, err := s.api.Post(
		"/server_api/user/profile",
		uuids,
		ctx,
	)

	return rawJSON, err
}
