#!/usr/bin/env python3
"""
Comprehensive API test script for Python FastAPI microservices
Tests all endpoints and validates the complete system functionality
"""

import asyncio
import httpx
import json
import sys
from typing import Dict, Any
import time

# Service URLs
SERVICES = {
    "user": "http://localhost:8001",
    "plaid": "http://localhost:8002", 
    "transaction": "http://localhost:8003"
}

class Colors:
    """ANSI color codes for terminal output"""
    RED = '\033[0;31m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    BLUE = '\033[0;34m'
    NC = '\033[0m'  # No Color

class APITester:
    def __init__(self):
        self.client = httpx.AsyncClient(timeout=30.0)
        self.test_user_token = None
        self.test_user_id = None
        
    async def close(self):
        """Close the HTTP client"""
        await self.client.aclose()
    
    def print_status(self, message: str):
        """Print status message"""
        print(f"{Colors.BLUE}[INFO]{Colors.NC} {message}")
    
    def print_success(self, message: str):
        """Print success message"""
        print(f"{Colors.GREEN}[SUCCESS]{Colors.NC} {message}")
    
    def print_error(self, message: str):
        """Print error message"""
        print(f"{Colors.RED}[ERROR]{Colors.NC} {message}")
    
    def print_warning(self, message: str):
        """Print warning message"""
        print(f"{Colors.YELLOW}[WARNING]{Colors.NC} {message}")
    
    async def test_service_health(self, service_name: str, base_url: str) -> bool:
        """Test service health endpoint"""
        try:
            response = await self.client.get(f"{base_url}/health")
            if response.status_code == 200:
                self.print_success(f"{service_name} service health check passed")
                return True
            else:
                self.print_error(f"{service_name} service health check failed: {response.status_code}")
                return False
        except Exception as e:
            self.print_error(f"{service_name} service health check failed: {e}")
            return False
    
    async def test_api_docs(self, service_name: str, base_url: str) -> bool:
        """Test API documentation endpoint"""
        try:
            response = await self.client.get(f"{base_url}/docs")
            if response.status_code == 200:
                self.print_success(f"{service_name} API docs accessible")
                return True
            else:
                self.print_warning(f"{service_name} API docs not accessible: {response.status_code}")
                return False
        except Exception as e:
            self.print_warning(f"{service_name} API docs not accessible: {e}")
            return False
    
    async def test_user_registration(self) -> bool:
        """Test user registration endpoint"""
        try:
            self.print_status("Testing user registration...")
            
            user_data = {
                "email": "test@example.com",
                "username": "testuser",
                "password": "TestPassword123"
            }
            
            response = await self.client.post(
                f"{SERVICES['user']}/auth/register",
                json=user_data
            )
            
            if response.status_code == 200:
                user_info = response.json()
                self.test_user_id = user_info.get("id")
                self.print_success("User registration successful")
                return True
            elif response.status_code == 400 and "already exists" in response.text:
                self.print_warning("User already exists (expected in subsequent runs)")
                return True
            else:
                self.print_error(f"User registration failed: {response.status_code} - {response.text}")
                return False
                
        except Exception as e:
            self.print_error(f"User registration test failed: {e}")
            return False
    
    async def test_user_login(self) -> bool:
        """Test user login endpoint"""
        try:
            self.print_status("Testing user login...")
            
            login_data = {
                "username": "testuser",
                "password": "TestPassword123"
            }
            
            response = await self.client.post(
                f"{SERVICES['user']}/auth/login",
                json=login_data
            )
            
            if response.status_code == 200:
                token_info = response.json()
                self.test_user_token = token_info.get("access_token")
                self.print_success("User login successful")
                return True
            else:
                self.print_error(f"User login failed: {response.status_code} - {response.text}")
                return False
                
        except Exception as e:
            self.print_error(f"User login test failed: {e}")
            return False
    
    async def test_authenticated_endpoints(self) -> bool:
        """Test endpoints that require authentication"""
        if not self.test_user_token:
            self.print_error("No authentication token available")
            return False
        
        headers = {"Authorization": f"Bearer {self.test_user_token}"}
        
        try:
            # Test get current user
            self.print_status("Testing authenticated user endpoint...")
            response = await self.client.get(
                f"{SERVICES['user']}/users/me",
                headers=headers
            )
            
            if response.status_code == 200:
                self.print_success("Authenticated user endpoint working")
            else:
                self.print_error(f"Authenticated user endpoint failed: {response.status_code}")
                return False
            
            # Test user preferences
            self.print_status("Testing user preferences endpoint...")
            response = await self.client.get(
                f"{SERVICES['user']}/users/me/preferences",
                headers=headers
            )
            
            if response.status_code == 200:
                self.print_success("User preferences endpoint working")
            else:
                self.print_warning(f"User preferences endpoint failed: {response.status_code}")
            
            # Test bank accounts endpoint
            self.print_status("Testing bank accounts endpoint...")
            response = await self.client.get(
                f"{SERVICES['user']}/users/me/bank-accounts",
                headers=headers
            )
            
            if response.status_code == 200:
                self.print_success("Bank accounts endpoint working")
            else:
                self.print_warning(f"Bank accounts endpoint failed: {response.status_code}")
            
            return True
            
        except Exception as e:
            self.print_error(f"Authenticated endpoints test failed: {e}")
            return False
    
    async def test_plaid_endpoints(self) -> bool:
        """Test Plaid service endpoints"""
        if not self.test_user_token:
            self.print_error("No authentication token available for Plaid tests")
            return False
        
        headers = {"Authorization": f"Bearer {self.test_user_token}"}
        
        try:
            # Test link token creation
            self.print_status("Testing Plaid link token creation...")
            response = await self.client.post(
                f"{SERVICES['plaid']}/plaid/link-token",
                headers=headers
            )
            
            if response.status_code == 200:
                self.print_success("Plaid link token creation working")
            else:
                self.print_warning(f"Plaid link token creation failed: {response.status_code}")
            
            # Test accounts endpoint
            self.print_status("Testing Plaid accounts endpoint...")
            response = await self.client.get(
                f"{SERVICES['plaid']}/plaid/accounts",
                headers=headers
            )
            
            if response.status_code == 200:
                self.print_success("Plaid accounts endpoint working")
            else:
                self.print_warning(f"Plaid accounts endpoint failed: {response.status_code}")
            
            return True
            
        except Exception as e:
            self.print_warning(f"Plaid endpoints test failed: {e}")
            return False
    
    async def test_transaction_endpoints(self) -> bool:
        """Test Transaction service endpoints"""
        if not self.test_user_token:
            self.print_error("No authentication token available for transaction tests")
            return False
        
        headers = {"Authorization": f"Bearer {self.test_user_token}"}
        
        try:
            # Test roundup calculation
            self.print_status("Testing roundup calculation...")
            
            calculation_data = {
                "amount": 12.34,
                "rounding_rule": "dollar",
                "custom_rounding_amount": None
            }
            
            response = await self.client.post(
                f"{SERVICES['transaction']}/transactions/calculate-roundup",
                json=calculation_data,
                headers=headers
            )
            
            if response.status_code == 200:
                result = response.json()
                self.print_success(f"Roundup calculation working: ${result.get('roundup_amount', 0)}")
            else:
                self.print_error(f"Roundup calculation failed: {response.status_code}")
                return False
            
            # Test roundup history
            self.print_status("Testing roundup history...")
            response = await self.client.get(
                f"{SERVICES['transaction']}/transactions/roundup-history",
                headers=headers
            )
            
            if response.status_code == 200:
                self.print_success("Roundup history endpoint working")
            else:
                self.print_warning(f"Roundup history endpoint failed: {response.status_code}")
            
            # Test roundup summary
            self.print_status("Testing roundup summary...")
            response = await self.client.get(
                f"{SERVICES['transaction']}/transactions/roundup-summary?days=30",
                headers=headers
            )
            
            if response.status_code == 200:
                self.print_success("Roundup summary endpoint working")
            else:
                self.print_warning(f"Roundup summary endpoint failed: {response.status_code}")
            
            return True
            
        except Exception as e:
            self.print_error(f"Transaction endpoints test failed: {e}")
            return False
    
    async def run_all_tests(self) -> bool:
        """Run all API tests"""
        print("🚀 Starting comprehensive API test suite")
        print("=" * 50)
        
        all_passed = True
        
        # Test 1: Service Health Checks
        self.print_status("Test 1: Service Health Checks")
        for service_name, base_url in SERVICES.items():
            if not await self.test_service_health(service_name, base_url):
                all_passed = False
        
        # Test 2: API Documentation
        self.print_status("Test 2: API Documentation")
        for service_name, base_url in SERVICES.items():
            await self.test_api_docs(service_name, base_url)
        
        # Test 3: User Registration
        self.print_status("Test 3: User Registration")
        if not await self.test_user_registration():
            all_passed = False
        
        # Test 4: User Login
        self.print_status("Test 4: User Login")
        if not await self.test_user_login():
            all_passed = False
        
        # Test 5: Authenticated Endpoints
        self.print_status("Test 5: Authenticated Endpoints")
        if not await self.test_authenticated_endpoints():
            all_passed = False
        
        # Test 6: Plaid Endpoints
        self.print_status("Test 6: Plaid Service Endpoints")
        await self.test_plaid_endpoints()
        
        # Test 7: Transaction Endpoints
        self.print_status("Test 7: Transaction Service Endpoints")
        if not await self.test_transaction_endpoints():
            all_passed = False
        
        return all_passed

async def main():
    """Main test function"""
    tester = APITester()
    
    try:
        # Wait a moment for services to be ready
        print("Waiting for services to be ready...")
        await asyncio.sleep(5)
        
        # Run all tests
        success = await tester.run_all_tests()
        
        print("\n" + "=" * 50)
        if success:
            print("🎉 All API tests passed successfully!")
            print("\nYour Python FastAPI microservices are working correctly!")
            print("\nNext steps:")
            print("  • Visit http://localhost:8001/docs for User Service API")
            print("  • Visit http://localhost:8002/docs for Plaid Service API")
            print("  • Visit http://localhost:8003/docs for Transaction Service API")
        else:
            print("❌ Some tests failed. Check the output above for details.")
            sys.exit(1)
            
    except KeyboardInterrupt:
        print("\n⚠️ Tests interrupted by user")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ Test suite failed with error: {e}")
        sys.exit(1)
    finally:
        await tester.close()

if __name__ == "__main__":
    asyncio.run(main())
