# README for Python FastAPI Backend

## Overview

This is the Python FastAPI backend for the Untether financial application. It consists of microservices for user management, Plaid integration, and transaction processing.

## Architecture

```
backend/
├── shared/                 # Shared utilities and models
│   ├── config.py          # Configuration settings
│   ├── database.py        # Database connection and session
│   ├── auth.py           # Authentication utilities
│   ├── models.py         # SQLAlchemy models
│   ├── schemas.py        # Pydantic schemas
│   └── utils.py          # Common utilities
├── services/
│   ├── user-service/      # User management and auth
│   ├── plaid-service/     # Bank account integration
│   └── transaction-service/ # Transaction processing
├── docker-compose.yml     # Service orchestration
├── requirements.txt       # Python dependencies
└── Makefile             # Development commands
```

## Services

### User Service (Port 8001)
- User registration and authentication
- JWT token management
- User preferences
- Bank account linking

### Plaid Service (Port 8002)
- Bank account connections
- Transaction fetching
- Webhook handling

### Transaction Service (Port 8003)
- Roundup calculations
- Transaction processing
- Analytics

## Quick Start

1. **Set up environment:**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

2. **Start services:**
   ```bash
   make dev
   ```

3. **Verify setup:**
   ```bash
   make status
   ```

4. **Run tests:**
   ```bash
   ./test-setup.sh
   ```

## 📚 Documentation

- **[Complete Setup Guide](docs/COMPLETE_SETUP_GUIDE.md)** - Detailed setup instructions
- **[Troubleshooting Guide](docs/TROUBLESHOOTING.md)** - Common issues and solutions

## Development

### Running Individual Services

```bash
# Start all services
make dev

# Stop all services
make down

# View logs
make logs

# Check status
make status
```

### Testing

```bash
# Comprehensive test suite
./test-setup.sh

# API testing
./test-api.py

# Quick validation
./validate-setup.sh
```

## API Documentation

Once services are running, visit:
- User Service: http://localhost:8001/docs
- Plaid Service: http://localhost:8002/docs
- Transaction Service: http://localhost:8003/docs

## Environment Variables

See `env.example` for all available configuration options.

## Docker Commands

```bash
# Development environment
make dev          # Start all services
make down         # Stop all services
make logs         # View logs
make status       # Check status

# Production environment
make prod         # Start production services
make down-prod    # Stop production services
make logs-prod    # View production logs
make status-prod  # Check production status

# Build and cleanup
make build        # Build development images
make build-prod   # Build production images
make clean        # Clean up containers and volumes
```
