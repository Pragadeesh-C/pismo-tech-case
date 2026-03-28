#!/usr/bin/env bash
set -euo pipefail

echo "Starting Pismo Transaction Routine API..."
echo ""

# Check Docker installed
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker first."
    echo "https://docs.docker.com/get-docker/"
    exit 1
fi

# Check Docker running
if ! docker info &> /dev/null; then
    echo "Docker is not running. Please start Docker."
    exit 1
fi

# Start services
docker compose up --build -d

echo "Waiting for API to be ready..."

until curl -s http://localhost:8080/health > /dev/null; do
  sleep 1
done

echo ""
echo "Services started!"
echo "  API:      http://localhost:8080"
echo "  Health:   http://localhost:8080/health"
echo "  Swagger:  http://localhost:8080/swagger/index.html"
echo ""
echo "Try it:"
echo "  curl -s -X POST http://localhost:8080/api/v1/accounts \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"document_number\": \"12345678900\"}' | jq"
echo ""
echo "  curl -s http://localhost:8080/api/v1/accounts/1 | jq"
echo ""
echo "  curl -s -X POST http://localhost:8080/api/v1/transactions \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"account_id\": 1, \"operation_type_id\": 4, \"amount\": 120.34}' | jq"
echo ""
echo "Useful commands:"
echo "  docker compose logs -f api"
echo "  docker compose down"
echo "  docker compose down -v"