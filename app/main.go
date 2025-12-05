package main

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/app"
	"github.com/labstack/echo/v4"
)

func SetupRouter() *app.App {
	e := echo.New()
	app := app.NewApp()
	app.Init(e)

	return app
}

func main() {
	app := SetupRouter()
	app.Echo.Start(":8080")
}
