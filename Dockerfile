# ---- Stage 1: Build ----
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
# Download the dependencies and verify the integrity of the modules.
RUN go mod download && go mod verify

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o api ./cmd/api

# ---- Stage 2: Run ----
FROM alpine:latest

WORKDIR /app

# Install ca-certificates and add a user with a specific UID.
RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 appuser

# Copy the built binary from the builder stage.
COPY --from=builder --chown=appuser:appuser /app/api .

# Expose the port the API server will listen on.
EXPOSE 8080

# Start the API server.
CMD ["./api"]