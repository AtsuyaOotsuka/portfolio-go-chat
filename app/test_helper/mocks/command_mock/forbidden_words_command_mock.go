package command_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/stretchr/testify/mock"
)

type ForbiddenWordsCommandMock struct {
	mock.Mock
}

func (m *ForbiddenWordsCommandMock) Run(args []string) {
	m.Called(args)
}

func (m *ForbiddenWordsCommandMock) SetUp(mongo usecase.MongoUseCaseInterface, timeOut int) {
	m.Called(mongo, timeOut)
}
