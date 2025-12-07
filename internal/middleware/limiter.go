package middleware

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
)

func LimiterMiddleware(limiter *redis_rate.Limiter, limit redis_rate.Limit) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		res, err := limiter.Allow(ctx, c.IP(), limit)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error handling requests",
			})
		}

		if res.Allowed == 0 {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "too many requests",
			})
		}

		return c.Next()
	}
}
