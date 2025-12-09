package service

import (
	"context"
	"errors"
	"fmt"
	"payment-broker/internal/helper"
	"payment-broker/internal/lib"
	model "payment-broker/internal/model/db"
	"payment-broker/internal/model/dto"
	"payment-broker/internal/repository"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TenantService interface {
	CheckAPIKey(ctx context.Context, APIKey string) (*dto.TenantCheckAPIKey, error)
	CreateTenant(name, accountID, webhookURL string) (*model.Tenant, error)
	GetAllTenants() ([]model.Tenant, error)
	GetTenantByID(id uint) (*model.Tenant, error)
	DeleteTenant(id uint) error
}

type tenantService struct {
	logger           *zap.Logger
	resty            *resty.Client
	redisLib         lib.RedisLib
	tenantRepository repository.TenantRepository
}

func NewTenantService(logger *zap.Logger, resty *resty.Client, redisLib lib.RedisLib,
	tenantRepository repository.TenantRepository) TenantService {
	return &tenantService{
		logger:           logger,
		resty:            resty,
		redisLib:         redisLib,
		tenantRepository: tenantRepository,
	}
}

func (s *tenantService) GetTenantByID(id uint) (*model.Tenant, error) {
	return s.tenantRepository.FindByID(id)
}

func (s *tenantService) GetAllTenants() ([]model.Tenant, error) {
	return s.tenantRepository.FindAll()
}

func (s *tenantService) DeleteTenant(id uint) error {
	return s.tenantRepository.Delete(id)
}

func (s *tenantService) CreateTenant(name, accountID, webhookURL string) (*model.Tenant, error) {
	tenant := &model.Tenant{
		Name:       name,
		AccountID:  accountID,
		WebhookURL: webhookURL,
		APIKey:     helper.GenerateAPIKey(),
	}

	err := s.tenantRepository.Create(tenant)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *tenantService) CheckAPIKey(ctx context.Context, APIKey string) (*dto.TenantCheckAPIKey, error) {
	cached_accountID, err := s.redisLib.Get(ctx, APIKey)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			tenant, err := s.tenantRepository.CheckAPIKey(APIKey)

			if err != nil {
				s.logger.Error("tenantRepository.CheckAPIKey", zap.String("api_key", APIKey), zap.Error(err))
				return nil, fmt.Errorf("failed to check API key: %w", err)
			}

			if tenant.AccountID == "" {
				return nil, fmt.Errorf("api_key not registered yet")
			}

			if err := s.redisLib.Set(ctx, APIKey, fmt.Sprintf("%d:%s", tenant.ID, tenant.AccountID), 24*7*time.Hour); err != nil {
				s.logger.Error("redisLib.Set", zap.String("api_key", APIKey), zap.Error(err))
			}

			id := fmt.Sprintf("%d", tenant.ID)

			return &dto.TenantCheckAPIKey{
				ID:        id,
				AccountID: tenant.AccountID,
			}, nil
		}

		s.logger.Error("redisLib.Get", zap.String("api_key", APIKey), zap.Error(err))
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	return &dto.TenantCheckAPIKey{
		ID:        strings.Split(cached_accountID, ":")[0],
		AccountID: strings.Split(cached_accountID, ":")[1],
	}, nil
}
