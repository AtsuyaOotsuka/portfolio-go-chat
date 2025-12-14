package provider

import "github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"

type Provider struct {
	mongo *usecase.Mongo
}

func NewProvider(
	mongo *usecase.Mongo,
) *Provider {
	return &Provider{
		mongo: mongo,
	}
}
