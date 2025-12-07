package service

import (
	"context"
	"fmt"
	"os"
	"payment-broker/internal/model/dto"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type XenditService interface {
	CreatePayment(ctx context.Context, body interface{}, accountID string) (*dto.XenditResponse, error)
}

type xenditService struct {
	restry      *resty.Client
	logger      *zap.Logger
	splitRuleID string
	baseURL     string
}

func NewXenditService(restry *resty.Client, logger *zap.Logger) XenditService {
	return &xenditService{
		restry:      restry,
		logger:      logger,
		splitRuleID: os.Getenv("XENDIT_SPLIT_RULE_ID"),
		baseURL:     os.Getenv("XENDIT_BASE_URL"),
	}
}

func (s *xenditService) CreatePayment(ctx context.Context, body interface{}, accountID string) (*dto.XenditResponse, error) {
	resp, err := s.restry.R().
		SetBasicAuth(os.Getenv("XENDIT_API_KEY"), "").
		SetHeader("for-user-id", accountID).
		SetHeader("with-split-rule", s.splitRuleID).
		SetBody(body).
		Post(s.baseURL + "/v2/invoices")

	if err != nil {
		s.logger.Error("failed to reach Xendit API", zap.Error(err))
		return nil, fmt.Errorf("failed to reach Xendit API: %w", err)
	}

	return &dto.XenditResponse{
		StatusCode: resp.StatusCode(),
		Header:     resp.Header(),
		Body:       resp.Body(),
	}, nil
}
