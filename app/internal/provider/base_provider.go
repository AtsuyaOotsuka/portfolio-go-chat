package provider

import "github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"

type Provider struct {
	mongo *usecase.Mongo
	redis *usecase.Redis
}

func NewProvider(
	mongo *usecase.Mongo,
	redis *usecase.Redis,
) *Provider {
	return &Provider{
		mongo: mongo,
		redis: redis,
	}
}
