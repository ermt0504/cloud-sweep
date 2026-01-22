# syntax=docker/dockerfile:1

# =============================================================================
# Stage 1: Build
# =============================================================================
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -o /bin/api ./cmd/api

# Build Worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -o /bin/worker ./cmd/worker

# =============================================================================
# Stage 2: API Runtime
# =============================================================================
FROM alpine:3.19 AS api

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 cloudsweep && \
    adduser -u 1000 -G cloudsweep -s /bin/sh -D cloudsweep

WORKDIR /app

COPY --from=builder /bin/api /app/api
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

USER cloudsweep

EXPOSE 8080

ENTRYPOINT ["/app/api"]

# =============================================================================
# Stage 3: Worker Runtime
# =============================================================================
FROM alpine:3.19 AS worker

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 cloudsweep && \
    adduser -u 1000 -G cloudsweep -s /bin/sh -D cloudsweep

WORKDIR /app

COPY --from=builder /bin/worker /app/worker
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

USER cloudsweep

ENTRYPOINT ["/app/worker"]
