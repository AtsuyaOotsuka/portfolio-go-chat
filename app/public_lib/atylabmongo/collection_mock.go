package atylabmongo

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCollectionStructMock struct {
	mock.Mock
}

func (m *MongoCollectionStructMock) InsertOne(ctx context.Context, document interface{}) (string, error) {
	args := m.Called(ctx, document)
	return args.String(0), args.Error(1)
}

func (m *MongoCollectionStructMock) Find(ctx context.Context, filter interface{}) (MongoCursorInterface, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(MongoCursorInterface), args.Error(1)
}

func (m *MongoCollectionStructMock) FindOne(ctx context.Context, filter interface{}, object interface{}) error {
	args := m.Called(ctx, filter, object)
	return args.Error(0)
}

func (m *MongoCollectionStructMock) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MongoCollectionStructMock) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MongoCollectionStructMock) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}
