# Troubleshooting Guide for Python FastAPI Microservices

## Quick Start Checklist

### Prerequisites
- [ ] Docker installed and running
- [ ] Docker Compose installed
- [ ] Python 3.11+ (for local development)
- [ ] Git (for version control)

### Step 1: Validate Setup
```bash
cd backend
./validate-setup.sh
```

### Step 2: Start Services
```bash
# Development environment
make dev

# Or manually
docker-compose up --build
```

### Step 3: Run Tests
```bash
# Comprehensive test suite
./test-setup.sh

# API tests only
python3 test-api.py
```

## Common Issues & Solutions

### 1. Docker Not Running
**Error:** `Docker is not running. Please start Docker first.`

**Solution:**
- Start Docker Desktop
- Or start Docker daemon: `sudo systemctl start docker`

### 2. Port Already in Use
**Error:** `Port 8001/8002/8003/5432 is already in use`

**Solutions:**
```bash
# Stop existing containers
docker-compose down

# Kill processes using ports
sudo lsof -ti:8001 | xargs kill -9
sudo lsof -ti:8002 | xargs kill -9
sudo lsof -ti:8003 | xargs kill -9
sudo lsof -ti:5432 | xargs kill -9
```

### 3. Build Failures
**Error:** `Failed to build Docker images`

**Solutions:**
```bash
# Clean Docker cache
docker system prune -a

# Rebuild without cache
docker-compose build --no-cache

# Check Dockerfile syntax
docker build -t test-build services/user-service/
```

### 4. Database Connection Issues
**Error:** `Database connection failed`

**Solutions:**
```bash
# Check PostgreSQL container
docker-compose logs postgres

# Restart database
docker-compose restart postgres

# Check database health
docker-compose exec postgres pg_isready -U postgres
```

### 5. Import Errors
**Error:** `ModuleNotFoundError: No module named 'shared'`

**Solutions:**
- Ensure you're running from the correct directory
- Check that `shared/__init__.py` exists
- Verify Python path in Dockerfile

### 6. Authentication Issues
**Error:** `Authentication failed`

**Solutions:**
- Check JWT secret key in `.env` file
- Verify token format in requests
- Check user exists in database

### 7. Service Health Check Failures
**Error:** `Health check failed`

**Solutions:**
```bash
# Check service logs
docker-compose logs user-service
docker-compose logs plaid-service
docker-compose logs transaction-service

# Check service status
docker-compose ps

# Restart specific service
docker-compose restart user-service
```

## Development Workflow

### Starting Development
```bash
# 1. Start services
make dev

# 2. Check status
make status

# 3. View logs
make logs

# 4. Run tests
./test-setup.sh
```

### Making Changes
```bash
# 1. Edit code
# 2. Services auto-reload (development mode)
# 3. Test changes
python3 test-api.py
```

### Stopping Services
```bash
# Stop and remove containers
make down

# Stop and remove volumes (WARNING: deletes data)
make clean
```

## API Testing

### Manual Testing
```bash
# Test health endpoints
curl http://localhost:8001/health
curl http://localhost:8002/health
curl http://localhost:8003/health

# Test API docs
open http://localhost:8001/docs
open http://localhost:8002/docs
open http://localhost:8003/docs
```

### Automated Testing
```bash
# Run comprehensive tests
./test-setup.sh

# Run API tests
python3 test-api.py
```

## Database Management

### Running Migrations
```bash
# Run migrations
make migrate

# Create new migration
make create-migration name=add_new_field
```

### Database Access
```bash
# Connect to database
docker-compose exec postgres psql -U postgres -d untether

# Backup database
docker-compose exec postgres pg_dump -U postgres untether > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres untether < backup.sql
```

## Production Deployment

### Production Setup
```bash
# Build production images
make build-prod

# Start production services
make prod

# Check production status
make status-prod
```

### Environment Variables
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

## Monitoring & Debugging

### View Logs
```bash
# All services
make logs

# Specific service
docker-compose logs user-service
docker-compose logs plaid-service
docker-compose logs transaction-service

# Follow logs
docker-compose logs -f user-service
```

### Check Resource Usage
```bash
# Container resource usage
docker stats

# Service status
docker-compose ps

# Health checks
curl http://localhost:8001/health
curl http://localhost:8002/health
curl http://localhost:8003/health
```

### Debug Container Issues
```bash
# Enter container
docker-compose exec user-service bash

# Check Python environment
docker-compose exec user-service python -c "import sys; print(sys.path)"

# Check installed packages
docker-compose exec user-service pip list
```

## Performance Optimization

### Build Optimization
```bash
# Use build cache
docker-compose build

# Parallel builds
docker-compose build --parallel
```

### Runtime Optimization
```bash
# Set resource limits
docker-compose -f docker-compose.prod.yml up -d

# Monitor performance
docker stats
```

## Security Considerations

### Development Security
- Use strong passwords in `.env`
- Don't commit `.env` files
- Use non-root users in containers
- Enable HTTPS in production

### Production Security
- Use secrets management
- Enable firewall rules
- Regular security updates
- Monitor for vulnerabilities

## Getting Help

### Check Logs First
```bash
docker-compose logs
```

### Common Commands
```bash
# Restart everything
docker-compose down && docker-compose up --build

# Clean everything
make clean

# Check Docker status
docker system df
docker system prune
```

### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=DEBUG
docker-compose up
```

If you're still having issues, check the logs and run the validation script to identify the specific problem.
