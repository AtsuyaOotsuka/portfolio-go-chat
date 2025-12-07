package atylabmongo

import "github.com/stretchr/testify/mock"

type MongoDatabaseStructMock struct {
	mock.Mock
}

func (m *MongoDatabaseStructMock) Collection(name string) MongoCollectionInterface {
	args := m.Called(name)
	return args.Get(0).(MongoCollectionInterface)
}
