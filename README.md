# Pismo Tech Case — Transaction Routine API

## Overview

This service implements a simplified transaction system, handling account creation and financial transactions with proper validation, error handling, and clean architecture principles.

---

## Features

* Create and retrieve accounts
* Create transactions with operation type handling (debit/credit)
* Input validation and domain-level error handling
* Consistent API error responses
* Request logging with structured logs
* Swagger API documentation

---

## Tech Stack

* **Language**: Go 1.26
* **Framework**: Gin (https://github.com/gin-gonic/gin)
* **Database**: PostgreSQL (via [pgx](https://github.com/jackc/pgx) connection pool)
* **SQL Code Generation**: sqlc 
* **Logging**: zerolog (https://github.com/rs/zerolog)
* **Containerisation**: Docker + Docker Compose

---

## API Documentation

- Swagger (OpenAPI)

---

## Architecture

The application follows a layered architecture:

```
Handler → Service → Repository → Database
```

* **Handler**: HTTP layer (Gin)
* **Service**: Business logic and domain rules
* **Repository**: Database access via sqlc

### Project Structure

```
├── cmd/
│   ├── api/             # Application entrypoint
│   └── docs/            # Auto-generated Swagger docs
├── internal/
│   ├── api/
│   │   ├── handler/     # HTTP handlers (request/response)
│   │   ├── middleware/   # CORS, request logging
│   │   ├── models/      # Request/response DTOs, error codes
│   │   └── route/       # Route registration
│   ├── config/          # Environment config loading
│   ├── constants/       # Operation types, shared constants
│   ├── database/        # Connection pool, migrations
│   ├── mocks/           # Test mocks (repository)
│   ├── repository/      # sqlc-generated data access
│   └── service/         # Business logic, domain models
└── db/
    ├── migrations/      # SQL migration files
    └── queries/         # sqlc query definitions
```

---

## Quick Start

### Using Docker Compose (recommended)

```bash
git clone https://github.com/pragadeesh-c/pismo-tech-case.git
cd pismo-tech-case
./run.sh
```

This starts the API on `http://localhost:8080` with PostgreSQL.

---

### Running Locally

**Prerequisites**: Go 1.26+, PostgreSQL running locally

```bash
# 1. Set up environment
cp .env.example .env
# Edit .env with your database credentials

# 2. Install dependencies
go mod tidy

# 3. Run the service
go run ./cmd/api
```

---

## API Endpoints

All endpoints are prefixed with `/api/v1` except for /health

| Method | Path                        | Description              |
| ------ | --------------------------- | ------------------------ |
| POST   | /api/v1/accounts            | Create a new account     |
| GET    | /api/v1/accounts/:accountId | Get account by ID        |
| POST   | /api/v1/transactions        | Create a new transaction |
| GET    | /health                     | Health check             |

---

## Environment Variables

| Variable        | Default      | Description                          |
| --------------- | ------------ | ------------------------------------ |
| DATABASE_URL    | — (required) | PostgreSQL connection string         |
| SERVER_PORT     | 8080         | HTTP server port                     |
| SERVER_GIN_MODE | debug        | Gin mode (debug / release)           |
| LOG_ENV         | production   | development for console, else JSON   |
| LOG_LEVEL       | info         | Log level (debug, info, warn, error) |

---

## Running Tests

```bash
go test ./...
```

* Unit tests cover service layer with mocked repository
---

## Swagger Documentation

Once running, visit:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Development Commands

### Generate sqlc code

```bash
sqlc generate
```

Regenerates the type-safe Go code in `internal/repository/` from the SQL queries in `db/queries/`.

### Generate Swagger docs

```bash
swag init -g ./cmd/api/main.go -o cmd/docs
```

Regenerates the Swagger spec from handler annotations into `cmd/docs/`.

### Database migrations

Migrations are auto-applied at startup via `golang-migrate`. To manage them manually:

```bash
# Create a new migration
migrate create -ext sql -dir db/migrations -seq <migration_name>

# Apply all pending migrations
migrate -path db/migrations -database "$DATABASE_URL" up

# Rollback the last migration
migrate -path db/migrations -database "$DATABASE_URL" down 1
```

---

## Design Decisions

* **Dependency injection**: All layers receive dependencies through constructors, enabling easier testing and separation of concerns
* **pgx over database/sql**: Better PostgreSQL-native type support and connection pooling
* **sqlc for queries**: Type-safe query generation and compile-time validation
* **Monetary values**: Stored using PostgreSQL NUMERIC to avoid floating-point precision issues
* **Migrations at startup**: Auto-applied via golang-migrate to keep schema in sync
* **Graceful shutdown**: Drains in-flight requests on SIGINT/SIGTERM before exit

---

## Future Enhancements

- Implement idempotency for transaction creation to ensure safe retries and prevent duplicate processing
- Improve validation and request schema enforcement
- Add integration tests with real database
- Implement pagination and filtering for transaction history