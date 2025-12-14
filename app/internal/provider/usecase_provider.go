package provider

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
)

func (p *Provider) bindMongoSvc() *usecase.MongoUseCaseStruct {
	return usecase.NewMongoUseCaseStruct(
		atylabmongo.NewMongoConnectionStruct(),
		p.mongo,
	)
}
