#!/usr/bin/env bash
set -euo pipefail

echo "Starting Pismo Transaction Routine API..."
echo ""

# Check Docker is available
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker first."
    echo "https://docs.docker.com/get-docker/"
    exit 1
fi

# Build and start
docker compose up --build -d

echo ""
echo " Services started!"
echo "   API:      http://localhost:8080"
echo "   Health:   http://localhost:8080/health"
echo ""
echo "  Useful commands:"
echo "   docker compose logs -f api    # View API logs"
echo "   docker compose down           # Stop all services"
echo "   docker compose down -v        # Stop and remove volumes"
