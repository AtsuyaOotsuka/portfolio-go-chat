package provider

import "github.com/AtsuyaOotsuka/portfolio-go-chat/internal/handler"

func (p *Provider) BindHealthCheckHandler() *handler.HealthCheckHandler {
	return handler.NewHealthCheckHandler()
}
