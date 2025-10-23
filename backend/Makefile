# Makefile for Python backend services
.PHONY: help install dev test clean build up down logs migrate

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

migrate: ## Run database migrations
	cd services/user-service && alembic upgrade head

create-migration: ## Create new migration (usage: make create-migration name=migration_name)
	cd services/user-service && alembic revision --autogenerate -m "$(name)"

user-service: ## Start only user service
	docker-compose up user-service

plaid-service: ## Start only plaid service
	docker-compose up plaid-service

transaction-service: ## Start only transaction service
	docker-compose up transaction-service

seed: ## Seed database with test data
	./seed.sh