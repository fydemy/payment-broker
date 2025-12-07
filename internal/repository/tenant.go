package repository

import (
	model "payment-broker/internal/model/db"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CheckAPIKey(APIKey string) (*string, error)
}

type tenantRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewTenantRepository(logger *zap.Logger, db *gorm.DB) TenantRepository {
	return &tenantRepository{
		logger: logger,
		db:     db,
	}
}

func (r *tenantRepository) CheckAPIKey(APIKey string) (*string, error) {
	var tenant model.Tenant
	err := r.db.Select("account_id").Where("api_key = ?", APIKey).First(&tenant).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("tenantRepository.CheckAPIKey", zap.String("api_key", APIKey), zap.Error(err))
			return nil, nil
		}
		r.logger.Error("tenantRepository.CheckAPIKey", zap.String("api_key", APIKey), zap.Error(err))
		return nil, err
	}

	r.logger.Info("tenantRepository.CheckAPIKey", zap.String("api_key", APIKey), zap.String("account_id", tenant.AccountID))
	return &tenant.AccountID, nil
}
