package usecase

import (
	"fmt"
	"os"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
)

type MongoUseCaseInterface interface {
	MongoInit() (*Mongo, error)
}

type MongoUseCaseStruct struct {
	mongoConnectorPkg atylabmongo.MongoConnectorInterface
}

func NewMongoUseCaseStruct(
	mongoConnectorPkg atylabmongo.MongoConnectorInterface,
) *MongoUseCaseStruct {
	return &MongoUseCaseStruct{
		mongoConnectorPkg: mongoConnectorPkg,
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

type Mongo struct {
	MongoConnector *atylabmongo.MongoConnector
}

func (s *MongoUseCaseStruct) MongoInit() (*Mongo, error) {
	uri, err := makeUri()
	fmt.Println("MongoDB URI:", uri)
	if err != nil {
		return nil, err
	}

	mongoConnector, err := s.mongoConnectorPkg.NewMongoConnect("chatapp", uri)
	if err != nil {
		return &Mongo{}, err
	}

	mongo := &Mongo{
		MongoConnector: mongoConnector,
	}
	return mongo, nil
}
