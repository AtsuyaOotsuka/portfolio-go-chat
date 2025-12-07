package atylabmongo

import "github.com/stretchr/testify/mock"

type MongoPkgStructMock struct {
	mock.Mock
}

func (m *MongoPkgStructMock) MakeConnector(
	dbName string,
	uri string,
) (MongoConnectorInterface, error) {
	args := m.Called(dbName, uri)
	return args.Get(0).(MongoConnectorInterface), args.Error(1)
}
