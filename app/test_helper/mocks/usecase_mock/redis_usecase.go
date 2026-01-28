package usecase_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/stretchr/testify/mock"
)

type RedisUseCaseMock struct {
	mock.Mock
}

func (m *RedisUseCaseMock) RedisInit() (*usecase.Redis, error) {
	args := m.Called()
	return args.Get(0).(*usecase.Redis), args.Error(1)
}
