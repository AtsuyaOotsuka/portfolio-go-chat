package usecase

import (
	"fmt"
	"os"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
)

type MongoUseCaseInterface interface {
	MongoInit() (*Mongo, error)
}

type Mongo struct {
	MongoConnector *atylabmongo.MongoConnector
	IsConnected    bool
}

func NewMongo() *Mongo {
	return &Mongo{
		IsConnected: false,
	}
}

type MongoUseCaseStruct struct {
	mongoConnectorPkg atylabmongo.MongoConnectorInterface
	mongo             *Mongo
}

func NewMongoUseCaseStruct(
	mongoConnectorPkg atylabmongo.MongoConnectorInterface,
	mongo *Mongo,
) *MongoUseCaseStruct {
	return &MongoUseCaseStruct{
		mongoConnectorPkg: mongoConnectorPkg,
		mongo:             mongo,
	}
}

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

func (s *MongoUseCaseStruct) MongoInit() (*Mongo, error) {
	if s.mongo != nil && s.mongo.IsConnected {
		return s.mongo, nil
	}

	uri, err := makeUri()
	fmt.Println("MongoDB URI:", uri)
	if err != nil {
		return nil, err
	}

	mongoConnector, err := s.mongoConnectorPkg.NewMongoConnect("chatapp", uri)
	if err != nil {
		return &Mongo{}, err
	}

	s.mongo = &Mongo{
		MongoConnector: mongoConnector,
		IsConnected:    true,
	}
	return s.mongo, nil
}
