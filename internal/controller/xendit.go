package controller

import (
	"payment-broker/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type XenditController interface {
	CreatePayment(c *fiber.Ctx) error
}

type xenditController struct {
	logger        *zap.Logger
	xenditService service.XenditService
}

func NewXenditController(logger *zap.Logger, xenditService service.XenditService) XenditController {
	return &xenditController{
		logger:        logger,
		xenditService: xenditService,
	}
}

func (t *xenditController) CreatePayment(c *fiber.Ctx) error {
	ctx := c.Context()
	body := c.Locals("body")
	accountID := c.Get("X-Account-ID")

	resp, err := t.xenditService.CreatePayment(ctx, body, accountID)

	if err != nil {
		t.logger.Error("Xendit API error", zap.Error(err))
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Failed to process payment",
		})
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	return c.Status(resp.StatusCode).Send(resp.Body)
}
