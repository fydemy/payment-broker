package router

import (
	"payment-broker/internal/controller"
	"payment-broker/internal/middleware"
	"payment-broker/internal/service"

	"github.com/gofiber/fiber/v2"
)

func NewXenditRouter(app fiber.Router, tenantService service.TenantService, xenditController controller.XenditController, webhookController controller.WebhookController) {
	xenditAPI := app.Group("/xendit")

	xenditAPIAction := xenditAPI.Group("/action")

	xenditAPIAction.Use(middleware.XenditMiddleware(tenantService))
	xenditAPIAction.Post("/invoices", xenditController.CreatePayment)
	xenditAPIAction.Post("/recurring/plans", xenditController.CreateSubscription)
	xenditAPIAction.Post("/payouts", xenditController.CreatePayout)
	xenditAPIAction.Post("/customers", xenditController.CreateCustomer)

	xenditAPIWebhook := xenditAPI.Group("/webhook")
	xenditAPIWebhook.Post("/", middleware.XenditWebhookMiddleware(), webhookController.WebhookHandler)
}
