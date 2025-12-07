package router

import (
	"payment-broker/internal/controller"
	"payment-broker/internal/middleware"
	"payment-broker/internal/service"

	"github.com/gofiber/fiber/v2"
)

func NewXenditRouter(app fiber.Router, tenantService service.TenantService, xenditController controller.XenditController) {
	xenditAPI := app.Group("/xendit")

	xenditAPI.Use(middleware.XenditMiddleware(tenantService))
	xenditAPI.Post("/invoices", xenditController.CreatePayment)
}
