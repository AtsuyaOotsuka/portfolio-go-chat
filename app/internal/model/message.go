package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	RoomID        string
	Sender        string
	Message       string
	CreatedAt     time.Time
	IsReadUserIds []string
}
