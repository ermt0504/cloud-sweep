.PHONY: build build-api build-worker run-api run-worker test lint clean deps docker-up docker-down docker-build migrate swagger

# Variables
BINARY_API=bin/api
BINARY_WORKER=bin/worker
GO=go
GOFLAGS=-ldflags="-s -w"

# Build
build: build-api build-worker

build-api:
	$(GO) build $(GOFLAGS) -o $(BINARY_API) ./cmd/api

build-worker:
	$(GO) build $(GOFLAGS) -o $(BINARY_WORKER) ./cmd/worker

# Run
run-api:
	$(GO) run ./cmd/api

run-worker:
	$(GO) run ./cmd/worker

# Test
test:
	$(GO) test -v -race -cover ./...

test-coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	golangci-lint run ./...

# Dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker build -t cloudsweep-api --target api .
	docker build -t cloudsweep-worker --target worker .

# Database
migrate:
	$(GO) run ./cmd/migrate

migrate-down:
	$(GO) run ./cmd/migrate down

# Swagger
swagger:
	swag init -g docs/swagger.go -o docs --parseDependency --parseInternal

swagger-install:
	$(GO) install github.com/swaggo/swag/cmd/swag@latest

swagger-fmt:
	swag fmt

# Development
dev-api:
	air -c .air.api.toml

dev-worker:
	air -c .air.worker.toml

# Help
help:
	@echo "Commandes disponibles:"
	@echo "  make build          - Compile tous les binaires"
	@echo "  make run-api        - Lance l'API"
	@echo "  make run-worker     - Lance le worker"
	@echo "  make test           - Execute les tests"
	@echo "  make lint           - Analyse statique"
	@echo "  make deps           - Telecharge les dependances"
	@echo "  make docker-up      - Demarre les conteneurs"
	@echo "  make docker-down    - Arrete les conteneurs"
	@echo "  make clean          - Nettoie les artefacts"
	@echo "  make swagger        - Genere la documentation Swagger"
	@echo "  make swagger-install - Installe swag CLI"
