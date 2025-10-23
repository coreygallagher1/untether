#!/usr/bin/env python3
"""
Database seeding script for Untether
Creates default test data for development
"""

import sys
import os
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from shared.database import SessionLocal, engine
from shared.models import User, UserPreferences, PlaidItem, BankAccount, RoundupCalculation
from shared.auth import get_password_hash
from sqlalchemy import text
from datetime import datetime, timedelta
import random

def clear_existing_data():
    """Clear existing test data"""
    print("🧹 Clearing existing test data...")
    
    db = SessionLocal()
    try:
        # Delete in reverse order of dependencies
        db.execute(text("DELETE FROM roundup_calculations WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'test%@example.com')"))
        db.execute(text("DELETE FROM bank_accounts WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'test%@example.com')"))
        db.execute(text("DELETE FROM plaid_items WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'test%@example.com')"))
        db.execute(text("DELETE FROM user_preferences WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'test%@example.com')"))
        db.execute(text("DELETE FROM users WHERE email LIKE 'test%@example.com'"))
        
        db.commit()
        print("✅ Existing test data cleared")
    except Exception as e:
        print(f"❌ Error clearing data: {e}")
        db.rollback()
    finally:
        db.close()

def create_test_user():
    """Create one test user with related data"""
    print("👤 Creating test user...")
    
    db = SessionLocal()
    try:
        # Test User: Complete profile
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
        
        # Plaid item
        plaid_item = PlaidItem(
            user_id=test_user.id,
            access_token="test_access_token_123",
            item_id="test_item_123",
            institution_id="ins_123",
            institution_name="Test Bank",
            status="active",
            created_at=datetime.utcnow() - timedelta(days=25)
        )
        db.add(plaid_item)
        db.flush()
        
        # Bank accounts
        bank_account1 = BankAccount(
            user_id=test_user.id,
            plaid_item_id=plaid_item.id,
            plaid_account_id="acc_123",
            account_name="Test Checking Account",
            account_type="depository",
            account_subtype="checking",
            mask="1234",
            is_active=True,
            created_at=datetime.utcnow() - timedelta(days=25)
        )
        db.add(bank_account1)
        
        bank_account2 = BankAccount(
            user_id=test_user.id,
            plaid_item_id=plaid_item.id,
            plaid_account_id="acc_456",
            account_name="Test Savings Account",
            account_type="depository",
            account_subtype="savings",
            mask="5678",
            is_active=True,
            created_at=datetime.utcnow() - timedelta(days=25)
        )
        db.add(bank_account2)
        
        # Sample roundup calculations
        for i in range(3):
            roundup = RoundupCalculation(
                user_id=test_user.id,
                transaction_id=f"txn_{i+1}",
                amount=round(random.uniform(10.00, 100.00), 2),
                rounding_rule="next_dollar",
                custom_rounding_amount=None,
                rounded_amount=round(random.uniform(10.00, 100.00), 2),
                roundup_amount=round(random.uniform(0.01, 0.99), 2),
                created_at=datetime.utcnow() - timedelta(days=random.randint(1, 20))
            )
            db.add(roundup)
        
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
        
        # Count plaid items
        plaid_count = db.query(PlaidItem).join(User).filter(User.email == 'test@example.com').count()
        print(f"   - Plaid items: {plaid_count}")
        
        # Count bank accounts
        bank_count = db.query(BankAccount).join(User).filter(User.email == 'test@example.com').count()
        print(f"   - Bank accounts: {bank_count}")
        
        # Count roundup calculations
        roundup_count = db.query(RoundupCalculation).join(User).filter(User.email == 'test@example.com').count()
        print(f"   - Roundup calculations: {roundup_count}")
        
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
        
        # Create test user and related data
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
