package mongo_svc

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomSvcInterface interface {
	GetRoomList(uuid string, target string, ctx *atylabmongo.MongoCtxSvc) ([]model.Room, error)
	CreateRoom(room model.Room, ctx *atylabmongo.MongoCtxSvc) (string, error)
	GetRoomByID(roomID string, ctx *atylabmongo.MongoCtxSvc) (model.Room, error)
	JoinRoom(roomID string, uuid string, ctx *atylabmongo.MongoCtxSvc) error
	LeaveRoom(roomID string, uuid string, ctx *atylabmongo.MongoCtxSvc) error
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

func (s *RoomSvcStruct) GetRoomList(uuid string, target string, ctx *atylabmongo.MongoCtxSvc) ([]model.Room, error) {
	var filter bson.M
	switch target {
	case "all":
		filter = bson.M{
			"$or": []bson.M{
				{"isprivate": false},
				{"members": uuid}, // 参加済みの場合はプライベートでも表示
			},
		}
	case "joined":
		filter = bson.M{"members": uuid} // 参加済みのものだけ
	default:
		return nil, fmt.Errorf("invalid target: %s", target)
	}

	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return []model.Room{}, err
	}

	collection := mongo.MongoConnector.Db.Collection("rooms")

	cursor, err := collection.Find(ctx.Ctx, filter)
	if err != nil {
		fmt.Println("Failed to find rooms:", err)
		return []model.Room{}, err
	}
	defer cursor.Close(ctx.Ctx)

	var rooms []model.Room
	for cursor.Next(ctx.Ctx) {
		var room model.Room
		if err := cursor.Decode(&room); err != nil {
			fmt.Println("Failed to decode room:", err)
			return []model.Room{}, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (s *RoomSvcStruct) CreateRoom(room model.Room, ctx *atylabmongo.MongoCtxSvc) (string, error) {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return "", err
	}

	collection := mongo.MongoConnector.Db.Collection("rooms")

	InsertedID, err := collection.InsertOne(ctx.Ctx, room)
	if err != nil {
		return "", err
	}

	return InsertedID, nil
}

func (s *RoomSvcStruct) GetRoomByID(roomID string, ctx *atylabmongo.MongoCtxSvc) (model.Room, error) {
	var err error
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return model.Room{}, err
	}

	collection := mongo.MongoConnector.Db.Collection("rooms")

	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return model.Room{}, err
	}

	var room model.Room
	err = collection.FindOne(ctx.Ctx, bson.M{"_id": id}, &room)
	if err != nil {
		return model.Room{}, err
	}

	return room, nil

}

func (s *RoomSvcStruct) JoinRoom(roomID string, uuid string, ctx *atylabmongo.MongoCtxSvc) error {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return err
	}

	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	collection := mongo.MongoConnector.Db.Collection("rooms")

	_, err = collection.UpdateOne(
		ctx.Ctx,
		bson.M{"_id": id},
		bson.M{"$addToSet": bson.M{"members": uuid}},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomSvcStruct) LeaveRoom(roomID string, uuid string, ctx *atylabmongo.MongoCtxSvc) error {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return err
	}

	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	collection := mongo.MongoConnector.Db.Collection("rooms")

	_, err = collection.UpdateOne(
		ctx.Ctx,
		bson.M{"_id": id},
		bson.M{"$pull": bson.M{"members": uuid}},
	)
	if err != nil {
		return err
	}

	return nil
}
