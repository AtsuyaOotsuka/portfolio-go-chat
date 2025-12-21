package command_mock

import "github.com/stretchr/testify/mock"

type VersionCommandMock struct {
	mock.Mock
}

func (m *VersionCommandMock) Run(args []string) {
	m.Called(args)
}
