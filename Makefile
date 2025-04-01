# Untether Makefile
# This Makefile provides commands for building, testing, and managing the Untether microservices.

# Load environment variables from .env file
include .env
export

# Declare all phony targets
.PHONY: all build test proto clean docker-up docker-down migrate-up migrate-down \
	build-services docker-build \
	proto-services \
	test-services \
	logs-services \
	db-create db-drop db-reset db-query \
	lint fmt fmt-check check-deps \
	clean clean-docker \
	dev-setup dev-cleanup

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOTIDY=$(GOCMD) mod tidy

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
	$(MAKE) proto-services
	$(MAKE) docker-build

dev-cleanup:
	@echo "Cleaning up development environment..."
	$(MAKE) docker-down
	$(GOCLEAN)
	rm -rf user-service plaid-service roundup-service
	rm -rf services/*/proto/*.pb.go
	@echo "Development environment cleaned up"

# Build Commands
# -------------
build: build-services

build-services:
	@echo "Building all services..."
	$(GOBUILD) -o user-service ./services/user/cmd
	$(GOBUILD) -o plaid-service ./services/plaid/cmd
	$(GOBUILD) -o roundup-service ./services/roundup/cmd

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

docker-build:
	@echo "Building Docker images..."
	docker build -t untether-user-service -f services/user/Dockerfile .
	docker build -t untether-plaid-service -f services/plaid/Dockerfile .
	docker build -t untether-roundup-service -f services/roundup/Dockerfile .

# Service Logs
# -----------
logs-services:
	@echo "Viewing service logs..."
	docker compose logs -f

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
test: test-services

test-services:
	@echo "Running all service tests..."
	$(GOTEST) -v ./services/...

# Proto Commands
# -------------
proto: proto-services

proto-services:
	@echo "Generating all service protos..."
	$(PROTOC) -I./services/user/proto --go_out=./services/user/proto --go_opt=paths=source_relative \
		--go-grpc_out=./services/user/proto --go-grpc_opt=paths=source_relative \
		./services/user/proto/user.proto
	$(PROTOC) -I./services/plaid/proto --go_out=./services/plaid/proto --go_opt=paths=source_relative \
		--go-grpc_out=./services/plaid/proto --go-grpc_opt=paths=source_relative \
		./services/plaid/proto/plaid.proto
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

# Cleanup Commands
# --------------
clean: clean-docker

clean-docker:
	@echo "Cleaning Docker resources..."
	docker compose down -v
	docker system prune -f

# Help
# ----
help:
	@echo "Available commands:"
	@echo "  make all              - Run tests and build all services"
	@echo "  make dev-setup        - Set up development environment"
	@echo "  make dev-cleanup      - Clean up development environment"
	@echo "  make build           - Build all services"
	@echo "  make docker-up       - Start Docker services"
	@echo "  make docker-down     - Stop Docker services"
	@echo "  make logs-services   - View all service logs"
	@echo "  make db-reset        - Reset database"
	@echo "  make test           - Run all tests"
	@echo "  make proto          - Generate all protos"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make clean          - Clean up resources" 