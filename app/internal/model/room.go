package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string
	OwnerID   string
	CreatedAt time.Time
	Members   []string
	IsPrivate bool
}
