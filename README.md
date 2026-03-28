# Pismo Tech Case — Transaction Routine API

A Go-based RESTful service for managing **accounts** and **transactions**, built as part of the Pismo technical assessment.

## Tech Stack

- **Language**: Go 1.26
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL 18 (via [pgx](https://github.com/jackc/pgx) connection pool)
- **Logging**: [zerolog](https://github.com/rs/zerolog)
- **Containerisation**: Docker + Docker Compose

## Quick Start

### Using Docker Compose (recommended)

```bash
git clone https://github.com/pragadeesh-c/pismo-tech-case.git
cd pismo-tech-case
./run.sh
```

This starts the API on `http://localhost:8080` with PostgreSQL.

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

## API Endpoints
| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/accounts` | Create a new account |
| `GET`  | `/health` | Health check |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | — (required) | PostgreSQL connection string |
| `SERVER_PORT` | `8080` | HTTP server port |
| `SERVER_GIN_MODE` | `debug` | Gin mode (`debug` / `release`) |
| `LOG_ENV` | `production` | `development` for console, else JSON |
| `LOG_LEVEL` | `info` | Log level (`debug`, `info`, `warn`, `error`) |

## Running Tests

    go test ./...

## Swagger Documentation

Once running, visit: `http://localhost:8080/swagger/index.html`

## Design Decisions

- **Dependency injection**: All layers receive dependencies through constructors, making unit testing straightforward
- **pgx over database/sql**: Better PostgreSQL-native type support and connection pooling
- **Migrations at startup**: Auto-applied via golang-migrate to keep schema in sync
- **Graceful shutdown**: Drains in-flight requests on SIGINT/SIGTERM before exit