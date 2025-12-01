# Adding Database Migrations Later

If you need to add Alembic migrations back (for production, team collaboration, or data preservation), follow these steps:

## Quick Setup

1. **Install Alembic** (if not already in requirements.txt):
   ```bash
   pip install alembic
   ```

2. **Initialize Alembic in a service**:
   ```bash
   cd services/user-service
   alembic init migrations
   ```

3. **Configure `migrations/env.py`**:
   ```python
   import sys
   import os
   sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', 'shared'))
   
   from shared.database import Base
   from shared.models import User, UserPreferences, PlaidItem, BankAccount, RoundupCalculation
   
   target_metadata = Base.metadata
   ```

4. **Configure `alembic.ini`**:
   ```ini
   sqlalchemy.url = postgresql://postgres:password@localhost:5432/untether
   ```

5. **Create initial migration**:
   ```bash
   alembic revision --autogenerate -m "initial"
   ```

6. **Run migrations**:
   ```bash
   alembic upgrade head
   ```

## When to Use Migrations

- ✅ Production deployments with existing data
- ✅ Team collaboration (coordinated schema changes)
- ✅ Need to rollback changes
- ✅ Multiple environments (dev/staging/prod)

## Current Approach (Simple Init)

For now, we use `Base.metadata.create_all()` which:
- Creates tables directly from models
- Perfect for rapid development
- No migration files to manage
- Just drop/recreate database when models change

Use: `make init-db`

