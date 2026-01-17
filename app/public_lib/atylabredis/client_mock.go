package atylabredis

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type RedisClientStructMock struct {
	mock.Mock
}

func (r *RedisClientStructMock) Get(
	ctx context.Context,
	key string,
) (string, error) {
	args := r.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (r *RedisClientStructMock) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	args := r.Called(ctx, key, value, expiration)
	return args.Error(0)
}
