package usecase

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabredis"
)

type RedisUseCaseInterface interface {
	RedisInit() (*Redis, error)
}

type Redis struct {
	RedisConnector *atylabredis.RedisConnector
	IsConnected    bool
}

func NewRedis() *Redis {
	return &Redis{
		IsConnected: false,
	}
}

type RedisUseCaseStruct struct {
	redisConnectorPkg atylabredis.RedisConnectorInterface
	redis             *Redis
}

func NewRedisUseCaseStruct(
	redisConnectorPkg atylabredis.RedisConnectorInterface,
	redis *Redis,
) *RedisUseCaseStruct {
	return &RedisUseCaseStruct{
		redisConnectorPkg: redisConnectorPkg,
		redis:             redis,
	}
}

func (s *RedisUseCaseStruct) RedisInit() (*Redis, error) {
	if s.redis != nil && s.redis.IsConnected {
		return s.redis, nil
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASS")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, fmt.Errorf("failed to convert REDIS_DB to int: %w", err)
	}

	redisConnector, err := s.redisConnectorPkg.NewRedisConnect(
		redisAddr,
		redisPass,
		redisDB,
	)

	if err != nil {
		return nil, err
	}

	s.redis = &Redis{
		RedisConnector: redisConnector,
		IsConnected:    true,
	}

	return s.redis, nil
}
