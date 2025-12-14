package mongo_svc

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageSvcInterface interface {
	SendMessage(message model.Message, ctx *atylabmongo.MongoCtxSvc) (string, error)
	GetMessageList(roomID string, ctx *atylabmongo.MongoCtxSvc) ([]model.Message, error)
	ReadMessages(messageIds []string, roomId string, userId string, ctx *atylabmongo.MongoCtxSvc) error
	IsSender(messageID string, roomID string, userID string, ctx *atylabmongo.MongoCtxSvc) error
	DeleteMessage(messageID string, roomID string, ctx *atylabmongo.MongoCtxSvc) error
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

func (s *MessageSvcStruct) SendMessage(message model.Message, ctx *atylabmongo.MongoCtxSvc) (string, error) {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return "", err
	}

	collection := mongo.MongoConnector.Db.Collection("messages")
	InsertedID, err := collection.InsertOne(ctx.Ctx, message)
	if err != nil {
		return "", err
	}

	return InsertedID, nil
}

func (s *MessageSvcStruct) GetMessageList(roomID string, ctx *atylabmongo.MongoCtxSvc) ([]model.Message, error) {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return []model.Message{}, err
	}

	collection := mongo.MongoConnector.Db.Collection("messages")
	filter := bson.M{"roomid": roomID}

	cursor, err := collection.Find(ctx.Ctx, filter)
	if err != nil {
		fmt.Println("Failed to find messages:", err)
		return []model.Message{}, err
	}

	defer cursor.Close(ctx.Ctx)

	var messages []model.Message
	for cursor.Next(ctx.Ctx) {
		var message model.Message
		if err := cursor.Decode(&message); err != nil {
			fmt.Println("Failed to decode message:", err)
			return []model.Message{}, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (s *MessageSvcStruct) ReadMessages(messageIds []string, roomId string, userId string, ctx *atylabmongo.MongoCtxSvc) error {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return err
	}

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

	_, err = collection.UpdateMany(ctx.Ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *MessageSvcStruct) IsSender(messageID string, roomID string, userID string, ctx *atylabmongo.MongoCtxSvc) error {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return err
	}

	collection := mongo.MongoConnector.Db.Collection("messages")
	messageObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":    messageObjectID,
		"roomid": roomID,
		"sender": userID,
	}

	var result model.Message
	err = collection.FindOne(ctx.Ctx, filter, &result)
	if err != nil {
		fmt.Println("User is not the sender of the message:", err)
		return err
	}

	return nil
}

func (s *MessageSvcStruct) DeleteMessage(messageID string, roomID string, ctx *atylabmongo.MongoCtxSvc) error {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return err
	}

	collection := mongo.MongoConnector.Db.Collection("messages")
	messageObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx.Ctx, bson.M{"_id": messageObjectID, "roomid": roomID})
	if err != nil {
		return err
	}

	return nil
}
