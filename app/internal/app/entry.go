package app

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/routing"
)

func (a *App) entryGlobalMiddleware() {
	// 前処理系ミドルウェアをここに追加
	a.Echo.Use(a.middleware.Jwt)
	a.Echo.Use(a.middleware.Csrf)
	// 後処理系ミドルウェアをここに追加
	// 例: a.Echo.Use(a.middleware.Logging)
}

func (a *App) entryRoutes() {
	routing := routing.NewRouting(a.Echo, a.middleware)

	routing.HealthCheckRoute(
		a.provider.BindHealthCheckHandler(),
	)

	routing.RoomRoute(
		a.provider.BindRoomHandler(),
	)
}
