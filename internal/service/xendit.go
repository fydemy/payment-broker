package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"payment-broker/internal/model/dto"
	"payment-broker/internal/repository"
	"strconv"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type XenditService interface {
	CreateRequest(ctx context.Context, body interface{}, url string, accountID string) (*dto.XenditResponse, error)
	ProxyWebhook(ctx context.Context, tenantID string, body interface{}, api_key string) error
}

type xenditService struct {
	resty            *resty.Client
	logger           *zap.Logger
	splitRuleID      string
	tenantRepository repository.TenantRepository
	baseURL          string
}

func NewXenditService(resty *resty.Client, logger *zap.Logger) XenditService {
	return &xenditService{
		resty:       resty,
		logger:      logger,
		splitRuleID: os.Getenv("XENDIT_SPLIT_RULE_ID"),
		baseURL:     os.Getenv("XENDIT_BASE_URL"),
	}
}

func (s *xenditService) CreateRequest(ctx context.Context, body interface{}, url string, accountID string) (*dto.XenditResponse, error) {
	resp, err := s.resty.R().
		SetBasicAuth(os.Getenv("XENDIT_API_KEY"), "").
		SetHeader("for-user-id", accountID).
		SetHeader("with-split-rule", s.splitRuleID).
		SetBody(body).
		Post(s.baseURL + url)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to reach Xendit API %s", url), zap.Error(err))
		return nil, fmt.Errorf("failed to reach Xendit API %s: %w", url, err)
	}

	return &dto.XenditResponse{
		StatusCode: resp.StatusCode(),
		Header:     resp.Header(),
		Body:       resp.Body(),
	}, nil
}

func (s *xenditService) ProxyWebhook(ctx context.Context, tenantID string, body interface{}, api_key string) error {
	id, err := strconv.ParseUint(tenantID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid tenant ID format: %w", err)
	}

	webhook_url, err := s.tenantRepository.CheckTenant(uint(id))

	if err != nil {
		s.logger.Error("tenantRepository.CheckTenant", zap.String("tenantID", tenantID), zap.Error(err))
		return fmt.Errorf("failed to check tenant: %w", err)
	}

	if webhook_url == nil {
		return fmt.Errorf("tenant not registered yet")
	}

	resp, err := s.resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Api-Key", api_key).
		SetBody(body).
		Post(*webhook_url)

	if err != nil {
		s.logger.Error("resty.Post", zap.String("webhook_url", *webhook_url), zap.Error(err))
		return fmt.Errorf("failed to post to webhook: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		s.logger.Error("resty.Post", zap.String("webhook_url", *webhook_url), zap.Int("status_code", resp.StatusCode()))
		return fmt.Errorf("failed to post to webhook: %w", err)
	}

	return nil
}
