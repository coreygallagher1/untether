# Plaid service - FastAPI application
from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy.orm import Session
from typing import List, Optional
import sys
import os
import httpx
import json

# Add shared modules to path
sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', 'shared'))

from shared.database import get_db
from shared.models import User, BankAccount, PlaidItem
from shared.schemas import BankAccountResponse, LinkTokenResponse, ExchangeTokenRequest, ExchangeTokenResponse, PlaidAccountsResponse, PlaidBalanceResponse
from shared.auth import get_current_user_token
from shared.config import settings
from shared.utils import create_error_response, log_error

# Create FastAPI app
app = FastAPI(
    title="Plaid Service",
    description="Bank account integration and Plaid API service",
    version="1.0.0"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.allowed_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Plaid API client
class PlaidClient:
    def __init__(self):
        self.client_id = settings.plaid_client_id
        self.secret = settings.plaid_secret
        self.environment = settings.plaid_environment
        self.base_url = f"https://{self.environment}.plaid.com"
    
    async def create_link_token(self, user_id: str) -> str:
        """Create a Plaid link token for a user."""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/link/token/create",
                json={
                    "client_id": self.client_id,
                    "secret": self.secret,
                    "client_name": "Untether",
                    "products": ["transactions", "auth"],
                    "country_codes": ["US"],
                    "language": "en",
                    "user": {
                        "client_user_id": user_id
                    }
                }
            )
            response.raise_for_status()
            data = response.json()
            return data["link_token"]
    
    async def exchange_public_token(self, public_token: str) -> str:
        """Exchange a public token for an access token."""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/item/public_token/exchange",
                json={
                    "client_id": self.client_id,
                    "secret": self.secret,
                    "public_token": public_token
                }
            )
            response.raise_for_status()
            data = response.json()
            return data["access_token"]
    
    async def get_accounts(self, access_token: str) -> List[dict]:
        """Get accounts for an access token."""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/accounts/get",
                json={
                    "client_id": self.client_id,
                    "secret": self.secret,
                    "access_token": access_token
                }
            )
            response.raise_for_status()
            data = response.json()
            return data["accounts"]
    
    async def get_balance(self, access_token: str, account_id: str) -> float:
        """Get balance for a specific account."""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/accounts/balance/get",
                json={
                    "client_id": self.client_id,
                    "secret": self.secret,
                    "access_token": access_token,
                    "account_ids": [account_id]
                }
            )
            response.raise_for_status()
            data = response.json()
            accounts = data["accounts"]
            if accounts:
                return accounts[0]["balances"]["available"] or 0.0
            return 0.0
    
    async def get_transactions(self, access_token: str, start_date: str, end_date: str) -> List[dict]:
        """Get transactions for a date range."""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/transactions/get",
                json={
                    "client_id": self.client_id,
                    "secret": self.secret,
                    "access_token": access_token,
                    "start_date": start_date,
                    "end_date": end_date
                }
            )
            response.raise_for_status()
            data = response.json()
            return data["transactions"]

# Initialize Plaid client
plaid_client = PlaidClient()

# Dependency to get current user from token
def get_current_user(token_data: dict = Depends(get_current_user_token), db: Session = Depends(get_db)):
    """Get current user from JWT token."""
    try:
        username = token_data.get("sub")
        if username is None:
            raise create_error_response("Invalid token", status.HTTP_401_UNAUTHORIZED)
        
        user = db.query(User).filter(User.username == username).first()
        if user is None:
            raise create_error_response("User not found", status.HTTP_404_NOT_FOUND)
        
        return user
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "get_current_user")
        raise create_error_response("Authentication failed", status.HTTP_401_UNAUTHORIZED)

# Plaid endpoints
@app.post("/plaid/link-token")
async def create_link_token(
    current_user: User = Depends(get_current_user)
):
    """Create a Plaid link token for the current user."""
    try:
        link_token = await plaid_client.create_link_token(str(current_user.id))
        return {"link_token": link_token}
    except Exception as e:
        log_error(e, "create_link_token")
        raise create_error_response("Failed to create link token", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.post("/plaid/exchange-token")
async def exchange_public_token(
    public_token: str,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Exchange a public token for an access token and store it."""
    try:
        access_token = await plaid_client.exchange_public_token(public_token)
        
        # Store the access token in the database
        # Note: In a real implementation, you'd want to encrypt this
        from shared.models import PlaidItem
        plaid_item = PlaidItem(
            user_id=current_user.id,
            access_token=access_token,
            item_id="",  # You'd get this from the exchange response
            institution_id="",  # You'd get this from the exchange response
            status="active"
        )
        db.add(plaid_item)
        db.commit()
        
        return {"access_token": access_token, "success": True}
    except Exception as e:
        log_error(e, "exchange_public_token")
        raise create_error_response("Failed to exchange token", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.get("/plaid/accounts")
async def get_accounts(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get accounts for the current user."""
    try:
        # Get user's access token
        from shared.models import PlaidItem
        plaid_item = db.query(PlaidItem).filter(
            PlaidItem.user_id == current_user.id,
            PlaidItem.status == "active"
        ).first()
        
        if not plaid_item:
            raise create_error_response("No active Plaid connection found", status.HTTP_404_NOT_FOUND)
        
        accounts = await plaid_client.get_accounts(plaid_item.access_token)
        
        # Convert to our format
        formatted_accounts = []
        for account in accounts:
            balance = await plaid_client.get_balance(plaid_item.access_token, account["account_id"])
            formatted_accounts.append({
                "account_id": account["account_id"],
                "name": account["name"],
                "type": account["type"],
                "subtype": account["subtype"],
                "mask": account.get("mask", ""),
                "balance": balance
            })
        
        return {"accounts": formatted_accounts}
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "get_accounts")
        raise create_error_response("Failed to get accounts", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.get("/plaid/accounts/{account_id}/balance")
async def get_account_balance(
    account_id: str,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get balance for a specific account."""
    try:
        # Get user's access token
        from shared.models import PlaidItem
        plaid_item = db.query(PlaidItem).filter(
            PlaidItem.user_id == current_user.id,
            PlaidItem.status == "active"
        ).first()
        
        if not plaid_item:
            raise create_error_response("No active Plaid connection found", status.HTTP_404_NOT_FOUND)
        
        balance = await plaid_client.get_balance(plaid_item.access_token, account_id)
        return {"account_id": account_id, "balance": balance}
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "get_account_balance")
        raise create_error_response("Failed to get balance", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.get("/plaid/transactions")
async def get_transactions(
    start_date: str,
    end_date: str,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get transactions for a date range."""
    try:
        # Get user's access token
        from shared.models import PlaidItem
        plaid_item = db.query(PlaidItem).filter(
            PlaidItem.user_id == current_user.id,
            PlaidItem.status == "active"
        ).first()
        
        if not plaid_item:
            raise create_error_response("No active Plaid connection found", status.HTTP_404_NOT_FOUND)
        
        transactions = await plaid_client.get_transactions(
            plaid_item.access_token, start_date, end_date
        )
        
        return {"transactions": transactions}
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "get_transactions")
        raise create_error_response("Failed to get transactions", status.HTTP_500_INTERNAL_SERVER_ERROR)

# Webhook endpoint for Plaid updates
@app.post("/plaid/webhook")
async def handle_webhook(
    webhook_data: dict,
    db: Session = Depends(get_db)
):
    """Handle Plaid webhook events."""
    try:
        webhook_type = webhook_data.get("webhook_type")
        
        if webhook_type == "TRANSACTIONS":
            # Handle transaction updates
            pass
        elif webhook_type == "ITEM":
            # Handle item updates (like login errors)
            pass
        
        return {"status": "success"}
    except Exception as e:
        log_error(e, "handle_webhook")
        raise create_error_response("Webhook processing failed", status.HTTP_500_INTERNAL_SERVER_ERROR)

# Health check
@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {"status": "healthy", "service": "plaid-service"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8002)
