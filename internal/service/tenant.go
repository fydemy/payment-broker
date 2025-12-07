package service

import (
	"context"
	"errors"
	"fmt"
	"payment-broker/internal/lib"
	"payment-broker/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TenantService interface {
	CheckAPIKey(ctx context.Context, APIKey string) (*string, error)
}

type tenantService struct {
	logger           *zap.Logger
	redisLib         lib.RedisLib
	tenantRepository repository.TenantRepository
}

func NewTenantService(logger *zap.Logger, redisLib lib.RedisLib,
	tenantRepository repository.TenantRepository) TenantService {
	return &tenantService{
		logger:           logger,
		redisLib:         redisLib,
		tenantRepository: tenantRepository,
	}
}

func (s *tenantService) CheckAPIKey(ctx context.Context, APIKey string) (*string, error) {
	cached_accountID, err := s.redisLib.Get(ctx, APIKey)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			accountID, err := s.tenantRepository.CheckAPIKey(APIKey)

			if err != nil {
				s.logger.Error("tenantRepository.CheckAPIKey", zap.String("api_key", APIKey), zap.Error(err))
				return nil, fmt.Errorf("failed to check API key: %w", err)
			}

			if accountID == nil {
				return nil, fmt.Errorf("api_key not registered yet")
			}

			if err := s.redisLib.Set(ctx, APIKey, *accountID, 24*time.Hour); err != nil {
				s.logger.Error("redisLib.Set", zap.String("api_key", APIKey), zap.Error(err))
			}

			return accountID, nil
		}

		s.logger.Error("redisLib.Get", zap.String("api_key", APIKey), zap.Error(err))
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	return &cached_accountID, nil
}
