package atylabmongo

import "github.com/stretchr/testify/mock"

type MongoClientStructMock struct {
	mock.Mock
}

func (m *MongoClientStructMock) Database(name string) MongoDatabaseInterface {
	args := m.Called(name)
	return args.Get(0).(MongoDatabaseInterface)
}
