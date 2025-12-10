package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	RoomID        string             `bson:"roomid"`
	Sender        string             `bson:"sender"`
	Message       string             `bson:"message"`
	CreatedAt     time.Time          `bson:"createdAt"`
	IsReadUserIds []string           `bson:"isReadUserIds"`
}
