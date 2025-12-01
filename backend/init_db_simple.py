#!/usr/bin/env python3
"""
Simple database initialization script.
Use this instead of migrations during rapid development.

This creates all tables from your SQLAlchemy models automatically.
Just drop and recreate your database when you change models.
"""
import sys
import os

# Add shared to path
sys.path.insert(0, os.path.join(os.path.dirname(__file__)))

from shared.database import init_db, engine
from shared.config import settings

def main():
    print(f"Initializing database: {settings.database_url}")
    print("Creating all tables from models...")
    
    try:
        init_db()
        print("✅ Database initialized successfully!")
        print("\nTables created:")
        print("  - users")
        print("  - user_preferences")
        print("  - plaid_items")
        print("  - bank_accounts")
        print("  - roundup_calculations")
    except Exception as e:
        print(f"❌ Error initializing database: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

