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

func TestNewMessageSvcStruct(t *testing.T) {
	atylabMongo := usecase.NewMongoUseCaseStruct(atylabmongo.NewMongoConnectionStruct())
	svc := NewMessageSvcStruct(atylabMongo)
	if svc == nil {
		t.Error("expected non-nil MessageSvcStruct")
		return
	}
	assert.Equal(t, atylabMongo, svc.mongo, "expected mongo field to be set correctly")
}

func TestSendMessage(t *testing.T) {
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
				mongoDatabaseMock.On("Collection", "messages").Return(mongoCollectionMock)

				mongoConnectorStruct := &atylabmongo.MongoConnector{
					Db:     mongoDatabaseMock,
					Ctx:    context.TODO(),
					Cancel: func() {},
				}

				mongoConnectionStructMock := setupInitMock(tt.initErr, mongoConnectorStruct)
				mongoUseCase := usecase.NewMongoUseCaseStruct(mongoConnectionStructMock)
				messageSvc := NewMessageSvcStruct(mongoUseCase)

				message := model.Message{
					RoomID:  "room1",
					Sender:  "user1",
					Message: "Hello, World!",
				}

				messageID, err := messageSvc.SendMessage(message)
				if (err != nil) != tt.initErr && (err != nil) != tt.InsertOneErr {
					t.Errorf("CreateRoom() [%s] error = %v, wantErr %v", tt.name, err, tt.initErr)
				}
				if messageID != "mocked_id" && !tt.initErr && !tt.InsertOneErr {
					t.Errorf("expected messageID to be 'mocked_id', got %s", messageID)
				}
				if m, ok := mongoConnectionStructMock.(interface{ AssertExpectations(*testing.T) }); ok {
					m.AssertExpectations(t)
				}
			})
		}
	})
}

func TestGetMessageList(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		tests := []struct {
			name       string
			initErr    bool
			findOneErr bool
			decodeErr  bool
			returnErr  bool
		}{
			{"success", false, false, false, false},
			{"error", true, false, false, true},
			{"findone_error", false, true, false, true},
			{"decode_error", false, false, true, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mongoCollectionMock := new(atylabmongo.MongoCollectionStructMock)
				var message model.Message
				mongoCursorMock := new(atylabmongo.MongoCursorStructMock)
				mongoCursorMock.On("Next", mock.Anything).Return(true).Once()
				mongoCursorMock.On("Next", mock.Anything).Return(false).Once()
				if tt.decodeErr {
					mongoCursorMock.On("Decode", &message).Return(assert.AnError)
				} else {
					mongoCursorMock.On("Decode", &message).Return(nil)
				}
				mongoCursorMock.On("Close", mock.Anything).Return(nil)

				filter := bson.M{"roomid": "room1"}
				if tt.findOneErr {
					mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, assert.AnError)
				} else {
					mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, nil)
				}

				mongoDatabaseMock := new(atylabmongo.MongoDatabaseStructMock)
				mongoDatabaseMock.On("Collection", "messages").Return(mongoCollectionMock)

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
				messageSvc := NewMessageSvcStruct(mongoUseCase)

				messages, err := messageSvc.GetMessageList("room1")
				if (err != nil) != tt.returnErr {
					t.Errorf("GetMessageList() [%s] error = %v, initErr %v", tt.name, err, tt.initErr)
				}
				if len(messages) != 1 && !tt.returnErr {
					t.Errorf("expected 1 message, got %d", len(messages))
				}

				if tt.returnErr {
					return
				}

				if m, ok := mongoConnectionStructMock.(interface{ AssertExpectations(*testing.T) }); ok {
					m.AssertExpectations(t)
				}
			})
		}
	})
}
