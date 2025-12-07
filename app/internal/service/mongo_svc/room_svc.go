package mongo_svc

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"go.mongodb.org/mongo-driver/bson"
)

type RoomSvcInterface interface {
	GetRoomList(uuid string, target string) ([]model.Room, error)
	CreateRoom(room model.Room) (string, error)
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

func (s *RoomSvcStruct) GetRoomList(uuid string, target string) ([]model.Room, error) {
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

	defer mongo.MongoConnector.Cancel()

	collection := mongo.MongoConnector.Db.Collection("rooms")

	cursor, err := collection.Find(mongo.MongoConnector.Ctx, filter)
	if err != nil {
		fmt.Println("Failed to find rooms:", err)
		return []model.Room{}, err
	}
	defer cursor.Close(mongo.MongoConnector.Ctx)

	var rooms []model.Room
	for cursor.Next(mongo.MongoConnector.Ctx) {
		var room model.Room
		if err := cursor.Decode(&room); err != nil {
			fmt.Println("Failed to decode room:", err)
			return []model.Room{}, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (s *RoomSvcStruct) CreateRoom(room model.Room) (string, error) {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return "", err
	}

	defer mongo.MongoConnector.Cancel()

	collection := mongo.MongoConnector.Db.Collection("rooms")

	InsertedID, err := collection.InsertOne(mongo.MongoConnector.Ctx, room)
	if err != nil {
		return "", err
	}

	return InsertedID, nil
}
