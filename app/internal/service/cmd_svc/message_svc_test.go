package cmd_svc

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/consts"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

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

func TestContainsForbiddenWords(t *testing.T) {
	defaultConstsForbiddenWords := consts.ForbiddenWords
	defer func() {
		consts.ForbiddenWords = defaultConstsForbiddenWords
	}()
	consts.ForbiddenWords = []string{"forbidden"}

	messageSvc := new(MessageSvcStruct)

	assert.True(t, messageSvc.ContainsForbiddenWords("This message contains a forbidden word."))
	assert.False(t, messageSvc.ContainsForbiddenWords("This message is clean."))
}
