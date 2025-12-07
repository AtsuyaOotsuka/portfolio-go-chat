package atylabmongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}) (string, error)
	Find(ctx context.Context, filter interface{}) (cursor MongoCursorInterface, err error)
	FindOne(ctx context.Context, filter interface{}, object interface{}) error
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
}

type MongoCollectionStruct struct {
	coll *mongo.Collection
}

func NewMongoCollectionStruct(coll *mongo.Collection) *MongoCollectionStruct {
	return &MongoCollectionStruct{
		coll: coll,
	}
}

func (r *MongoCollectionStruct) InsertOne(ctx context.Context, document interface{}) (string, error) {
	result, err := r.coll.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("InsertedID is not an ObjectID: %#v", result.InsertedID)
	}

	return id.Hex(), nil
}

func (r *MongoCollectionStruct) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return r.coll.UpdateMany(ctx, filter, update)
}

func (r *MongoCollectionStruct) Find(ctx context.Context, filter interface{}) (cursor MongoCursorInterface, err error) {
	cursor, err = r.coll.Find(ctx, filter)
	return
}

func (r *MongoCollectionStruct) FindOne(ctx context.Context, filter interface{}, object interface{}) error {
	err := r.coll.FindOne(ctx, filter).Decode(object)
	return err
}

func (r *MongoCollectionStruct) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return r.coll.UpdateOne(ctx, filter, update)
}

func (r *MongoCollectionStruct) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return r.coll.DeleteOne(ctx, filter)
}
