package app

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/provider"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/labstack/echo/v4"
)

type App struct {
	Echo       *echo.Echo
	middleware *middleware.Middleware
	provider   *provider.Provider
	mongo      *usecase.Mongo
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(
	e *echo.Echo,
	mongo *usecase.Mongo,
) {
	a.Echo = e
	a.mongo = mongo
	a.initProviders(mongo)
	a.initMiddlewares()
	a.entryGlobalMiddleware()
	a.entryRoutes()
}

func (a *App) Shutdown() {
	fmt.Println("Shutting down the application...")
	// ここにシャットダウン処理を追加
	fmt.Println("Application shut down completed.")
}
