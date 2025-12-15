package app

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/provider"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
)

func (a *App) initProviders(
	mongo *usecase.Mongo,
) {
	a.provider = provider.NewProvider(
		mongo,
	)

}

func (a *App) initMiddlewares(
	mongo *usecase.Mongo,
) {
	// ミドルウェアの初期化
	a.middleware = &middleware.Middleware{
		Csrf: a.provider.BindCsrfMiddleware().Handler(),
		Jwt:  a.provider.BindJwtMiddleware().Handler(),
		Room: a.provider.BindRoomMiddleware().Handler(),
	}
}
