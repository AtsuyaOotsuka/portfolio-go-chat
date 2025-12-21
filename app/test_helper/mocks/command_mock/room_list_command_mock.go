package command_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/stretchr/testify/mock"
)

type RoomListCommandMock struct {
	mock.Mock
}

func (m *RoomListCommandMock) Run(args []string) {
	m.Called(args)
}

func (m *RoomListCommandMock) SetUp(mongo usecase.MongoUseCaseInterface) {
	m.Called(mongo)
}
