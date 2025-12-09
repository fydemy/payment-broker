package controller

import (
	"encoding/json"
	"payment-broker/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WebhookController interface {
	WebhookHandler(c *fiber.Ctx) error
}

type webhookController struct {
	logger        *zap.Logger
	xenditService service.XenditService
}

func NewWebhookController(logger *zap.Logger, xenditService service.XenditService) WebhookController {
	return &webhookController{
		logger:        logger,
		xenditService: xenditService,
	}
}

// WebhookHandler godoc
// @Summary      Handle Xendit Webhook
// @Description  Handling Event and UnEvent Webhook from Xendit
// @Tags         webhook
// @Accept       json
// @Produce      json
// @Param        X-Api-Key  header    string                   true  "API Key"
// @Param        body       body      map[string]interface{}   true  "Webhook payload"
// @Success      200        {object}  map[string]interface{}   "Webhook processed successfully"
// @Failure      400        {object}  map[string]interface{}   "Invalid webhook body or missing tenant"
// @Failure      500        {object}  map[string]interface{}   "Internal server error"
// @Router       /xendit/webhook [post]
func (t *webhookController) WebhookHandler(c *fiber.Ctx) error {
	rawBody := c.Body()
	apiKey := c.Get("X-Api-Key")

	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		t.logger.Error("Failed to parse webhook body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid webhook body",
		})
	}

	var tenantID string

	if _, hasEvent := body["event"]; hasEvent {
		if data, ok := body["data"].(map[string]interface{}); ok {
			if referenceID, ok := data["reference_id"].(string); ok && referenceID != "" {
				parts := strings.Split(referenceID, ":")
				if len(parts) > 0 {
					tenantID = parts[0]
				}
			}
		}
	} else {
		if referenceID, ok := body["reference_id"].(string); ok && referenceID != "" {
			parts := strings.Split(referenceID, ":")
			if len(parts) > 0 {
				tenantID = parts[0]
			}
		}
	}

	if tenantID == "" {
		t.logger.Error("Failed to extract tenant_id from webhook", zap.Any("body", body))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing tenant identifier in webhook",
		})
	}

	return t.xenditService.ProxyWebhook(c.Context(), tenantID, rawBody, apiKey)
}
