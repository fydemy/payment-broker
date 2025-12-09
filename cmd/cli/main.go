package main

import (
	"fmt"
	"payment-broker/internal/app"
	"payment-broker/internal/repository"
	"payment-broker/internal/service"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load(".env")

	fmt.Print("🚀 Payment Broker - Tenant Manager")
	fmt.Print("===================================\n")

	logger := zap.NewNop()

	db := app.InitDB(logger)

	tenantRepo := repository.NewTenantRepository(logger, db)
	tenantService := service.NewTenantService(logger, nil, nil, tenantRepo)
	cliService := service.NewCLIService(tenantService)

	cliService.MainMenu()
}
