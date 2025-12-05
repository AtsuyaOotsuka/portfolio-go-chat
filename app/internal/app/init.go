package app

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/provider"
)

func (a *App) initProviders() {
	a.provider = provider.NewProvider()

}

func (a *App) initMiddlewares() {
	// ミドルウェアの初期化
	a.middleware = middleware.NewMiddleware(a.Echo)

}
