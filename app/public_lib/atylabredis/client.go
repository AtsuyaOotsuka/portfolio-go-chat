package atylabredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClientInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type RedisClientStruct struct {
	rdb *redis.Client
}

func NewRedisClientStruct(rdb *redis.Client) *RedisClientStruct {
	return &RedisClientStruct{
		rdb: rdb,
	}
}

func (r *RedisClientStruct) Get(
	ctx context.Context,
	key string,
) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisClientStruct) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	err := r.rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
