# Payment Broker

A multi-tenant payment broker service that integrates with XenPlatform feature.

## Features

- Multi-tenant support with API key authentication
- Redis caching for improved performance
- Rate limiting middleware
- Payment and Webhook processing via Xendit API

## Tech Stack

- **Go** - Programming language
- **Fiber v2** - Web framework
- **PostgreSQL** - Database
- **Redis** - Caching & rate limiting
- **GORM** - ORM
- **Zap** - Structured logging

## Prerequisites

- Go 1.24.5+
- PostgreSQL
- Redis

## Installation

```bash
git clone https://github.com/fydemy/payment-broker.git
cd broker
go mod download
cp .env.example .env
```

## Project Structure

```
├── cmd/
│   └── main.go           # Application entry point
│   └── cli
│       └── main.go       # Interactive CLI entry point
├── docs                  # Swagger API docs
├── internal/
│   ├── app/              # Application initialization
│   ├── controller/       # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── service/          # Business logic
│   ├── repository/       # Data access layer
│   ├── model/            # Data models
│   ├── lib/              # Utility libraries
│   └── router/           # Route definitions
├── .env.example          # Example environment variables
├── Makefile              # Automation script for running
└── go.mod                # Go module definition
```

## API Endpoints

- [Create customer](https://docs.xendit.co/apidocs/create-customer-request)
- [Create invoice](https://archive.developers.xendit.co/api-reference/#invoices)
- [Create payout](https://docs.xendit.co/apidocs/create-payout)
- [Create subscription](https://docs.xendit.co/apidocs/create-recurring-plan)
  Webhook
