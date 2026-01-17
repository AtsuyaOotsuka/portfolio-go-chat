package service

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabredis"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/svc_mock/mongo_svc_mock"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/usecase_mock"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetRoom(t *testing.T) {
	room := model.Room{}
	mongo_svc_mock := new(mongo_svc_mock.RoomSvcMock)
	api_mock := new(atylabapi.ApiPostStructMock)
	redis := &usecase.RedisUseCaseStruct{}

	mongo_svc_mock.On("GetRoomByID", "roomId", mock.Anything).Return(room, nil)

	roomSvc := NewRoomSvc(redis, mongo_svc_mock, api_mock)

	got, err := roomSvc.GetRoom("roomId", nil)
	assert.NoError(t, err)
	assert.Equal(t, room, got)
}

func TestIsMember(t *testing.T) {
	mongo_svc_mock := new(mongo_svc_mock.RoomSvcMock)
	api_mock := new(atylabapi.ApiPostStructMock)
	redis := &usecase.RedisUseCaseStruct{}

	roomSvc := NewRoomSvc(redis, mongo_svc_mock, api_mock)

	room := model.Room{
		Members: []string{"uuid1", "uuid2"},
	}

	assert.True(t, roomSvc.IsMember(room, "uuid1"))
	assert.False(t, roomSvc.IsMember(room, "uuid3"))
}

func TestIsOwner(t *testing.T) {
	mongo_svc_mock := new(mongo_svc_mock.RoomSvcMock)
	api_mock := new(atylabapi.ApiPostStructMock)
	redis := &usecase.RedisUseCaseStruct{}

	roomSvc := NewRoomSvc(redis, mongo_svc_mock, api_mock)

	room := model.Room{
		OwnerID: "ownerUuid",
	}

	assert.True(t, roomSvc.IsOwner(room, "ownerUuid"))
	assert.False(t, roomSvc.IsOwner(room, "otherUuid"))
}

func createResultData(returnDataType string) []byte {
	switch returnDataType {
	case "all_hit":
		return []byte(`[
			{"uuid": "uuid1", "name": "User One"},
			{"uuid": "uuid2", "name": "User Two"}
		]`)
	default:
		return []byte(`[]`)
	}
}

func TestGetMemberInfos(t *testing.T) {
	expected := map[string]map[string]any{
		"api_response": {
			"api_call":         true,
			"api_error":        false,
			"redis_init_error": false,
			"redis_hit":        false,
			"redis_set":        true,
			"redis_get_error":  false,
			"redis_set_error":  false,
			"return_type":      createResultData("all_hit"),
			"success":          true,
		},
		"redis_cache_hit": {
			"api_call":         false,
			"api_error":        false,
			"redis_init_error": false,
			"redis_hit":        true,
			"redis_set":        false,
			"redis_get_error":  false,
			"redis_set_error":  false,
			"return_type":      createResultData("all_hit"),
			"success":          true,
		},
		"api_response_error": {
			"api_call":         true,
			"api_error":        true,
			"redis_init_error": false,
			"redis_hit":        false,
			"redis_set":        false,
			"redis_get_error":  false,
			"redis_set_error":  false,
			"return_type":      []byte{},
			"success":          false,
		},
		"redis_get_error": {
			"api_call":         true,
			"api_error":        false,
			"redis_init_error": false,
			"redis_hit":        false,
			"redis_set":        true,
			"redis_get_error":  true,
			"redis_set_error":  false,
			"return_type":      createResultData("all_hit"),
			"success":          true,
		},
		"redis_set_error": {
			"api_call":         true,
			"api_error":        false,
			"redis_init_error": false,
			"redis_hit":        false,
			"redis_set":        true,
			"redis_get_error":  false,
			"redis_set_error":  true,
			"return_type":      createResultData("all_hit"),
			"success":          true,
		},
		"redis_init_error_for_get": {
			"api_call":         true,
			"api_error":        false,
			"redis_init_error": true,
			"redis_hit":        false,
			"redis_set":        false,
			"redis_get_error":  false,
			"redis_set_error":  false,
			"return_type":      createResultData("all_hit"),
			"success":          false,
		},
		"redis_init_error_for_set": {
			"api_call":         true,
			"api_error":        false,
			"redis_init_error": false,
			"redis_hit":        false,
			"redis_set":        true,
			"redis_get_error":  false,
			"redis_set_error":  false,
			"return_type":      createResultData("all_hit"),
			"success":          true,
		},
		"broken_api_response": {
			"api_call":         true,
			"api_error":        false,
			"redis_init_error": false,
			"redis_hit":        false,
			"redis_set":        true,
			"redis_get_error":  false,
			"redis_set_error":  false,
			"return_type":      []byte(`invalid json`),
			"success":          false,
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {
			mongo_svc_mock := new(mongo_svc_mock.RoomSvcMock)
			api_mock := new(atylabapi.ApiPostStructMock)

			if expect["api_call"].(bool) {
				var api_error error
				if expect["api_error"].(bool) {
					api_error = assert.AnError
				}

				api_mock.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(
					expect["return_type"].([]byte), api_error,
				)
			}

			room := model.Room{
				ID:      primitive.NewObjectID(),
				Members: []string{"uuid1", "uuid2"},
			}

			redisClient := new(atylabredis.RedisClientStructMock)

			var redisHitData string
			if expect["redis_hit"].(bool) {
				redisHitData = string(expect["return_type"].([]byte))
			}
			if expect["redis_get_error"].(bool) {
				redisHitData = ""
			}

			var redisGetError error
			if expect["redis_get_error"].(bool) {
				redisGetError = assert.AnError
			}

			redisClient.On("Get", mock.Anything, mock.Anything).Return(redisHitData, redisGetError)

			if expect["redis_set"].(bool) {
				var redisSetError error
				if expect["redis_set_error"].(bool) {
					redisSetError = assert.AnError
				}

				redisClient.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redisSetError)
			}

			redis := new(usecase_mock.RedisUseCaseMock)

			var redisInitError error
			if expect["redis_init_error"].(bool) {
				redisInitError = assert.AnError
			}
			redis.On("RedisInit").Return(&usecase.Redis{
				RedisConnector: &atylabredis.RedisConnector{
					Client: redisClient,
				},
				IsConnected: true,
			}, redisInitError)

			roomSvc := NewRoomSvc(redis, mongo_svc_mock, api_mock)

			ctx := atylabapi.NewApiCtxSvc()

			members, err := roomSvc.GetMemberInfos(room, ctx)
			if expect["success"].(bool) {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				return
			}
			assert.Len(t, members, 2)
		})
	}
}
