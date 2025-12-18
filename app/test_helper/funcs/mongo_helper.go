package funcs

import (
	"context"
	"fmt"
	"os"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func makeUri() (string, error) {
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	user := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASS")

	if user != "" && pass != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port), nil
	}
	return "", fmt.Errorf("incomplete MongoDB connection information")
}

func SetUpMongoTestDatabase() *TestMongoStruct {
	ctx := context.Background()
	uri, err := makeUri()
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db := client.Database("chatapp")
	return &TestMongoStruct{
		DB:     db,
		Ctx:    ctx,
		Client: client,
	}
}

type TestMongoStruct struct {
	DB     *mongo.Database
	Ctx    context.Context
	Client *mongo.Client
}

func (m *TestMongoStruct) Disconnect() error {
	fmt.Println("Disconnecting MongoDB client...")
	return m.Client.Disconnect(m.Ctx)
}

func (m *TestMongoStruct) MongoCleanUp() error {

	var err error

	err = m.DB.Collection(model.RoomCollectionName).Drop(m.Ctx)
	if err != nil {
		return err
	}
	err = m.DB.Collection(model.MessageCollectionName).Drop(m.Ctx)
	if err != nil {
		return err
	}

	fmt.Println("MongoDB cleaned up for tests.")
	return nil
}

func (m *TestMongoStruct) Insert(collectionName string, doc interface{}) (string, error) {
	result, err := m.DB.Collection(collectionName).InsertOne(m.Ctx, doc)
	if err != nil {
		return "", err
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("InsertedID is not an ObjectID: %#v", result.InsertedID)
	}
	InsertId := id.Hex()
	fmt.Printf("Inserted Room with ID: %s\n", InsertId)
	return InsertId, nil
}

func (m *TestMongoStruct) FindOneContents(collectionName string, id string) (*mongo.SingleResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := primitive.M{"_id": objID}
	result := m.DB.Collection(collectionName).FindOne(m.Ctx, filter)
	return result, nil
}

func (m *TestMongoStruct) ExistContents(collectionName string, filter interface{}) (bool, error) {
	count, err := m.DB.Collection(collectionName).CountDocuments(m.Ctx, filter)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (m *TestMongoStruct) CountContents(collectionName string, filter interface{}) (int64, error) {
	count, err := m.DB.Collection(collectionName).CountDocuments(m.Ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
