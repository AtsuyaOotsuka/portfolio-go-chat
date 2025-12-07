package mongo_svc

import (
	"context"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewRoomSvcStruct(t *testing.T) {
	atylabMongo := usecase.NewMongoUseCaseStruct(atylabmongo.NewMongoConnectionStruct())
	svc := NewRoomSvcStruct(atylabMongo)
	if svc == nil {
		t.Error("expected non-nil RoomSvcStruct")
		return
	}
	if svc.mongo != atylabMongo {
		t.Error("expected mongo field to be set correctly")
	}
}

var mongoSvcEnvs = funcs.Envs{
	"MONGO_HOST": "localhost",
	"MONGO_PORT": "27017",
	"MONGO_USER": "testuser",
	"MONGO_PASS": "testpass",
}

func TestGetRoomsFor(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {

		tests := []struct {
			name       string
			initErr    bool
			request    string
			findOneErr bool
			decodeErr  bool
			returnErr  bool
		}{
			{"success_all", false, "all", false, false, false},
			{"success_joined", false, "joined", false, false, false},
			{"error", true, "all", false, false, true},
			{"invalid_target", false, "invalid_target", false, false, true},
			{"findone_error", false, "all", true, false, true},
			{"decode_error", false, "all", false, true, true},
		}

		for _, tt := range tests {
			mongoCollectionMock := new(atylabmongo.MongoCollectionStructMock)
			var room model.Room
			mongoCursorMock := new(atylabmongo.MongoCursorStructMock)
			mongoCursorMock.On("Next", mock.Anything).Return(true).Once()
			mongoCursorMock.On("Next", mock.Anything).Return(false).Once()
			if tt.decodeErr {
				mongoCursorMock.On("Decode", &room).Return(assert.AnError)
			} else {
				mongoCursorMock.On("Decode", &room).Return(nil)
			}
			mongoCursorMock.On("Close", mock.Anything).Return(nil)

			var filter bson.M
			if tt.request == "all" {
				filter = bson.M{
					"$or": []bson.M{
						{"isprivate": false},
						{"members": "user123"},
					},
				}
			} else {
				filter = bson.M{"members": "user123"} // 参加済みのものだけ
			}

			if tt.findOneErr {
				mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, assert.AnError)
			} else {
				mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, nil)
			}

			mongoDatabaseMock := new(atylabmongo.MongoDatabaseStructMock)
			mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)

			mongoConnectorStruct := &atylabmongo.MongoConnector{
				Db:     mongoDatabaseMock,
				Ctx:    context.TODO(),
				Cancel: func() {},
			}
			if tt.initErr {
				mongoConnectorStruct = nil
			}

			mongoConnectionStructMock := setupInitMock(tt.initErr, mongoConnectorStruct)

			mongoUseCase := usecase.NewMongoUseCaseStruct(mongoConnectionStructMock)

			roomSvc := NewRoomSvcStruct(mongoUseCase)

			rooms, err := roomSvc.GetRoomList("user123", tt.request)
			if (err != nil) != tt.returnErr {
				t.Errorf("GetRooms() [%s] error = %v, initErr %v", tt.name, err, tt.initErr)
			}
			if len(rooms) != 1 && !tt.returnErr {
				t.Errorf("expected 1 room, got %d", len(rooms))
			}

			if tt.returnErr {
				continue
			}

			if m, ok := mongoConnectionStructMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		}
	})
}

func TestCreateRoom(t *testing.T) {

	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		tests := []struct {
			name         string
			initErr      bool
			InsertOneErr bool
		}{
			{"success", false, false},
			{"error", true, false},
			{"insert_error", false, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mongoCollectionMock := new(atylabmongo.MongoCollectionStructMock)
				if tt.InsertOneErr {
					mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("", assert.AnError)
				} else {
					mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("mocked_id", nil)
				}
				mongoDatabaseMock := new(atylabmongo.MongoDatabaseStructMock)
				mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)

				mongoConnectorStruct := &atylabmongo.MongoConnector{
					Db:     mongoDatabaseMock,
					Ctx:    context.TODO(),
					Cancel: func() {},
				}
				mongoConnectionStructMock := setupInitMock(tt.initErr, mongoConnectorStruct)

				mongoUseCase := usecase.NewMongoUseCaseStruct(mongoConnectionStructMock)
				roomSvc := NewRoomSvcStruct(mongoUseCase)

				room := model.Room{Name: "Test Room"}

				roomId, err := roomSvc.CreateRoom(room)
				if (err != nil) != tt.initErr && (err != nil) != tt.InsertOneErr {
					t.Errorf("CreateRoom() [%s] error = %v, wantErr %v", tt.name, err, tt.initErr)
				}
				if roomId != "mocked_id" && !tt.initErr && !tt.InsertOneErr {
					t.Errorf("expected roomId to be 'mocked_id', got %s", roomId)
				}

				if m, ok := mongoConnectionStructMock.(interface{ AssertExpectations(*testing.T) }); ok {
					m.AssertExpectations(t)
				}
			})
		}
	})
}

func setupInitMock(wantErr bool, mongoConnectorStruct *atylabmongo.MongoConnector) atylabmongo.MongoConnectorInterface {
	var returnErr error
	if wantErr {
		returnErr = assert.AnError
	}
	mongoConnectionStructMock := new(atylabmongo.MongoConnectionStructMock)
	mongoConnectionStructMock.On("NewMongoConnect", "chatapp", mock.Anything).Return(mongoConnectorStruct, returnErr)
	return mongoConnectionStructMock
}
