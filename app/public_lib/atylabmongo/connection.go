package atylabmongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnector struct {
	Db MongoDatabaseInterface
}

type MongoConnectorInterface interface {
	NewMongoConnect(database string, mongoUri string) (*MongoConnector, error)
}

func NewMongoConnectionStruct() *MongoConnectionStruct {
	return &MongoConnectionStruct{}
}

type MongoConnectionStruct struct{}

func (m *MongoConnectionStruct) NewMongoConnect(database string, mongoUri string) (*MongoConnector, error) {
	client, err := m.connect(mongoUri)
	if err != nil {
		return nil, err
	}

	mongoConnector := &MongoConnector{}
	mongoClient := NewMongoClientStruct(client)
	mongoConnector.Db = mongoClient.Database(database)
	fmt.Println("Connected to MongoDB!")

	return mongoConnector, nil
}

func (m *MongoConnectionStruct) connect(mongoUri string) (*mongo.Client, error) {
	// タイムアウト付きのcontext
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	defer fmt.Println("Created new MongoDB context")

	clientOptions := options.Client().ApplyURI(mongoUri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client, nil
}
