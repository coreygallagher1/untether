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

## Project Structure

```
untether/
├── cmd/                    # Application entrypoints
├── internal/              # Private application code
│   ├── proto/            # Generated protobuf code
│   └── user/             # User service implementation
├── migrations/           # Database migrations
├── proto/               # Protobuf definitions
├── configs/             # Configuration files
├── docs/               # Project documentation
└── scripts/            # Utility scripts
```

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Protocol Buffers compiler (protoc)
- Make

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

3. Start the development environment:
   ```bash
   make dev-setup
   ```
   This will:
   - Install dependencies
   - Generate protobuf code
   - Start Docker services
   - Run database migrations

4. Build and run the service:
   ```bash
   make build
   ./untether
   ```

## Development

### Available Make Commands

- `make build` - Build the application
- `make test` - Run tests
- `make proto` - Generate protobuf code
- `make docker-up` - Start Docker services
- `make docker-down` - Stop Docker services
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback database migrations
- `make dev-setup` - Set up development environment
- `make dev-cleanup` - Clean up development environment

### Running Tests

```bash
make test
```

# Run with coverage report
make test-coverage

### Database Migrations

To create a new migration:
```bash
migrate create -ext sql -dir migrations -seq create_new_table
```

To run migrations:
```bash
make migrate-up
```

To rollback migrations:
```bash
make migrate-down
```

## API Documentation

The service uses gRPC for communication. You can test the endpoints using `grpcurl`:

```bash
# Create a user
grpcurl -plaintext -d '{"email": "test@example.com", "first_name": "Test", "last_name": "User"}' localhost:50051 proto.UserService/CreateUser

# Get a user
grpcurl -plaintext -d '{"id": "user-id"}' localhost:50051 proto.UserService/GetUser
```
