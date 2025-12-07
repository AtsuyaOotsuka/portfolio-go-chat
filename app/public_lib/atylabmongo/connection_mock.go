package atylabmongo

import (
	"github.com/stretchr/testify/mock"
)

type MongoConnectionStructMock struct {
	mock.Mock
}

func (m *MongoConnectionStructMock) NewMongoConnect(database string, mongoUri string) (*MongoConnector, error) {
	args := m.Called(database, mongoUri)
	return args.Get(0).(*MongoConnector), args.Error(1)
}

func (m *MongoConnectionStructMock) Cancel() {
	m.Called()
}
