package usecase

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabredis"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	r := NewRedis()
	if r == nil {
		t.Errorf("NewRedis() returned nil")
	}
	if r.IsConnected != false {
		t.Errorf("NewRedis() expected IsConnected to be false, got %v", r.IsConnected)
	}
}

func TestNewRedisUseCaseStruct(t *testing.T) {
	redisConnectorPkg := atylabredis.NewRedisConnectorStruct()
	redis := NewRedis()
	useCase := NewRedisUseCaseStruct(redisConnectorPkg, redis)

	if useCase == nil {
		t.Errorf("NewRedisUseCaseStruct() returned nil")
	}
	if useCase.redisConnectorPkg != redisConnectorPkg {
		t.Errorf("NewRedisUseCaseStruct() expected redisConnectorPkg to be set correctly")
	}
	if useCase.redis != redis {
		t.Errorf("NewRedisUseCaseStruct() expected redis to be set correctly")
	}
}

var redisSvcEnvs = funcs.Envs{
	"REDIS_ADDR": "localhost:6379",
	"REDIS_PASS": "testpass",
	"REDIS_DB":   "0",
}

func TestRedisInit(t *testing.T) {
	redisConnectorPkg := new(atylabredis.RedisConnectorStructMock)
	redisConnectorPkg.On("NewRedisConnect", "localhost:6379", "testpass", 0).Return(
		&atylabredis.RedisConnector{},
		nil,
	)
	redis := NewRedis()
	useCase := NewRedisUseCaseStruct(redisConnectorPkg, redis)

	funcs.WithEnvMap(redisSvcEnvs, t, func() {
		r, err := useCase.RedisInit()
		if err != nil {
			t.Errorf("RedisInit() returned an error: %v", err)
		}
		if r == nil {
			t.Errorf("RedisInit() returned nil")
		}
		if !r.IsConnected {
			t.Errorf("RedisInit() expected IsConnected to be true, got false")
		}
	})
}

func TestRedisInit_AlreadyConnected(t *testing.T) {
	redisConnectorPkg := new(atylabredis.RedisConnectorStructMock)
	redis := &Redis{
		IsConnected: true,
	}
	useCase := NewRedisUseCaseStruct(redisConnectorPkg, redis)

	r, err := useCase.RedisInit()
	if err != nil {
		t.Errorf("RedisInit() returned an error: %v", err)
	}
	if r == nil {
		t.Errorf("RedisInit() returned nil")
	}
	if !r.IsConnected {
		t.Errorf("RedisInit() expected IsConnected to be true, got false")
	}

	redisConnectorPkg.AssertNotCalled(t, "NewRedisConnect")
}

func TestRedisInit_ConnectionError(t *testing.T) {
	redisConnectorPkg := new(atylabredis.RedisConnectorStructMock)
	redisConnectorPkg.On("NewRedisConnect", "localhost:6379", "testpass", 0).Return(
		&atylabredis.RedisConnector{},
		assert.AnError,
	)
	redis := NewRedis()
	useCase := NewRedisUseCaseStruct(redisConnectorPkg, redis)

	funcs.WithEnvMap(redisSvcEnvs, t, func() {
		r, err := useCase.RedisInit()
		if err == nil {
			t.Errorf("RedisInit() expected to return an error, got nil")
		}
		if r != nil {
			t.Errorf("RedisInit() expected to return nil, got %v", r)
		}
	})
}

func TestRedisInit_InvalidDBEnv(t *testing.T) {
	redisConnectorPkg := new(atylabredis.RedisConnectorStructMock)
	redis := NewRedis()
	useCase := NewRedisUseCaseStruct(redisConnectorPkg, redis)

	invalidDBEnvs := funcs.Envs{
		"REDIS_ADDR": "localhost:6379",
		"REDIS_PASS": "testpass",
		"REDIS_DB":   "invalid_int",
	}

	funcs.WithEnvMap(invalidDBEnvs, t, func() {
		r, err := useCase.RedisInit()
		if err == nil {
			t.Errorf("RedisInit() expected to return an error, got nil")
		}
		if r != nil {
			t.Errorf("RedisInit() expected to return nil, got %v", r)
		}
	})
}
