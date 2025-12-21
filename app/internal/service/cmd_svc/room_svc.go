package cmd_svc

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"go.mongodb.org/mongo-driver/bson"
)

type RoomSvcInterface interface {
	ListRooms(ctx *atylabmongo.MongoCtxSvc) ([]model.Room, error)
}

type RoomSvcStruct struct {
	mongo usecase.MongoUseCaseInterface
}

func NewRoomSvcStruct(
	mongo usecase.MongoUseCaseInterface,
) *RoomSvcStruct {
	return &RoomSvcStruct{
		mongo: mongo,
	}
}

func (s *RoomSvcStruct) ListRooms(ctx *atylabmongo.MongoCtxSvc) ([]model.Room, error) {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return []model.Room{}, err
	}

	collection := mongo.MongoConnector.Db.Collection(model.RoomCollectionName)
	cursor, err := collection.Find(ctx.Ctx, bson.M{})
	if err != nil {
		fmt.Println("Failed to find rooms:", err)
		return []model.Room{}, err
	}
	defer cursor.Close(ctx.Ctx)

	var rooms []model.Room
	if err = cursor.All(ctx.Ctx, &rooms); err != nil {
		fmt.Println("Failed to decode rooms:", err)
		return []model.Room{}, err
	}
	return rooms, nil
}
