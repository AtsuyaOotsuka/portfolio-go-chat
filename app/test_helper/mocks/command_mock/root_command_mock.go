package command_mock

import "github.com/stretchr/testify/mock"

type RootCommandMock struct {
	mock.Mock
}

func (m *RootCommandMock) Run() {
	m.Called()
}
