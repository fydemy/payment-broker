# Payment Broker

A multi-tenant payment broker service that integrates with Xendit payment gateway.

## Features

- Multi-tenant support with API key authentication
- Redis caching for improved performance
- Rate limiting middleware
- Payment processing via Xendit API

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

1. Clone the repository

```bash
git clone https://github.com/fydemy/payment-broker.git
cd broker
```

2. Install dependencies

```bash
go mod download
```

3. Set up environment variables

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
REDIS_ADDR=localhost:6379
REDIS_PWD=

DB_DSN="host=localhost user=postgres dbname=api port=5432"

APP_ENV=development
APP_PORT=8080

XENDIT_API_KEY=your_xendit_api_key
```

## Running the Application

From the project root:

```bash
go run cmd/main.go
```

Or use Air for hot reload:

```bash
cd cmd
air
```

## Project Structure

```
.
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── app/              # Application initialization
│   ├── controller/       # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── service/          # Business logic
│   ├── repository/       # Data access layer
│   ├── model/            # Data models
│   ├── lib/              # Utility libraries
│   └── router/           # Route definitions
├── .env                  # Environment variables (git ignored)
└── go.mod                # Go module definition
```

## API Endpoints

### Create Payment

```
POST /v1/xendit/invoices
Headers:
  Content-Type: application/json
  X-API-Key: <your-api-key>

Body:
{
  // Xendit invoice payload
}
```
