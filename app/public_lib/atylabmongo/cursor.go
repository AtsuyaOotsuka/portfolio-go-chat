package atylabmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCursorInterface interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
}

type MongoCursorStruct struct {
	cursor *mongo.Cursor
}

func (r *MongoCursorStruct) Next(ctx context.Context) bool {
	return r.cursor.Next(ctx)
}

func (r *MongoCursorStruct) Decode(val interface{}) error {
	return r.cursor.Decode(val)
}

func (r *MongoCursorStruct) Close(ctx context.Context) error {
	return r.cursor.Close(ctx)
}
