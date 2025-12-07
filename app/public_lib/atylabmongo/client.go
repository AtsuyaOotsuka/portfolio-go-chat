package atylabmongo

import "go.mongodb.org/mongo-driver/mongo"

type MongoClientInterface interface {
	Database(name string) MongoDatabaseInterface
}

type MongoClientStruct struct {
	client *mongo.Client
}

func NewMongoClientStruct(
	client *mongo.Client,
) *MongoClientStruct {
	return &MongoClientStruct{
		client: client,
	}
}

func (r *MongoClientStruct) Database(name string) MongoDatabaseInterface {
	return NewMongoDatabaseStruct(r.client.Database(name))
}
