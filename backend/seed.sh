#!/bin/bash
# Simple database seeding script

echo "🌱 Seeding database with test user..."

# Create test user
docker compose exec postgres psql -U postgres -d untether -c "
INSERT INTO users (email, username, first_name, last_name, hashed_password, is_active, is_verified, created_at) 
VALUES ('test@example.com', 'testuser', 'Test', 'User', '\$2b\$12\$yfscfjlRHAJbyBjvKd2PAOFB/Ahs.whrWNYWSrS/lT3VSJe9KLr/u', true, true, NOW()) 
ON CONFLICT (email) DO NOTHING;
"

# Create user preferences
docker compose exec postgres psql -U postgres -d untether -c "
INSERT INTO user_preferences (user_id, roundup_enabled, custom_rounding_amount, donation_percentage, created_at) 
SELECT id, true, 0.00, 50, NOW() FROM users WHERE email = 'test@example.com' 
AND NOT EXISTS (SELECT 1 FROM user_preferences WHERE user_id = users.id);
"

echo "✅ Test user created successfully!"
echo ""
echo "📋 Test Account:"
echo "   Email: test@example.com"
echo "   Username: testuser" 
echo "   Password: TestPass123!"
echo ""
echo "🔗 You can now test the application!"
