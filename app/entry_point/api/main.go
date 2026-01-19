package main

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/app"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabredis"
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

func SetupRedis() (*usecase.Redis, error) {
	redis := usecase.NewRedis()
	redisUseCase := usecase.NewRedisUseCaseStruct(
		atylabredis.NewRedisConnectorStruct(),
		redis,
	)
	_, err := redisUseCase.RedisInit()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return nil, err
	}
	return redis, nil
}

func SetupRouter(
	mongo *usecase.Mongo,
	redis *usecase.Redis,
) *app.App {
	e := echo.New()
	e.Validator = &usecase.CustomValidator{Validator: validator.New()}
	app := app.NewApp()
	app.Init(e, mongo, redis)

	return app
}

func main() {
	mongo, err := SetupMongo()
	if err != nil {
		fmt.Println("Error setting up MongoDB:", err)
		return
	}
	redis, err := SetupRedis()
	if err != nil {
		fmt.Println("Error setting up Redis:", err)
		return
	}
	app := SetupRouter(mongo, redis)
	defer app.Shutdown()

	for _, route := range app.Echo.Routes() {
		fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}
	app.Echo.Start(":8080")
}
