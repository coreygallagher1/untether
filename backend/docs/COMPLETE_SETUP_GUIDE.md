# 🚀 Complete Setup Guide: Python FastAPI Microservices

## Overview
Your Python FastAPI microservices are ready! This guide will help you get everything running and tested.

## ✅ What's Been Created

### **Core Architecture**
- ✅ **3 FastAPI services** (User, Plaid, Transaction)
- ✅ **Shared modules** (config, database, auth, models, schemas, utils)
- ✅ **Docker setup** with multi-stage builds and security
- ✅ **Database models** with proper relationships
- ✅ **Authentication system** with JWT tokens
- ✅ **Comprehensive validation** and error handling

### **Testing & Validation**
- ✅ **Test scripts** for comprehensive validation
- ✅ **API testing** with automated endpoints
- ✅ **Health checks** for all services
- ✅ **Troubleshooting guide** for common issues

## 🚀 Quick Start

### **Step 1: Prerequisites**
Make sure you have:
- Docker installed and running
- Docker Compose installed
- Git (for version control)

### **Step 2: Start Services**
```bash
cd backend

# Start development environment
make dev

# Or manually
docker-compose up --build
```

### **Step 3: Verify Everything Works**
```bash
# Run comprehensive tests
./test-setup.sh

# Or run API tests
python3 test-api.py
```

## 📋 Detailed Setup Instructions

### **1. Environment Setup**
```bash
cd backend

# Validate setup (creates .env if needed)
./validate-setup.sh

# Check Docker is running
docker --version
docker-compose --version
```

### **2. Start Development Environment**
```bash
# Start all services
make dev

# Check status
make status

# View logs
make logs
```

### **3. Test the System**
```bash
# Comprehensive test suite
./test-setup.sh

# API endpoint testing
python3 test-api.py

# Check service health
curl http://localhost:8001/health
curl http://localhost:8002/health
curl http://localhost:8003/health
```

## 🔧 Available Commands

### **Development Commands**
```bash
make dev          # Start development environment
make status       # Check service status
make logs         # View logs
make down         # Stop services
make clean        # Clean up everything
```

### **Production Commands**
```bash
make prod         # Start production environment
make build-prod   # Build production images
make status-prod  # Check production status
make logs-prod    # View production logs
```

### **Database Commands**
```bash
make init-db     # Initialize database from models (simple, no migrations)
# Note: For production migrations, see docs/ADDING_MIGRATIONS_LATER.md
```

## 🌐 Service Endpoints

### **User Service (Port 8001)**
- **Health:** http://localhost:8001/health
- **API Docs:** http://localhost:8001/docs
- **Register:** POST /auth/register
- **Login:** POST /auth/login
- **User Info:** GET /users/me

### **Plaid Service (Port 8002)**
- **Health:** http://localhost:8002/health
- **API Docs:** http://localhost:8002/docs
- **Link Token:** POST /plaid/link-token
- **Accounts:** GET /plaid/accounts

### **Transaction Service (Port 8003)**
- **Health:** http://localhost:8003/health
- **API Docs:** http://localhost:8003/docs
- **Calculate Roundup:** POST /transactions/calculate-roundup
- **Roundup History:** GET /transactions/roundup-history

## 🧪 Testing

### **Automated Tests**
```bash
# Comprehensive test suite
./test-setup.sh

# API endpoint tests
python3 test-api.py

# Import validation
python3 test-imports.py
```

### **Manual Testing**
```bash
# Test user registration
curl -X POST "http://localhost:8001/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "username": "testuser", "password": "TestPassword123"}'

# Test user login
curl -X POST "http://localhost:8001/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "TestPassword123"}'

# Test roundup calculation
curl -X POST "http://localhost:8003/transactions/calculate-roundup" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"amount": 12.34, "rounding_rule": "dollar"}'
```

## 🔍 Troubleshooting

### **Common Issues**

#### **Docker Not Running**
```bash
# Start Docker Desktop or Docker daemon
sudo systemctl start docker  # Linux
# Or start Docker Desktop application
```

#### **Port Already in Use**
```bash
# Stop existing containers
docker-compose down

# Kill processes using ports
sudo lsof -ti:8001 | xargs kill -9
sudo lsof -ti:8002 | xargs kill -9
sudo lsof -ti:8003 | xargs kill -9
sudo lsof -ti:5432 | xargs kill -9
```

#### **Build Failures**
```bash
# Clean Docker cache
docker system prune -a

# Rebuild without cache
docker-compose build --no-cache
```

#### **Database Issues**
```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

### **Debug Commands**
```bash
# Check service logs
docker-compose logs user-service
docker-compose logs plaid-service
docker-compose logs transaction-service

# Check service status
docker-compose ps

# Enter container for debugging
docker-compose exec user-service bash
```

## 📊 System Architecture

### **Services**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Service  │    │  Plaid Service  │    │Transaction Service│
│   Port: 8001    │    │   Port: 8002    │    │   Port: 8003    │
│                 │    │                 │    │                 │
│ • Authentication│    │ • Bank Accounts │    │ • Roundup Calc  │
│ • User Mgmt     │    │ • Transactions  │    │ • Analytics     │
│ • Preferences   │    │ • Webhooks      │    │ • History       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    │   Port: 5432    │
                    │                 │
                    │ • Users         │
                    │ • Accounts      │
                    │ • Transactions  │
                    │ • Calculations  │
                    └─────────────────┘
```

### **Shared Modules**
```
shared/
├── __init__.py      # Package initialization
├── config.py        # Configuration management
├── database.py      # Database connection
├── auth.py          # Authentication utilities
├── models.py        # SQLAlchemy models
├── schemas.py       # Pydantic schemas
└── utils.py         # Common utilities
```

## 🔒 Security Features

### **Authentication**
- JWT token-based authentication
- Password strength validation
- Secure password hashing with bcrypt
- Token expiration handling

### **Input Validation**
- Email format validation
- Username format validation
- Password strength requirements
- SQL injection prevention

### **Container Security**
- Non-root users in containers
- Read-only volumes where possible
- Minimal attack surface
- Resource limits

## 🚀 Production Deployment

### **Production Setup**
```bash
# Build production images
make build-prod

# Start production services
make prod

# Check production status
make status-prod
```

### **Environment Variables**
Create `.env.prod` with production values:
```bash
POSTGRES_PASSWORD=secure_password
SECRET_KEY=your_production_secret_key
PLAID_CLIENT_ID=your_production_plaid_id
PLAID_SECRET=your_production_plaid_secret
PLAID_ENVIRONMENT=production
LOG_LEVEL=WARNING
ENVIRONMENT=production
```

## 📈 Monitoring & Health Checks

### **Health Endpoints**
- All services have `/health` endpoints
- Docker health checks configured
- Automatic restart on failure
- Resource monitoring

### **Logging**
- Structured logging with timestamps
- Error tracking with context
- Service-specific log levels
- Log aggregation ready

## 🎯 Next Steps

### **Development**
1. **Start services:** `make dev`
2. **Test APIs:** Visit http://localhost:8001/docs
3. **Make changes:** Services auto-reload
4. **Run tests:** `./test-setup.sh`

### **Production**
1. **Set up environment:** Configure `.env.prod`
2. **Deploy:** `make prod`
3. **Monitor:** Check health endpoints
4. **Scale:** Add more service instances

### **Integration**
1. **Frontend:** Update Next.js to use Python APIs
2. **Plaid:** Configure production Plaid credentials
3. **Monitoring:** Set up logging and metrics
4. **CI/CD:** Add automated testing and deployment

## 🎉 Success!

Your Python FastAPI microservices are now ready for development and production use! The system includes:

- ✅ **Complete microservices architecture**
- ✅ **Secure authentication system**
- ✅ **Comprehensive API endpoints**
- ✅ **Production-ready Docker setup**
- ✅ **Automated testing suite**
- ✅ **Health monitoring**
- ✅ **Troubleshooting guides**

**Ready to start coding!** 🚀
