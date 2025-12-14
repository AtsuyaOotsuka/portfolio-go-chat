package mongo_svc

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewMessageSvcStruct(t *testing.T) {
	atylabMongo := usecase.NewMongoUseCaseStruct(atylabmongo.NewMongoConnectionStruct(), usecase.NewMongo())
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
					Db: mongoDatabaseMock,
				}

				mongoConnectionStructMock := setupInitMock(tt.initErr, mongoConnectorStruct)
				mongoUseCase := usecase.NewMongoUseCaseStruct(mongoConnectionStructMock, usecase.NewMongo())
				messageSvc := NewMessageSvcStruct(mongoUseCase)

				message := model.Message{
					RoomID:  "room1",
					Sender:  "user1",
					Message: "Hello, World!",
				}

				messageID, err := messageSvc.SendMessage(message, atylabmongo.NewMongoCtxSvc())
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
					Db: mongoDatabaseMock,
				}
				if tt.initErr {
					mongoConnectorStruct = nil
				}

				mongoConnectionStructMock := setupInitMock(tt.initErr, mongoConnectorStruct)
				mongoUseCase := usecase.NewMongoUseCaseStruct(mongoConnectionStructMock, usecase.NewMongo())
				messageSvc := NewMessageSvcStruct(mongoUseCase)

				messages, err := messageSvc.GetMessageList("room1", atylabmongo.NewMongoCtxSvc())
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

func TestReadMessages(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		tests := []struct {
			name                string
			id                  string
			initErr             bool
			updateManyCallCount int
			updateManyErr       bool
			returnErr           bool
		}{
			{"success", "60c72b2f9b1d4c3d88f0e6b1", false, 1, false, false},
			{"ObjectIDFromHex_error", "invalid_id", false, 0, false, true},
			{"initErr", "60c72b2f9b1d4c3d88f0e6b1", true, 0, false, true},
			{"updateone_error", "60c72b2f9b1d4c3d88f0e6b1", false, 1, true, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mongoCollectionMock := new(atylabmongo.MongoCollectionStructMock)
				mongoDatabaseMock := new(atylabmongo.MongoDatabaseStructMock)

				if tt.updateManyCallCount > 0 {
					if tt.updateManyErr {
						mongoCollectionMock.On("UpdateMany", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, assert.AnError)
					} else {
						mongoCollectionMock.On("UpdateMany", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)
					}
				}

				mongoDatabaseMock.On("Collection", "messages").Return(mongoCollectionMock)

				mongoConnectorStruct := &atylabmongo.MongoConnector{
					Db: mongoDatabaseMock,
				}

				mongoConnectionStructMock := setupInitMock(tt.initErr, mongoConnectorStruct)
				mongoUseCase := usecase.NewMongoUseCaseStruct(mongoConnectionStructMock, usecase.NewMongo())
				messageSvc := NewMessageSvcStruct(mongoUseCase)

				messageIds := []string{"60c72b2f9b1d4c3d88f0e6b1", tt.id}
				roomId := "room1"
				userId := "user1"

				err := messageSvc.ReadMessages(messageIds, roomId, userId, atylabmongo.NewMongoCtxSvc())
				if (err != nil) != tt.returnErr {
					t.Errorf("ReadMessages() [%s] error = %v, wantErr %v", tt.name, err, tt.returnErr)
				}

				if m, ok := mongoConnectionStructMock.(interface{ AssertExpectations(*testing.T) }); ok {
					m.AssertExpectations(t)
				}
			})
		}
	})
}

func TestIsSender(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		tests := []struct {
			name       string
			messageID  string
			roomID     string
			userID     string
			initErr    bool
			findOneErr bool
			returnErr  bool
		}{
			{"success", "60c72b2f9b1d4c3d88f0e6b1", "room1", "user1", false, false, false},
			{"ObjectIDFromHex_error", "invalid_id", "room1", "user1", false, false, true},
			{"initErr", "60c72b2f9b1d4c3d88f0e6b1", "room1", "user1", true, false, true},
			{"findone_error", "60c72b2f9b1d4c3d88f0e6b1", "room1", "user1", false, true, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mongoCollectionMock := new(atylabmongo.MongoCollectionStructMock)

				var message model.Message
				if tt.findOneErr {
					mongoCollectionMock.On("FindOne", mock.Anything, mock.Anything, &message).Return(assert.AnError)
				} else {
					mongoCollectionMock.On("FindOne", mock.Anything, mock.Anything, &message).Return(nil)
				}

				mongoDatabaseMock := new(atylabmongo.MongoDatabaseStructMock)
				mongoDatabaseMock.On("Collection", "messages").Return(mongoCollectionMock)

				mongoConnectorStruct := &atylabmongo.MongoConnector{
					Db: mongoDatabaseMock,
				}

				mongoConnectionStructMock := setupInitMock(tt.initErr, mongoConnectorStruct)
				mongoUseCase := usecase.NewMongoUseCaseStruct(mongoConnectionStructMock, usecase.NewMongo())
				messageSvc := NewMessageSvcStruct(mongoUseCase)

				err := messageSvc.IsSender(tt.messageID, tt.roomID, tt.userID, atylabmongo.NewMongoCtxSvc())
				if (err != nil) != tt.returnErr {
					t.Errorf("IsSender() [%s] error = %v, wantErr %v", tt.name, err, tt.returnErr)
				}

				if m, ok := mongoConnectionStructMock.(interface{ AssertExpectations(*testing.T) }); ok {
					m.AssertExpectations(t)
				}
			})
		}
	})
}

func TestDeleteMessage(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		tests := []struct {
			name         string
			messageID    string
			roomID       string
			initErr      bool
			deleteOneErr bool
			returnErr    bool
		}{
			{"success", "60c72b2f9b1d4c3d88f0e6b1", "room1", false, false, false},
			{"ObjectIDFromHex_error", "invalid_id", "room1", false, false, true},
			{"initErr", "60c72b2f9b1d4c3d88f0e6b1", "room1", true, false, true},
			{"deleteone_error", "60c72b2f9b1d4c3d88f0e6b1", "room1", false, true, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mongoCollectionMock := new(atylabmongo.MongoCollectionStructMock)

				if tt.deleteOneErr {
					mongoCollectionMock.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, assert.AnError)
				} else {
					mongoCollectionMock.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, nil)
				}

				mongoDatabaseMock := new(atylabmongo.MongoDatabaseStructMock)
				mongoDatabaseMock.On("Collection", "messages").Return(mongoCollectionMock)

				mongoConnectorStruct := &atylabmongo.MongoConnector{
					Db: mongoDatabaseMock,
				}

				mongoConnectionStructMock := setupInitMock(tt.initErr, mongoConnectorStruct)
				mongoUseCase := usecase.NewMongoUseCaseStruct(mongoConnectionStructMock, usecase.NewMongo())
				messageSvc := NewMessageSvcStruct(mongoUseCase)

				err := messageSvc.DeleteMessage(tt.messageID, tt.roomID, atylabmongo.NewMongoCtxSvc())
				if (err != nil) != tt.returnErr {
					t.Errorf("DeleteMessage() [%s] error = %v, wantErr %v", tt.name, err, tt.returnErr)
				}

				if m, ok := mongoConnectionStructMock.(interface{ AssertExpectations(*testing.T) }); ok {
					m.AssertExpectations(t)
				}
			})
		}
	})
}
