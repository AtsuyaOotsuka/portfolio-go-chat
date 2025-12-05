package routing

import "github.com/AtsuyaOotsuka/portfolio-go-chat/internal/handler"

func (r *Routing) HealthCheckRoute(
	handler handler.HealthCheckHandlerInterface,
) {
	r.echo.GET("/healthcheck", handler.Check)
	r.echo.POST("/healthcheck_post", handler.Check)
}
