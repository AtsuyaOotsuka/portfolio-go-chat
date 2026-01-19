package atylabredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConnector struct {
	Client RedisClientInterface
}

type RedisConnectorInterface interface {
	NewRedisConnect(
		addr string,
		password string,
		db int,
	) (*RedisConnector, error)
}

type RedisConnectorStruct struct{}

func NewRedisConnectorStruct() *RedisConnectorStruct {
	return &RedisConnectorStruct{}
}

func (r *RedisConnectorStruct) NewRedisConnect(
	addr string,
	password string,
	db int,
) (*RedisConnector, error) {
	database, err := r.connect(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err != nil {
		return nil, err
	}

	redisConnector := &RedisConnector{}
	redisClient := NewRedisClientStruct(database)
	redisConnector.Client = redisClient

	return redisConnector, nil
}

func (r *RedisConnectorStruct) connect(opts *redis.Options) (*redis.Client, error) {
	rdb := redis.NewClient(opts)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	_, err := rdb.Ping(ctx).Result()
	fmt.Println("Connected to Redis")
	return rdb, err
}
