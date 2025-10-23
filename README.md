# Untether

Untether is a holistic personal finance application featuring automated debt repayment and donations from roundups on daily transactions.

## Project Structure

```
Untether/
├── backend/                    # Python FastAPI Backend
│   ├── services/              # Microservices
│   │   ├── user-service/      # User management service
│   │   │   ├── main.py        # FastAPI application
│   │   │   ├── Dockerfile     # Service container definition
│   │   │   ├── requirements.txt # Python dependencies
│   │   │   └── migrations/    # Alembic database migrations
│   │   │
│   │   ├── plaid-service/     # Plaid integration service
│   │   │   ├── main.py        # FastAPI application
│   │   │   ├── Dockerfile     # Service container definition
│   │   │   ├── requirements.txt # Python dependencies
│   │   │   └── migrations/    # Alembic database migrations
│   │   │
│   │   └── transaction-service/ # Transaction processing service
│   │       ├── main.py        # FastAPI application
│   │       ├── Dockerfile     # Service container definition
│   │       └── requirements.txt # Python dependencies
│   │
│   ├── shared/                # Shared Python modules
│   │   ├── auth.py           # Authentication utilities
│   │   ├── config.py         # Configuration management
│   │   ├── database.py       # Database connection
│   │   ├── models.py         # SQLAlchemy ORM models
│   │   ├── schemas.py        # Pydantic schemas
│   │   └── utils.py          # Utility functions
│   │
│   ├── docker-compose.yml    # Docker orchestration
│   ├── docker-compose.prod.yml # Production Docker setup
│   ├── Makefile              # Build and development commands
│   ├── requirements.txt      # Python dependencies
│   ├── init-db.sql           # Database initialization
│   ├── env.example           # Environment variables template
│   ├── COMPLETE_SETUP_GUIDE.md # Setup instructions
│   ├── TROUBLESHOOTING.md    # Debug guide
│   ├── test-setup.sh         # Comprehensive testing script
│   ├── test-api.py           # API testing script
│   └── validate-setup.sh     # Quick validation script
│
└── frontend/                  # Next.js Frontend
    ├── src/                   # Source code
    │   ├── api/              # API client
    │   ├── app/              # App router
    │   └── components/       # React components
    ├── public/               # Static assets
    ├── package.json          # Node.js dependencies
    ├── tsconfig.json         # TypeScript configuration
    ├── next.config.js        # Next.js configuration
    ├── next.config.ts        # TypeScript Next.js configuration
    ├── postcss.config.mjs    # PostCSS configuration
    └── eslint.config.mjs     # ESLint configuration
```

## Getting Started

### Prerequisites

- Python 3.11 or later
- Node.js 18 or later
- Docker and Docker Compose
- PostgreSQL 15 or later
- Plaid API credentials (for development)

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Copy environment variables:
   ```bash
   cp env.example .env
   ```

3. Start the services:
   ```bash
   make dev
   ```

4. Verify services are running:
   ```bash
   make status
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

4. Open your browser to `http://localhost:3000`

## Development

- **Backend:** Python FastAPI microservices with SQLAlchemy ORM
- **Frontend:** Next.js with React, TypeScript, and Material-UI
- **Database:** PostgreSQL with Alembic migrations
- **Authentication:** JWT tokens with bcrypt password hashing
- **API:** RESTful APIs with automatic OpenAPI documentation
- **Containerization:** Docker with multi-stage builds

## Architecture

The project is organized into three main microservices:

1. **User Service** (`services/user-service/`)
   - Handles user authentication and profile management
   - Manages user preferences and settings
   - Provides user-related data access
   - JWT token generation and validation

2. **Plaid Service** (`services/plaid-service/`)
   - Integrates with Plaid API for bank account connections
   - Manages bank account linking and authentication
   - Handles transaction data synchronization
   - Webhook processing for real-time updates

3. **Transaction Service** (`services/transaction-service/`)
   - Manages transaction calculations and roundups
   - Handles savings account operations
   - Processes roundup transfers
   - Financial calculations and reporting

## Development Commands

The project provides a comprehensive set of Make commands for development:

### Service Management
- `make dev` - Start all services in development mode
- `make prod` - Start all services in production mode
- `make down` - Stop all services
- `make build` - Build Docker images
- `make logs` - View all service logs
- `make status` - Check service health status

### Database Operations
- `make clean` - Clean up containers and volumes
- Database migrations are handled automatically via Alembic

### Testing
- `./test-setup.sh` - Run comprehensive test suite
- `./test-api.py` - Test API endpoints
- `./validate-setup.sh` - Quick validation check

### Code Quality
- `make build-prod` - Build production images
- `make up-prod` - Start production environment
- `make logs-prod` - View production logs
- `make status-prod` - Check production status

## API Documentation

Once services are running, API documentation is available at:
- User Service: `http://localhost:8001/docs`
- Plaid Service: `http://localhost:8002/docs`
- Transaction Service: `http://localhost:8003/docs`

## Key Features

- **Automated Roundups:** Securely link bank accounts to automatically round up daily purchases and direct those roundups toward loans, donations, or savings goals.
- **Debt Repayment:** Seamlessly allocate spare change toward paying off student loans, credit card balances, and other debts.
- **Social Good:** Contribute automatically to selected charities and causes with every purchase, transforming everyday spending into meaningful impact.
- **Secure and Transparent:** Built with robust security measures and transparent processes, ensuring users have peace of mind and visibility into their financial progress.

## Mission Statement

**Untether empowers individuals to reclaim their financial freedom through intelligent, automated financial tools.** We believe that managing money shouldn't be a burden—our innovative platform simplifies the process of paying down debt, saving, and supporting meaningful causes. By automating everyday financial decisions, Untether enables users to effortlessly take control of their finances, create lasting habits, and make a positive social impact.

## License

MIT