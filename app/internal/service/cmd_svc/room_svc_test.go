package cmd_svc

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListRooms(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		tests := []struct {
			name       string
			initErr    bool
			returnErr  bool
			allCalled  bool
			allErr     bool
			findCalled bool
			findErr    bool
		}{
			{name: "Successful retrieval of rooms",
				initErr:    false,
				returnErr:  false,
				allCalled:  true,
				allErr:     false,
				findCalled: true,
				findErr:    false,
			},
			{name: "MongoDB initialization failure",
				initErr:    true,
				returnErr:  true,
				allCalled:  false,
				allErr:     false,
				findCalled: false,
				findErr:    false,
			},
			{name: "MongoDB Find failure",
				initErr:    false,
				returnErr:  true,
				allCalled:  false,
				allErr:     false,
				findCalled: true,
				findErr:    true,
			},
			{name: "MongoDB All failure",
				initErr:    false,
				returnErr:  true,
				allCalled:  true,
				allErr:     true,
				findCalled: true,
				findErr:    false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mongoCollectionMock := new(atylabmongo.MongoCollectionStructMock)
				var rooms []model.Room

				mongoCursorMock := new(atylabmongo.MongoCursorStructMock)
				if tt.allCalled {
					if tt.allErr {
						mongoCursorMock.On("All", mock.Anything, &rooms).Return(assert.AnError)
					} else {
						mongoCursorMock.On("All", mock.Anything, &rooms).Run(func(args mock.Arguments) {
							roomsPtr := args.Get(1).(*[]model.Room)
							*roomsPtr = []model.Room{
								{Name: "General"},
								{Name: "Random"},
							}
						}).Return(nil)
					}
				}

				if tt.findCalled {
					if tt.findErr {
						mongoCollectionMock.On("Find", mock.Anything, mock.Anything).Return(mongoCursorMock, assert.AnError)
					} else {
						mongoCollectionMock.On("Find", mock.Anything, mock.Anything).Return(mongoCursorMock, nil)
						mongoCursorMock.On("Close", mock.Anything).Return(nil)
					}
				}

				mongoDatabaseMock := new(atylabmongo.MongoDatabaseStructMock)
				mongoDatabaseMock.On("Collection", model.RoomCollectionName).Return(mongoCollectionMock)
				mongoConnectorStruct := &atylabmongo.MongoConnector{
					Db: mongoDatabaseMock,
				}
				mongoConnectionStructMock := setupInitMock(tt.initErr, mongoConnectorStruct)
				mongoUseCase := usecase.NewMongoUseCaseStruct(mongoConnectionStructMock, usecase.NewMongo())

				roomSvc := NewRoomSvcStruct(mongoUseCase)
				resultRooms, err := roomSvc.ListRooms(atylabmongo.NewMongoCtxSvc())
				if tt.returnErr {
					if err == nil {
						t.Errorf("ListRooms() expected error, got nil")
					}
					return
				}

				if err != nil {
					t.Errorf("ListRooms() error = %v", err)
					return
				}
				if len(resultRooms) != 2 {
					t.Errorf("ListRooms() rooms = %v, want non-empty slice", resultRooms)
				}
			})
		}
	})
}
