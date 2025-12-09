package main

import (
	"os"
	"payment-broker/internal/app"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// @title		Payment Broker
// @version	1.0
// @host		localhost:3000
// @BasePath	/v1
func main() {
	godotenv.Load(".env")

	logger := app.InitLogger()
	defer logger.Sync()

	cache := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PWD"),
		DB:       0,
	})

	db := app.InitDB(logger)
	router := app.InitApp(db, logger, cache)

	fapp := fiber.New()
	fapp.Use(recover.New())
	fapp.Use(swagger.New(swagger.Config{
		FilePath: "./docs/swagger.json",
	}))

	fapp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Content-Type, X-API-Key",
		AllowMethods: "POST",
	}))

	app.InitRouter(fapp, router, cache)

	if err := fapp.Listen(":" + os.Getenv("APP_PORT")); err != nil {
		logger.Error("failed to start server: ", zap.Error(err))
	}
}
