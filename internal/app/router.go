package app

import (
	"payment-broker/internal/middleware"
	"payment-broker/internal/router"

	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func InitRouter(f *fiber.App, app *App, redis *redis.Client) {
	limiter := redis_rate.NewLimiter(redis)

	api := f.Group("/v1")
	api.Use(middleware.LimiterMiddleware(limiter, redis_rate.PerSecond(5)))

	router.NewXenditRouter(api, app.Service.Tenant, app.Controller.Xendit)
}
