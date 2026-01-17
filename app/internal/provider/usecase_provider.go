package provider

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabredis"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
)

func (p *Provider) bindMongoSvc() *usecase.MongoUseCaseStruct {
	return usecase.NewMongoUseCaseStruct(
		atylabmongo.NewMongoConnectionStruct(),
		p.mongo,
	)
}

func (p *Provider) bindRedisSvc() *usecase.RedisUseCaseStruct {
	return usecase.NewRedisUseCaseStruct(
		atylabredis.NewRedisConnectorStruct(),
		p.redis,
	)
}
