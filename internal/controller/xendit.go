package controller

import (
	"encoding/json"
	"fmt"
	"payment-broker/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type XenditController interface {
	CreatePayment(c *fiber.Ctx) error
	CreateSubscription(c *fiber.Ctx) error
	CreatePayout(c *fiber.Ctx) error
	CreateCustomer(c *fiber.Ctx) error
}

type XenditLibController struct {
	logger        *zap.Logger
	xenditService service.XenditService
}

func NewXenditController(logger *zap.Logger, xenditService service.XenditService) XenditController {
	return &XenditLibController{
		logger:        logger,
		xenditService: xenditService,
	}
}

func (t *XenditLibController) handleXenditRequest(c *fiber.Ctx, endpoint string, errorMsg string) error {
	ctx := c.Context()
	rawBody := c.Body()

	accountID := c.Locals("X-Account-ID").(string)
	tenantID := c.Locals("X-Tenant-ID")

	var data map[string]interface{}
	if err := json.Unmarshal(rawBody, &data); err != nil {
		t.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	switch endpoint {
	case "/v2/invoices":
		if externalID, ok := data["external_id"].(string); ok && externalID != "" {
			data["external_id"] = fmt.Sprintf("%s:%s", tenantID, externalID)
		}
	default:
		if referenceID, ok := data["reference_id"].(string); ok && referenceID != "" {
			data["reference_id"] = fmt.Sprintf("%s:%s", tenantID, referenceID)
		}
	}

	resp, err := t.xenditService.CreateRequest(ctx, data, endpoint, accountID)

	if err != nil {
		t.logger.Error("Xendit API error", zap.Error(err))
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": errorMsg,
		})
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	return c.Status(resp.StatusCode).Send(resp.Body)
}

// CreatePayment godoc
// @Summary      Create Payment Invoice
// @Description  Create a new payment invoice via Xendit
// @Tags         action
// @Accept       json
// @Produce      json
// @Param        X-Api-Key  header    string                   true  "API Key"
// @Param        body          body      map[string]interface{}  true  "Payment invoice payload"
// @Success      200           {object}  map[string]interface{}  "Payment invoice created successfully"
// @Failure      400           {object}  map[string]interface{}  "Invalid request body"
// @Failure      502           {object}  map[string]interface{}  "Failed to process payment"
// @Router       /xendit/action/invoices [post]
func (t *XenditLibController) CreatePayment(c *fiber.Ctx) error {
	return t.handleXenditRequest(c, "/v2/invoices", "Failed to process payment")
}

// CreateSubscription godoc
// @Summary      Create Subscription Plan
// @Description  Create a new recurring subscription plan via Xendit
// @Tags         action
// @Accept       json
// @Produce      json
// @Param        X-Api-Key  header    string                   true  "API Key"
// @Param        body          body      map[string]interface{}  true  "Subscription plan payload"
// @Success      200           {object}  map[string]interface{}  "Subscription plan created successfully"
// @Failure      400           {object}  map[string]interface{}  "Invalid request body"
// @Failure      502           {object}  map[string]interface{}  "Failed to process subscription"
// @Router       /xendit/action/recurring/plans [post]
func (t *XenditLibController) CreateSubscription(c *fiber.Ctx) error {
	return t.handleXenditRequest(c, "/recurring/plans", "Failed to process subscription")
}

// CreatePayout godoc
// @Summary      Create Payout
// @Description  Create a new payout transaction via Xendit
// @Tags         action
// @Accept       json
// @Produce      json
// @Param        X-Api-Key  header    string                   true  "API Key"
// @Param        body          body      map[string]interface{}  true  "Payout payload"
// @Success      200           {object}  map[string]interface{}  "Payout created successfully"
// @Failure      400           {object}  map[string]interface{}  "Invalid request body"
// @Failure      502           {object}  map[string]interface{}  "Failed to process payout"
// @Router       /xendit/action/payouts [post]
func (t *XenditLibController) CreatePayout(c *fiber.Ctx) error {
	return t.handleXenditRequest(c, "/v2/payouts", "Failed to process payout")
}

// CreateCustomer godoc
// @Summary      Create Customer
// @Description  Create a new customer via Xendit
// @Tags         action
// @Accept       json
// @Produce      json
// @Param        X-Api-Key  header    string                   true  "API Key"
// @Param        body          body      map[string]interface{}  true  "Customer payload"
// @Success      200           {object}  map[string]interface{}  "Customer created successfully"
// @Failure      400           {object}  map[string]interface{}  "Invalid request body"
// @Failure      502           {object}  map[string]interface{}  "Failed to process customer"
// @Router       /xendit/action/customers [post]
func (t *XenditLibController) CreateCustomer(c *fiber.Ctx) error {
	return t.handleXenditRequest(c, "/customers", "Failed to process customer")
}
