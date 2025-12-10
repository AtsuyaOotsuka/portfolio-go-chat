package mongo_svc

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageSvcInterface interface {
	SendMessage(message model.Message) (string, error)
	GetMessageList(roomID string) ([]model.Message, error)
	ReadMessages(messageIds []string, roomId string, userId string) error
}

type MessageSvcStruct struct {
	mongo usecase.MongoUseCaseInterface
}

func NewMessageSvcStruct(
	mongo usecase.MongoUseCaseInterface,
) *MessageSvcStruct {
	return &MessageSvcStruct{
		mongo: mongo,
	}
}

func (s *MessageSvcStruct) SendMessage(message model.Message) (string, error) {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return "", err
	}
	defer mongo.MongoConnector.Cancel()

	collection := mongo.MongoConnector.Db.Collection("messages")
	InsertedID, err := collection.InsertOne(mongo.MongoConnector.Ctx, message)
	if err != nil {
		return "", err
	}

	return InsertedID, nil
}

func (s *MessageSvcStruct) GetMessageList(roomID string) ([]model.Message, error) {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return []model.Message{}, err
	}
	defer mongo.MongoConnector.Cancel()

	collection := mongo.MongoConnector.Db.Collection("messages")
	filter := bson.M{"roomid": roomID}

	cursor, err := collection.Find(mongo.MongoConnector.Ctx, filter)
	if err != nil {
		fmt.Println("Failed to find messages:", err)
		return []model.Message{}, err
	}

	defer cursor.Close(mongo.MongoConnector.Ctx)

	var messages []model.Message
	for cursor.Next(mongo.MongoConnector.Ctx) {
		var message model.Message
		if err := cursor.Decode(&message); err != nil {
			fmt.Println("Failed to decode message:", err)
			return []model.Message{}, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (s *MessageSvcStruct) ReadMessages(messageIds []string, roomId string, userId string) error {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return err
	}
	defer mongo.MongoConnector.Cancel()

	collection := mongo.MongoConnector.Db.Collection("messages")
	var chatObjectIDs []primitive.ObjectID
	for _, id := range messageIds {
		fmt.Println("Processing message ID:", id)
		chatObjectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		chatObjectIDs = append(chatObjectIDs, chatObjectID)
	}

	filter := bson.M{
		"_id":    bson.M{"$in": chatObjectIDs},
		"roomid": roomId,
	}

	update := bson.M{
		"$addToSet": bson.M{"isReadUserIds": userId},
	}

	_, err = collection.UpdateMany(mongo.MongoConnector.Ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
