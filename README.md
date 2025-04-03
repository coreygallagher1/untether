# Untether

Untether is a personal finance application that helps users manage their finances through automated roundups and smart transaction categorization.

## Project Structure

```
Untether/
├── backend/           # Backend services and API
│   ├── services/     # Microservices (user, plaid, transaction)
│   ├── proto/        # Protocol buffer definitions
│   ├── pkg/          # Shared packages
│   ├── configs/      # Configuration files
│   ├── docker-compose.yml
│   ├── Makefile
│   └── .env
│
└── frontend/         # Next.js frontend application
    ├── src/          # Source code
    ├── public/       # Static assets
    └── package.json
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
- `make clean-binaries` - Clean service binaries
- `make clean-docker` - Clean Docker resources
- `make clean-coverage` - Clean coverage reports

## Service Ports

- User Service: 50051
- Plaid Service: 50052
- Transaction Service: 50053

## API Documentation

The services use gRPC for communication. You can test the endpoints using `grpcurl`:

```bash
# Create a user
grpcurl -plaintext -d '{"email": "test@example.com", "first_name": "Test", "last_name": "User"}' localhost:50051 proto.UserService/CreateUser

# Get a user
grpcurl -plaintext -d '{"id": "user-id"}' localhost:50051 proto.UserService/GetUser

# Create a Plaid Link token
grpcurl -plaintext -d '{"user_id": "user-id"}' localhost:50052 proto.PlaidService/CreateLinkToken

# Exchange a public token for an access token
grpcurl -plaintext -d '{"public_token": "public-token"}' localhost:50052 proto.PlaidService/ExchangePublicToken

# Get accounts for an access token
grpcurl -plaintext -d '{"access_token": "access-token"}' localhost:50052 proto.PlaidService/GetAccounts

# Get balance for an account
grpcurl -plaintext -d '{"access_token": "access-token", "account_id": "account-id"}' localhost:50052 proto.PlaidService/GetBalance

# Calculate roundup for a transaction
grpcurl -plaintext -d '{"user_id": "user-id", "transaction_amount": 10.50}' localhost:50053 proto.TransactionService/RoundupTransaction
```

## Plaid Integration

The Plaid service provides a secure way to connect bank accounts using Plaid's API. The service supports:

- Creating Link tokens for account linking
- Exchanging public tokens for access tokens
- Retrieving account information
- Getting account balances

For development, use the sandbox environment and test credentials:
- Username: `user_good`
- Password: `pass_good`

For production, make sure to:
1. Use production credentials
2. Enable proper error handling
3. Implement rate limiting
4. Set up webhook notifications
5. Follow Plaid's security best practices
