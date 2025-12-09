package repository

import (
	"fmt"
	model "payment-broker/internal/model/db"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CheckAPIKey(APIKey string) (*model.Tenant, error)
	CheckTenant(tenantID uint) (*string, error)
	Create(tenant *model.Tenant) error
	FindAll() ([]model.Tenant, error)
	FindByID(id uint) (*model.Tenant, error)
	Delete(id uint) error
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

func (r *tenantRepository) Create(tenant *model.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *tenantRepository) FindAll() ([]model.Tenant, error) {
	var tenants []model.Tenant
	err := r.db.Find(&tenants).Error
	return tenants, err
}

func (r *tenantRepository) FindByID(id uint) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.First(&tenant, id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tenant{}, id).Error
}

func (r *tenantRepository) CheckAPIKey(APIKey string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.Select("account_id", "id").Where("api_key = ?", APIKey).First(&tenant).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("tenantRepository.CheckAPIKey", zap.String("api_key", APIKey), zap.Error(err))
			return nil, fmt.Errorf("invalid api key")
		}
		r.logger.Error("tenantRepository.CheckAPIKey", zap.String("api_key", APIKey), zap.Error(err))
		return nil, err
	}

	r.logger.Info("tenantRepository.CheckAPIKey", zap.String("api_key", APIKey), zap.String("account_id", tenant.AccountID))
	return &tenant, nil
}

func (r *tenantRepository) CheckTenant(tenantID uint) (*string, error) {
	var tenant model.Tenant
	err := r.db.Select("webhook_url").Where("id = ?", tenantID).First(&tenant).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("tenantRepository.CheckTenant", zap.Uint("tenant_id", tenantID), zap.Error(err))
			return nil, fmt.Errorf("invalid tenant id")
		}
		r.logger.Error("tenantRepository.CheckTenant", zap.Uint("tenant_id", tenantID), zap.Error(err))
		return nil, err
	}

	r.logger.Info("tenantRepository.CheckTenant", zap.Uint("tenant_id", tenantID), zap.String("webhook_url", tenant.WebhookURL))
	return &tenant.WebhookURL, nil
}
