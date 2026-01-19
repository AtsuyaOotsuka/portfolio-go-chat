package service

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
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
	redis        usecase.RedisUseCaseInterface
	mongoRoomSvc mongo_svc.RoomSvcInterface
	api          atylabapi.ApiPostInterface
}

func NewRoomSvc(
	redis usecase.RedisUseCaseInterface,
	mongoRoomSvc mongo_svc.RoomSvcInterface,
	api atylabapi.ApiPostInterface,
) RoomSvcInterface {
	return &RoomSvc{
		redis:        redis,
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
	rawJSON, err := s.getMemberInfos(room, ctx)
	if err != nil {
		return nil, err
	}
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

func (s *RoomSvc) getMemberInfos(room model.Room, ctx *atylabapi.ApiCtxSvc) ([]byte, error) {
	redisKey := "room:" + room.ID.Hex() + ":members"

	// redisにキャッシュがあればそちらを返す
	rawJSON, err := s.callRedisToGetMemberInfos(ctx, redisKey)
	if err == nil && len(rawJSON) > 0 {
		return rawJSON, nil
	}

	rawJSON, err = s.callApiToGetMemberInfos(room, ctx)
	if err != nil {
		return nil, err
	}

	// 取得した情報をredisにキャッシュする
	redis, err := s.redis.RedisInit()
	if err != nil {
		return nil, err
	}
	err = redis.RedisConnector.Client.Set(
		ctx.Ctx,
		redisKey,
		string(rawJSON),
		1*time.Minute, // キャッシュの有効期限を1分に設定
	)
	if err != nil {
		fmt.Println("Failed to cache member infos to Redis:", err)
		// キャッシュ登録に失敗しても情報自体は取得できていれば良いので、rawJSONを返す
		return rawJSON, nil
	}

	return rawJSON, nil
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
	fmt.Println("Fetched member infos from API:", string(rawJSON))

	return rawJSON, err
}

func (s *RoomSvc) callRedisToGetMemberInfos(ctx *atylabapi.ApiCtxSvc, redisKey string) ([]byte, error) {
	redis, err := s.redis.RedisInit()
	if err != nil {
		fmt.Println("Failed to initialize Redis:", err)
		return nil, err
	}

	cachedMembers, err := redis.RedisConnector.Client.Get(
		ctx.Ctx,
		redisKey,
	)
	fmt.Println("Fetched member infos from Redis cache:", cachedMembers)

	return []byte(cachedMembers), err
}
