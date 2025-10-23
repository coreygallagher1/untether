#!/bin/bash
# Quick validation script to check for common issues

echo "🔍 Quick validation check for Python FastAPI microservices"
echo "=========================================================="

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: docker-compose.yml not found. Please run this script from the backend directory."
    exit 1
fi

echo "✅ Found docker-compose.yml"

# Check if shared directory exists
if [ ! -d "shared" ]; then
    echo "❌ Error: shared directory not found."
    exit 1
fi

echo "✅ Found shared directory"

# Check if shared modules exist
required_files=(
    "shared/__init__.py"
    "shared/config.py"
    "shared/database.py"
    "shared/auth.py"
    "shared/models.py"
    "shared/schemas.py"
    "shared/utils.py"
)

for file in "${required_files[@]}"; do
    if [ ! -f "$file" ]; then
        echo "❌ Error: $file not found."
        exit 1
    fi
done

echo "✅ All shared modules found"

# Check if service directories exist
services=("user-service" "plaid-service" "transaction-service")

for service in "${services[@]}"; do
    if [ ! -d "services/$service" ]; then
        echo "❌ Error: services/$service directory not found."
        exit 1
    fi
    
    if [ ! -f "services/$service/main.py" ]; then
        echo "❌ Error: services/$service/main.py not found."
        exit 1
    fi
    
    if [ ! -f "services/$service/requirements.txt" ]; then
        echo "❌ Error: services/$service/requirements.txt not found."
        exit 1
    fi
    
    if [ ! -f "services/$service/Dockerfile" ]; then
        echo "❌ Error: services/$service/Dockerfile not found."
        exit 1
    fi
done

echo "✅ All service directories and files found"

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Error: Docker is not running. Please start Docker first."
    echo "   You can start Docker Desktop or run: sudo systemctl start docker"
    exit 1
fi

echo "✅ Docker is running"

# Check if .env file exists (optional)
if [ ! -f ".env" ]; then
    echo "⚠️  Warning: .env file not found. Creating a basic one..."
    cat > .env << EOF
# Database
POSTGRES_PASSWORD=password

# JWT Secret (CHANGE IN PRODUCTION!)
SECRET_KEY=test-secret-key-for-development-only

# Plaid Configuration (optional for testing)
PLAID_CLIENT_ID=your_plaid_client_id
PLAID_SECRET=your_plaid_secret
PLAID_ENVIRONMENT=sandbox

# Logging
LOG_LEVEL=INFO
ENVIRONMENT=development
EOF
    echo "✅ Created basic .env file"
else
    echo "✅ Found .env file"
fi

echo ""
echo "🎉 Validation complete! Everything looks good."
echo ""
echo "Next steps:"
echo "  1. Run './test-setup.sh' for comprehensive testing"
echo "  2. Or run 'make dev' to start development environment"
echo "  3. Or run 'python3 test-api.py' for API testing"
