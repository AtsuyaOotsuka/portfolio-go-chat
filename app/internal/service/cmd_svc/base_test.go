package cmd_svc

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mongoSvcEnvs = funcs.Envs{
	"MONGO_HOST": "localhost",
	"MONGO_PORT": "27017",
	"MONGO_USER": "testuser",
	"MONGO_PASS": "testpass",
}

func setupInitMock(wantErr bool, mongoConnectorStruct *atylabmongo.MongoConnector) atylabmongo.MongoConnectorInterface {
	var returnErr error
	if wantErr {
		returnErr = assert.AnError
	}
	mongoConnectionStructMock := new(atylabmongo.MongoConnectionStructMock)
	mongoConnectionStructMock.On("NewMongoConnect", "chatapp", mock.Anything).Return(mongoConnectorStruct, returnErr)
	return mongoConnectionStructMock
}
