package app

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/provider"
	"github.com/labstack/echo/v4"
)

type App struct {
	Echo       *echo.Echo
	middleware *middleware.Middleware
	provider   *provider.Provider
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(e *echo.Echo) {
	a.Echo = e
	a.initProviders()
	a.initMiddlewares()
	a.entryBeforeGlobalMiddleware()
	a.entryRoutes()
	a.entryAfterGlobalMiddleware()
}
