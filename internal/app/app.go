package app

import (
	"payment-broker/internal/controller"
	"payment-broker/internal/lib"
	"payment-broker/internal/repository"
	"payment-broker/internal/service"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Repository struct {
		Tenant repository.TenantRepository
	}

	Service struct {
		Tenant service.TenantService
		Xendit service.XenditService
	}

	Controller struct {
		Xendit controller.XenditController
	}
}

func InitApp(db *gorm.DB, logger *zap.Logger, redis *redis.Client) *App {
	app := &App{}

	redisLib := lib.NewRedisLib(redis)
	resty := resty.New().SetTimeout(10 * time.Second)

	app.Repository.Tenant = repository.NewTenantRepository(logger, db)
	app.Service.Tenant = service.NewTenantService(logger, redisLib, app.Repository.Tenant)
	app.Service.Xendit = service.NewXenditService(resty, logger)
	app.Controller.Xendit = controller.NewXenditController(logger, app.Service.Xendit)

	return app
}
