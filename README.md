# Untether

## Mission Statement

At untether, we empower you to break free from the constraints of debt. Through innovative automation, our products simplify financial management. We're committed to helping you untether yourself from debt and move toward a future of financial freedom.

## Overview

Untether is a financial automation platform that helps users manage their finances and break free from debt. The platform provides services for user management, transaction processing, and automated roundup savings.

## Project Structure

```
untether/
├── cmd/                    # Application entrypoints
├── internal/              # Private application code
│   ├── proto/            # Generated protobuf code
│   ├── roundup/          # Roundup service implementation
│   ├── service/          # Service implementations
│   ├── transaction/      # Transaction service implementation
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
