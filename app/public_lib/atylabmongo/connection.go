package atylabmongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnector struct {
	Db     MongoDatabaseInterface
	Ctx    context.Context
	Cancel context.CancelFunc
}

type MongoConnectorInterface interface {
	NewMongoConnect(database string, mongoUri string) (*MongoConnector, error)
}

func NewMongoConnectionStruct() *MongoConnectionStruct {
	return &MongoConnectionStruct{}
}

type MongoConnectionStruct struct{}

func (m *MongoConnectionStruct) NewMongoConnect(database string, mongoUri string) (*MongoConnector, error) {
	client, ctx, cancelFunc, err := m.connect(mongoUri)
	if err != nil {
		return &MongoConnector{}, err
	}

	mongoConnector := &MongoConnector{}
	mongoConnector.Ctx = ctx
	mongoConnector.Cancel = cancelFunc
	mongoClient := NewMongoClientStruct(client)
	mongoConnector.Db = mongoClient.Database(database)
	fmt.Println("Connected to MongoDB!")

	return mongoConnector, nil
}

func (m *MongoConnectionStruct) connect(mongoUri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	// タイムアウト付きのcontext
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)

	clientOptions := options.Client().ApplyURI(mongoUri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		defer cancelFunc()
		return nil, nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		defer cancelFunc()
		return nil, nil, nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client, ctx, cancelFunc, nil
}
