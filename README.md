# Untether

Untether is a holistic personal finance application featuring automated debt repayment and donations from roundups on daily transactions.

## Project Structure

```
Untether/
├── backend/          # Backend services and API
│   ├── services/     # Microservices
│   │   ├── user/     # User management service
│   │   │   ├── cmd/          # Service entry point
│   │   │   ├── internal/     # Service implementation
│   │   │   ├── proto/        # Protocol buffer definitions
│   │   │   ├── migrations/   # Database migrations
│   │   │   ├── tests/        # Test suites
│   │   │   └── Dockerfile    # Service container definition
│   │   │
│   │   ├── plaid/    # Plaid integration service
│   │   │   ├── cmd/          # Service entry point
│   │   │   ├── internal/     # Service implementation
│   │   │   ├── proto/        # Protocol buffer definitions
│   │   │   ├── migrations/   # Database migrations
│   │   │   ├── tests/        # Test suites
│   │   │   └── Dockerfile    # Service container definition
│   │   │
│   │   └── transaction/ # Transaction processing service
│   │       ├── cmd/          # Service entry point
│   │       ├── internal/     # Service implementation
│   │       ├── proto/        # Protocol buffer definitions
│   │       ├── migrations/   # Database migrations
│   │       ├── tests/        # Test suites
│   │       └── Dockerfile    # Service container definition
│   │
│   ├── pkg/          # Shared packages
│   │
│   ├── configs/      # Configuration files
│   │   └── config.go # Main configuration definitions
│   │
│   ├── docker-compose.yml # Docker orchestration
│   ├── Makefile      # Build and development commands
│   ├── go.mod        # Go module dependencies
│   ├── go.sum        # Go module checksums
│   └── .env          # Environment variables
│
└── frontend/         # Next.js frontend application
    ├── src/          # Source code
    ├── public/       # Static assets
    ├── .next/        # Next.js build output
    ├── node_modules/ # Node.js dependencies
    ├── package.json  # Project dependencies
    ├── tsconfig.json # TypeScript configuration
    ├── next.config.js # Next.js configuration
    ├── next.config.ts # TypeScript Next.js configuration
    ├── postcss.config.mjs # PostCSS configuration
    └── eslint.config.mjs # ESLint configuration
```

## Getting Started

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Set up the development environment:
   ```bash
   make dev-setup
   ```

3. Start the services:
   ```bash
   make docker-up
   ```

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

## Development

- Backend services are written in Go
- Frontend is built with Next.js and Material-UI
- gRPC is used for service-to-service communication
- PostgreSQL is used as the primary database
- Kafka is used for event streaming

## License

MIT

## Mission Statement

**Untether empowers individuals to reclaim their financial freedom through intelligent, automated financial tools.** We believe that managing money shouldn't be a burden—our innovative platform simplifies the process of paying down debt, saving, and supporting meaningful causes. By automating everyday financial decisions, Untether enables users to effortlessly take control of their finances, create lasting habits, and make a positive social impact.

## Overview

**Untether** is a comprehensive financial automation platform designed to help users effortlessly manage their money and break free from debt. Inspired by behavioral finance principles, our platform rounds up daily transactions and automatically allocates the spare change towards personal financial goals—such as loan repayment and charitable giving.

### Key Features
- **Automated Roundups:** Securely link bank accounts to automatically round up daily purchases and direct those roundups toward loans, donations, or savings goals.
- **Debt Repayment:** Seamlessly allocate spare change toward paying off student loans, credit card balances, and other debts.
- **Social Good:** Contribute automatically to selected charities and causes with every purchase, transforming everyday spending into meaningful impact.
- **Secure and Transparent:** Built with robust security measures and transparent processes, ensuring users have peace of mind and visibility into their financial progress.

## Architecture

The project is organized into three main microservices:

1. **User Service** (`services/user/`)
   - Handles user authentication and profile management
   - Manages user preferences and settings
   - Provides user-related data access

2. **Plaid Service** (`services/plaid/`)
   - Integrates with Plaid API for bank account connections
   - Manages bank account linking and authentication
   - Handles transaction data synchronization

3. **Transaction Service** (`services/transaction/`)
   - Manages transaction calculations and roundups
   - Handles savings account operations
   - Processes roundup transfers

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL 15 or later
- Protocol Buffers (protoc)
- golangci-lint
- migrate (database migration tool)
- Plaid API credentials (for development)

## Development Commands

The project provides a comprehensive set of Make commands for development:

### Service Management
- `make docker-up` - Start all services
- `make docker-down` - Stop all services
- `make docker-build` - Build Docker images
- `make logs-user` - View user service logs
- `make logs-plaid` - View plaid service logs
- `make logs-transaction` - View transaction service logs

### Database Operations
- `make migrate-up` - Run database migrations
- `make migrate-down` - Roll back migrations
- `make db-create` - Create database
- `make db-drop` - Drop database
- `make db-reset` - Reset database (drop, create, migrate)

### Testing
- `make test` - Run all tests
- `make test-integration` - Run integration tests
- `make test-coverage` - Generate test coverage report
- `make test-user` - Run user service tests
- `make test-plaid` - Run plaid service tests
- `make test-transaction` - Run transaction service tests

### Code Quality
- `make lint` - Run linter
- `make fmt` - Format code
- `make fmt-check` - Check code formatting
- `make health-check` - Check service health

### Cleanup
- `make clean` - Clean all build artifacts
- `make clean-binaries`