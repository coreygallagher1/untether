# Makefile for Python backend services
.PHONY: help install dev test clean build up down logs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install Python dependencies
	pip install -r requirements.txt

dev: ## Start development environment
	docker compose up --build -d

prod: ## Start production environment
	docker compose -f docker-compose.prod.yml up --build -d

test: ## Run tests
	pytest tests/

clean: ## Clean up containers and volumes
	docker compose down -v
	docker compose -f docker-compose.prod.yml down -v
	docker system prune -f

build: ## Build all services
	docker compose build

build-prod: ## Build production images
	docker compose -f docker-compose.prod.yml build

up: ## Start services
	docker compose up -d

up-prod: ## Start production services
	docker compose -f docker-compose.prod.yml up -d

down: ## Stop services
	docker compose down

down-prod: ## Stop production services
	docker compose -f docker-compose.prod.yml down

logs: ## Show logs
	docker compose logs -f

logs-prod: ## Show production logs
	docker compose -f docker-compose.prod.yml logs -f

status: ## Show service status
	docker compose ps

status-prod: ## Show production service status
	docker compose -f docker-compose.prod.yml ps

init-db: ## Initialize database from models (simple, no migrations)
	@echo "Initializing database from models..."
	@docker compose run --rm user-service python3 -c "import sys; sys.path.insert(0, '/app'); from shared.database import init_db; from shared.config import settings; print('Initializing database:', settings.database_url); print('Creating all tables from models...'); init_db(); print('✅ Database initialized successfully!'); print(''); print('Tables created:'); print('  - users'); print('  - user_preferences'); print('  - plaid_items'); print('  - bank_accounts'); print('  - roundup_calculations')"

user-service: ## Start only user service
	docker-compose up user-service

plaid-service: ## Start only plaid service
	docker-compose up plaid-service

transaction-service: ## Start only transaction service
	docker-compose up transaction-service

seed: ## Seed database with test data
	./seed.sh