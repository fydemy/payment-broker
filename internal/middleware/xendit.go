package middleware

import (
	"os"
	"payment-broker/internal/service"

	"github.com/gofiber/fiber/v2"
)

func XenditMiddleware(tenantService service.TenantService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		APIKey := c.Get("X-Api-Key")

		if APIKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key is required",
			})
		}

		tenant, err := tenantService.CheckAPIKey(ctx, APIKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to validate API key",
			})
		}

		if tenant.AccountID == "" || tenant.ID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		c.Locals("X-Account-ID", tenant.AccountID)
		c.Locals("X-Tenant-ID", tenant.ID)
		return c.Next()
	}
}

func XenditWebhookMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		callbackToken := c.Get("x-callback-token")

		if callbackToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Callback token is required",
			})
		}

		if callbackToken != os.Getenv("XENDIT_CALLBACK_TOKEN") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid callback token",
			})
		}
		return c.Next()
	}
}
