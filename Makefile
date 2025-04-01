# Untether Makefile
# This Makefile provides commands for building, testing, and managing the Untether microservices.

# Load environment variables from .env file
include .env
export

# Declare all phony targets
.PHONY: all build test proto clean docker-up docker-down migrate-up migrate-down test-integration test-coverage \
	build-user build-plaid build-roundup \
	run-user run-plaid run-roundup \
	proto-user proto-plaid proto-roundup \
	test-user test-plaid test-roundup \
	docker-build docker-build-user docker-build-plaid docker-build-roundup \
	dev-setup logs-user logs-plaid logs-roundup \
	db-create db-drop db-reset db-query \
	lint fmt fmt-check check-deps health-check \
	clean clean-binaries clean-docker clean-coverage

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOTIDY=$(GOCMD) mod tidy

# Binary names
USER_BINARY=user-service
PLAID_BINARY=plaid-service
ROUNDUP_BINARY=roundup-service

# Proto commands
PROTOC=protoc

# Migration commands
MIGRATE=migrate
MIGRATE_DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Service ports
USER_PORT=50051
PLAID_PORT=50052
ROUNDUP_PORT=50053

# Default target
all: test build

# Development Setup
# -----------------
dev-setup: check-deps
	@echo "Setting up development environment..."
	$(GOTIDY)
	$(MAKE) proto
	$(MAKE) docker-build

# Build Commands
# -------------
build: build-user build-plaid build-roundup

build-user:
	@echo "Building user service..."
	$(GOBUILD) -o $(USER_BINARY) ./services/user/cmd

build-plaid:
	@echo "Building plaid service..."
	$(GOBUILD) -o $(PLAID_BINARY) ./services/plaid/cmd

build-roundup:
	@echo "Building roundup service..."
	$(GOBUILD) -o $(ROUNDUP_BINARY) ./services/roundup/cmd

# Docker Commands
# --------------
docker-up:
	@echo "Starting Docker services..."
	docker compose up -d postgres
	@echo "Waiting for postgres to be ready..."
	@sleep 5
	@echo "Running database migrations..."
	docker compose exec -T postgres migrate -database "$(MIGRATE_DATABASE_URL)" -path /migrations/user up
	docker compose exec -T postgres migrate -database "$(MIGRATE_DATABASE_URL)" -path /migrations/plaid up
	docker compose exec -T postgres migrate -database "$(MIGRATE_DATABASE_URL)" -path /migrations/roundup up
	@echo "Starting remaining services..."
	docker compose up -d

docker-down:
	@echo "Stopping Docker services..."
	docker compose down

docker-build: docker-build-user docker-build-plaid docker-build-roundup

docker-build-user:
	@echo "Building user service Docker image..."
	docker build -t untether-user-service -f services/user/Dockerfile .

docker-build-plaid:
	@echo "Building plaid service Docker image..."
	docker build -t untether-plaid-service -f services/plaid/Dockerfile .

docker-build-roundup:
	@echo "Building roundup service Docker image..."
	docker build -t untether-roundup-service -f services/roundup/Dockerfile .

# Service Logs
# -----------
logs-user:
	@echo "Viewing user service logs..."
	docker compose logs -f user-service

logs-plaid:
	@echo "Viewing plaid service logs..."
	docker compose logs -f plaid-service

logs-roundup:
	@echo "Viewing roundup service logs..."
	docker compose logs -f roundup-service

# Database Commands
# ----------------
migrate-up:
	@echo "Running database migrations..."
	docker compose exec -T postgres migrate -database "$(MIGRATE_DATABASE_URL)" -path /migrations/user up
	docker compose exec -T postgres migrate -database "$(MIGRATE_DATABASE_URL)" -path /migrations/plaid up
	docker compose exec -T postgres migrate -database "$(MIGRATE_DATABASE_URL)" -path /migrations/roundup up

migrate-down:
	@echo "Rolling back database migrations..."
	docker compose exec -T postgres migrate -database "$(MIGRATE_DATABASE_URL)" -path /migrations/user down
	docker compose exec -T postgres migrate -database "$(MIGRATE_DATABASE_URL)" -path /migrations/plaid down
	docker compose exec -T postgres migrate -database "$(MIGRATE_DATABASE_URL)" -path /migrations/roundup down

db-create:
	@echo "Creating database..."
	docker compose exec postgres createdb -U $(DB_USER) $(DB_NAME)

db-drop:
	@echo "Dropping database..."
	docker compose exec postgres dropdb -U $(DB_USER) $(DB_NAME)

db-reset: db-drop db-create migrate-up

db-query:
	@echo "Executing database query..."
	docker compose exec postgres psql -U $(DB_USER) -d $(DB_NAME) -c "$(QUERY)"

# Test Commands
# ------------
test: test-user test-plaid test-roundup

test-user:
	@echo "Running user service tests..."
	$(GOTEST) -v ./services/user/...

test-plaid:
	@echo "Running plaid service tests..."
	$(GOTEST) -v ./services/plaid/...

test-roundup:
	@echo "Running roundup service tests..."
	$(GOTEST) -v ./services/roundup/...

test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./...

test-coverage:
	@echo "Generating test coverage report..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Proto Commands
# -------------
proto: proto-user proto-plaid proto-roundup

proto-user:
	@echo "Generating user service protos..."
	$(PROTOC) -I./services/user/proto --go_out=./services/user/proto --go_opt=paths=source_relative \
		--go-grpc_out=./services/user/proto --go-grpc_opt=paths=source_relative \
		./services/user/proto/user.proto

proto-plaid:
	@echo "Generating plaid service protos..."
	$(PROTOC) -I./services/plaid/proto --go_out=./services/plaid/proto --go_opt=paths=source_relative \
		--go-grpc_out=./services/plaid/proto --go-grpc_opt=paths=source_relative \
		./services/plaid/proto/plaid.proto

proto-roundup:
	@echo "Generating roundup service protos..."
	$(PROTOC) -I./services/roundup/proto --go_out=./services/roundup/proto --go_opt=paths=source_relative \
		--go-grpc_out=./services/roundup/proto --go-grpc_opt=paths=source_relative \
		./services/roundup/proto/roundup.proto

# Code Quality
# -----------
lint:
	@echo "Running linter..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

fmt-check:
	@echo "Checking code formatting..."
	@if [ -n "$$(go fmt ./...)" ]; then echo "Code is not formatted"; exit 1; fi

# Dependency Checks
# ----------------
check-deps:
	@echo "Checking dependencies..."
	@command -v protoc >/dev/null 2>&1 || { echo "protoc is required but not installed. Aborting." >&2; exit 1; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is required but not installed. Aborting." >&2; exit 1; }
	@command -v migrate >/dev/null 2>&1 || { echo "migrate is required but not installed. Aborting." >&2; exit 1; }

# Health Checks
# -----------
health-check:
	@echo "Checking service health..."
	@curl -s http://localhost:$(USER_PORT)/health || echo "User service not healthy"
	@curl -s http://localhost:$(PLAID_PORT)/health || echo "Plaid service not healthy"
	@curl -s http://localhost:$(ROUNDUP_PORT)/health || echo "Roundup service not healthy"

# Cleanup Commands
# --------------
clean: clean-binaries clean-docker clean-coverage

clean-binaries:
	@echo "Cleaning binaries..."
	$(GOCLEAN)
	rm -rf $(USER_BINARY) $(PLAID_BINARY) $(ROUNDUP_BINARY)

clean-docker:
	@echo "Cleaning Docker resources..."
	docker compose down -v
	docker system prune -f

clean-coverage:
	@echo "Cleaning coverage reports..."
	rm -rf coverage.out coverage.html

# Help
# ----
help:
	@echo "Available commands:"
	@echo "  make all              - Run tests and build all services"
	@echo "  make dev-setup        - Set up development environment"
	@echo "  make build           - Build all services"
	@echo "  make docker-up       - Start Docker services"
	@echo "  make docker-down     - Stop Docker services"
	@echo "  make docker-build    - Build Docker images"
	@echo "  make logs-*          - View service logs (user/plaid/roundup)"
	@echo "  make migrate-up      - Run database migrations"
	@echo "  make migrate-down    - Roll back migrations"
	@echo "  make db-reset        - Reset database (drop, create, migrate)"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Generate test coverage report"
	@echo "  make proto          - Generate protobuf files"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make health-check   - Check service health"
	@echo "  make clean          - Clean up build artifacts" 