package main

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/app"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func SetupRouter() *app.App {
	e := echo.New()
	e.Validator = &usecase.CustomValidator{Validator: validator.New()}
	app := app.NewApp()
	app.Init(e)

	return app
}

func main() {
	app := SetupRouter()
	for _, route := range app.Echo.Routes() {
		fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}
	app.Echo.Start(":8080")
}
