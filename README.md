# Untether

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

3. **Roundup Service** (`services/roundup/`)
   - Manages roundup transaction calculations
   - Handles savings account operations
   - Processes roundup transfers

## Project Structure

```
.
├── services/
│   ├── user/
│   │   ├── cmd/           # Service entry point
│   │   ├── internal/      # Internal service code
│   │   ├── migrations/    # Database migrations
│   │   └── proto/        # Protocol buffer definitions
│   ├── plaid/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── migrations/
│   │   └── proto/
│   └── roundup/
│       ├── cmd/
│       ├── internal/
│       ├── migrations/
│       └── proto/
├── pkg/                   # Shared packages
├── Makefile              # Build and development commands
└── docker-compose.yml    # Docker service definitions
```

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL 15 or later
- Protocol Buffers (protoc)
- golangci-lint
- migrate (database migration tool)
- Plaid API credentials (for development)

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/untether.git
   cd untether
   ```

2. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. Set up Plaid credentials:
   - Go to [Plaid Dashboard](https://dashboard.plaid.com/team/keys)
   - Create a new application or select an existing one
   - Copy your Client ID and Client Secret
   - Update your `.env` file with the credentials:
     ```
     PLAID_CLIENT_ID=your_client_id_here
     PLAID_CLIENT_SECRET=your_client_secret_here
     PLAID_ENVIRONMENT=sandbox  # Use sandbox for development
     ```

4. Set up the development environment:
   ```bash
   make dev-setup
   ```
   This will:
   - Check for required dependencies
   - Generate Protocol Buffer code
   - Build Docker images
   - Start required services

5. Run database migrations:
   ```bash
   make migrate-up
   ```

## Development Commands

The project provides a comprehensive set of Make commands for development:

### Service Management
- `make docker-up` - Start all services
- `make docker-down` - Stop all services
- `make docker-build` - Build Docker images
- `make logs-user` - View user service logs
- `make logs-plaid` - View plaid service logs
- `make logs-roundup` - View roundup service logs

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
- `make test-roundup` - Run roundup service tests

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
- Roundup Service: 50053

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
grpcurl -plaintext -d '{"user_id": "user-id", "transaction_amount": 10.50}' localhost:50053 proto.RoundupService/RoundupTransaction
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
