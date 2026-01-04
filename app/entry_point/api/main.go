package main

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/app"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func SetupMongo() (*usecase.Mongo, error) {
	mongo := usecase.NewMongo()
	mongoUseCase := usecase.NewMongoUseCaseStruct(
		atylabmongo.NewMongoConnectionStruct(),
		mongo,
	)
	_, err := mongoUseCase.MongoInit()
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return nil, err
	}
	return mongo, nil
}

func SetupRouter(
	mongo *usecase.Mongo,
) *app.App {
	e := echo.New()
	e.Validator = &usecase.CustomValidator{Validator: validator.New()}
	app := app.NewApp()
	app.Init(e, mongo)

	return app
}

func main() {
	mongo, err := SetupMongo()
	if err != nil {
		fmt.Println("Error setting up MongoDB:", err)
		return
	}
	app := SetupRouter(mongo)
	defer app.Shutdown()

	for _, route := range app.Echo.Routes() {
		fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}
	app.Echo.Start(":8080")
}
