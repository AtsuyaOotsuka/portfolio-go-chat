package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RoomCollectionName = "rooms"

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	OwnerID   string             `bson:"owner"`
	CreatedAt time.Time          `bson:"created_at"`
	Members   []string           `bson:"members"`
	IsPrivate bool               `bson:"is_private"`
}
