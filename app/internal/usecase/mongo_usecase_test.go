package usecase

import (
	"fmt"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
)

func TestNewMongoUseCaseStruct(t *testing.T) {
	mockMongoConnectorPkg := &atylabmongo.MongoConnectionStruct{}
	mongo := NewMongo()

	mongoUseCase := NewMongoUseCaseStruct(
		mockMongoConnectorPkg,
		mongo,
	)

	if mongoUseCase.mongoConnectorPkg != mockMongoConnectorPkg {
		t.Errorf("Expected mongoConnectorPkg to be set correctly")
	}
}

var mongoSvcEnvs = funcs.Envs{
	"MONGO_HOST": "localhost",
	"MONGO_PORT": "27017",
	"MONGO_USER": "testuser",
	"MONGO_PASS": "testpass",
}

func TestMakeUri(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		uri, err := makeUri()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expectedUri := "mongodb://testuser:testpass@localhost:27017"
		if uri != expectedUri {
			t.Errorf("Expected URI %s, got %s", expectedUri, uri)
		}
	})
}

func TestMakeUriNotEnv(t *testing.T) {
	uri, err := makeUri()
	if err == nil {
		t.Fatalf("Expected error due to missing env vars, got nil")
	}
	if uri != "" {
		t.Errorf("Expected empty URI on error, got %s", uri)
	}
}

func TestMongoInit(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		mockMongoConnectorPkg := &atylabmongo.MongoConnectionStructMock{}
		mockMongoConnectorPkg.On("NewMongoConnect", "chatapp", "mongodb://testuser:testpass@localhost:27017").Return(&atylabmongo.MongoConnector{}, nil)

		mongoUseCase := NewMongoUseCaseStruct(
			mockMongoConnectorPkg,
			NewMongo(),
		)

		mongo, err := mongoUseCase.MongoInit()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if mongo.MongoConnector == nil {
			t.Errorf("Expected MongoConnector to be initialized")
		}
	})
}

func TestMongoInitAlreadyConnected(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		mockMongoConnectorPkg := &atylabmongo.MongoConnectionStructMock{}

		mongoUseCase := NewMongoUseCaseStruct(
			mockMongoConnectorPkg,
			&Mongo{
				IsConnected: true,
			},
		)

		mongo, err := mongoUseCase.MongoInit()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if mongo.IsConnected != true {
			t.Errorf("Expected IsConnected to be true")
		}
	})
}

func TestMongoInitNewMongoConnectError(t *testing.T) {
	funcs.WithEnvMap(mongoSvcEnvs, t, func() {
		mockMongoConnectorPkg := &atylabmongo.MongoConnectionStructMock{}
		mockMongoConnectorPkg.On("NewMongoConnect", "chatapp", "mongodb://testuser:testpass@localhost:27017").Return(&atylabmongo.MongoConnector{}, fmt.Errorf("connection error"))

		mongoUseCase := NewMongoUseCaseStruct(
			mockMongoConnectorPkg,
			NewMongo(),
		)

		mongo, err := mongoUseCase.MongoInit()
		if err == nil {
			t.Fatalf("Expected error, got Mongo %v", mongo)
		}
		if mongo == nil {
			t.Errorf("Expected Mongo to be non-nil even on error")
		}
	})
}

func TestMongoInitNotEnv(t *testing.T) {
	mockMongoConnectorPkg := &atylabmongo.MongoConnectionStructMock{}
	mongoUseCase := NewMongoUseCaseStruct(
		mockMongoConnectorPkg,
		NewMongo(),
	)

	_, err := mongoUseCase.MongoInit()
	if err == nil {
		t.Fatalf("Expected error due to missing env vars, got nil")
	}
}
