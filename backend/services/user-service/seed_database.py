#!/usr/bin/env python3
"""
Simple database seeding script for Untether
Creates a basic test user for development
"""

import sys
import os
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from shared.database import SessionLocal
from shared.models import User, UserPreferences
from shared.auth import get_password_hash
from sqlalchemy import text
from datetime import datetime, timedelta

def clear_existing_data():
    """Clear existing test data"""
    print("🧹 Clearing existing test data...")
    
    db = SessionLocal()
    try:
        # Delete test user and preferences
        db.execute(text("DELETE FROM user_preferences WHERE user_id IN (SELECT id FROM users WHERE email = 'test@example.com')"))
        db.execute(text("DELETE FROM users WHERE email = 'test@example.com'"))
        
        db.commit()
        print("✅ Existing test data cleared")
    except Exception as e:
        print(f"❌ Error clearing data: {e}")
        db.rollback()
    finally:
        db.close()

def create_test_user():
    """Create one test user with basic data"""
    print("👤 Creating test user...")
    
    db = SessionLocal()
    try:
        # Test User: Basic profile
        test_user = User(
            email="test@example.com",
            username="testuser",
            first_name="Test",
            last_name="User",
            hashed_password=get_password_hash("TestPass123!"),
            is_active=True,
            is_verified=True,
            created_at=datetime.utcnow() - timedelta(days=30)
        )
        db.add(test_user)
        db.flush()  # Get the ID
        
        # User preferences
        user_prefs = UserPreferences(
            user_id=test_user.id,
            roundup_enabled=True,
            roundup_multiplier=2,
            notifications_enabled=True,
            created_at=datetime.utcnow() - timedelta(days=30)
        )
        db.add(user_prefs)
        
        db.commit()
        print("✅ Test user created successfully")
        print(f"   - Email: test@example.com")
        print(f"   - Username: testuser")
        print(f"   - Password: TestPass123!")
        
    except Exception as e:
        print(f"❌ Error creating test user: {e}")
        db.rollback()
        raise
    finally:
        db.close()

def verify_data():
    """Verify the seeded data"""
    print("🔍 Verifying seeded data...")
    
    db = SessionLocal()
    try:
        # Count users
        user_count = db.query(User).filter(User.email == 'test@example.com').count()
        print(f"   - Users created: {user_count}")
        
        # Count preferences
        pref_count = db.query(UserPreferences).join(User).filter(User.email == 'test@example.com').count()
        print(f"   - User preferences: {pref_count}")
        
        print("✅ Data verification complete")
        
    except Exception as e:
        print(f"❌ Error verifying data: {e}")
    finally:
        db.close()

def main():
    """Main seeding function"""
    print("🌱 Starting database seeding...")
    print("=" * 50)
    
    try:
        # Clear existing test data
        clear_existing_data()
        
        # Create test user
        create_test_user()
        
        # Verify everything was created
        verify_data()
        
        print("=" * 50)
        print("🎉 Database seeding completed successfully!")
        print("\n📋 Test Account Created:")
        print("   Email: test@example.com | Username: testuser | Password: TestPass123!")
        print("\n🔗 You can now test the application with this account!")
        
    except Exception as e:
        print(f"❌ Seeding failed: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()