# Testing Simple Database Initialization

## Quick Test

1. **Start Docker Desktop** (if not already running)

2. **Start the database:**
   ```bash
   docker compose up -d postgres
   ```

3. **Wait a few seconds for database to be ready:**
   ```bash
   sleep 5
   ```

4. **Run the initialization:**
   ```bash
   make init-db
   ```
   
   This runs the initialization inside the user-service container (which has all dependencies).
   
   Alternatively, you can run it manually:
   ```bash
   docker compose run --rm user-service python3 -c "import sys; sys.path.insert(0, '/app'); from shared.database import init_db; init_db(); print('✅ Done!')"
   ```

5. **Verify tables were created:**
   ```bash
   docker compose exec postgres psql -U postgres -d untether -c '\dt'
   ```

   You should see:
   - users
   - user_preferences
   - plaid_items
   - bank_accounts
   - roundup_calculations

## Full Test Script

Run the automated test:
```bash
./test_init.sh
```

## What to Expect

✅ **Success output:**
```
Initializing database: postgresql://postgres:password@localhost:5432/untether
Creating all tables from models...
✅ Database initialized successfully!

Tables created:
  - users
  - user_preferences
  - plaid_items
  - bank_accounts
  - roundup_calculations
```

## Troubleshooting

**If you get "connection refused":**
- Make sure Docker is running
- Wait a bit longer for postgres to start
- Check: `docker compose ps`

**If you get import errors:**
- Make sure you're in the `/backend` directory
- Check that `shared/` directory exists

**If tables already exist:**
- That's fine! The script is idempotent
- To start fresh: `docker compose down -v` then `docker compose up -d postgres`

## Next Steps

After initialization:
1. Start all services: `make dev`
2. Test API: `python3 test-api.py`
3. Check service health: `curl http://localhost:8001/health`

