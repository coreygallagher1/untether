#!/bin/bash
# Test script for simple database initialization

set -e

echo "🧪 Testing Simple Database Initialization"
echo "=========================================="
echo ""

# Step 1: Check if Docker is running
echo "Step 1: Checking Docker..."
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop first."
    exit 1
fi
echo "✅ Docker is running"
echo ""

# Step 2: Start services
echo "Step 2: Starting services..."
docker compose up -d postgres
echo "⏳ Waiting for database to be ready..."
sleep 5
echo "✅ Services started"
echo ""

# Step 3: Run initialization
echo "Step 3: Initializing database..."
python3 init_db_simple.py
echo ""

# Step 4: Verify tables were created
echo "Step 4: Verifying tables were created..."
TABLES=$(docker compose exec -T postgres psql -U postgres -d untether -t -c "SELECT tablename FROM pg_tables WHERE schemaname='public';" | tr -d ' ' | grep -v '^$')

if echo "$TABLES" | grep -q "users"; then
    echo "✅ users table exists"
else
    echo "❌ users table missing"
    exit 1
fi

if echo "$TABLES" | grep -q "user_preferences"; then
    echo "✅ user_preferences table exists"
else
    echo "❌ user_preferences table missing"
    exit 1
fi

if echo "$TABLES" | grep -q "plaid_items"; then
    echo "✅ plaid_items table exists"
else
    echo "❌ plaid_items table missing"
    exit 1
fi

if echo "$TABLES" | grep -q "bank_accounts"; then
    echo "✅ bank_accounts table exists"
else
    echo "❌ bank_accounts table missing"
    exit 1
fi

if echo "$TABLES" | grep -q "roundup_calculations"; then
    echo "✅ roundup_calculations table exists"
else
    echo "❌ roundup_calculations table missing"
    exit 1
fi

echo ""
echo "🎉 All tests passed! Database initialized successfully."
echo ""
echo "You can now:"
echo "  - Start your services: make dev"
echo "  - Test the API: python3 test-api.py"
echo "  - View tables: docker compose exec postgres psql -U postgres -d untether -c '\\dt'"

