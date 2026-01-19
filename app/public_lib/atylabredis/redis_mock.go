package atylabredis

import "github.com/stretchr/testify/mock"

type RedisConnectorStructMock struct {
	mock.Mock
}

func (r *RedisConnectorStructMock) NewRedisConnect(addr string, password string, db int) (*RedisConnector, error) {
	args := r.Called(addr, password, db)
	return args.Get(0).(*RedisConnector), args.Error(1)
}
