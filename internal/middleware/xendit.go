package middleware

import (
	"payment-broker/internal/service"

	"github.com/gofiber/fiber/v2"
)

func XenditMiddleware(tenantService service.TenantService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		APIKey := c.Get("X-API-Key")

		if APIKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key is required",
			})
		}

		accountID, err := tenantService.CheckAPIKey(ctx, APIKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to validate API key",
			})
		}

		if accountID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		c.Set("X-Account-ID", *accountID)
		return c.Next()
	}
}
