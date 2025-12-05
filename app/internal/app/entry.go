package app

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/routing"
)

func (a *App) entryBeforeGlobalMiddleware() {
	// 前処理系ミドルウェアをここに追加
	a.Echo.Use(a.middleware.Jwt)
	a.Echo.Use(a.middleware.Csrf)
}

func (a *App) entryAfterGlobalMiddleware() {
	// 後処理系ミドルウェアをここに追加
}

func (a *App) entryRoutes() {
	routing := routing.NewRouting(a.Echo, a.middleware)

	routing.HealthCheckRoute(
		a.provider.BindHealthCheckHandler(),
	)
}
