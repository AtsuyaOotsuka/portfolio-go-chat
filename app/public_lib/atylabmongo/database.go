package atylabmongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo Database
type MongoDatabaseInterface interface {
	Collection(name string) MongoCollectionInterface
}

type MongoDatabaseStruct struct {
	db *mongo.Database
}

func NewMongoDatabaseStruct(
	db *mongo.Database,
) *MongoDatabaseStruct {
	return &MongoDatabaseStruct{
		db: db,
	}
}

func (m *MongoDatabaseStruct) Collection(name string) MongoCollectionInterface {
	return NewMongoCollectionStruct(m.db.Collection(name))
}
