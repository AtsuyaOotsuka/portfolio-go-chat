package cmd_svc

import (
	"fmt"
	"strings"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/consts"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"go.mongodb.org/mongo-driver/bson"
)

type MessageSvcInterface interface {
	GetMessageList(roomID string, ctx *atylabmongo.MongoCtxSvc) ([]model.Message, error)
	ContainsForbiddenWords(message string) bool
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

func (s *MessageSvcStruct) GetMessageList(roomID string, ctx *atylabmongo.MongoCtxSvc) ([]model.Message, error) {
	mongo, err := s.mongo.MongoInit()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return []model.Message{}, err
	}

	collection := mongo.MongoConnector.Db.Collection(model.MessageCollectionName)
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

func (s *MessageSvcStruct) ContainsForbiddenWords(message string) bool {
	for _, word := range consts.ForbiddenWords {
		if strings.Contains(message, word) {
			return true
		}
	}
	return false
}
