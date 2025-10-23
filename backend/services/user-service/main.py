# User service - FastAPI application
from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy.orm import Session
from sqlalchemy import or_
from typing import List, Union
import sys
import os

# Add shared modules to path
sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', 'shared'))

from shared.database import get_db
from shared.models import User, UserPreferences, BankAccount
from shared.schemas import (
    UserCreate, UserUpdate, UserResponse, UserUpdateResponse,
    LoginRequest, PasswordChangeRequest, Token, UserPreferencesCreate, 
    UserPreferencesUpdate, UserPreferencesResponse,
    BankAccountCreate, BankAccountResponse
)
from shared.auth import verify_password, get_password_hash, create_access_token, get_current_user_token, validate_password_strength
from shared.config import settings
from shared.utils import create_error_response, log_error, validate_username, validate_email_format

# Create FastAPI app
app = FastAPI(
    title="User Service",
    description="User management and authentication service",
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

# Auth endpoints
@app.post("/auth/register", response_model=Token)
async def register_user(user: UserCreate, db: Session = Depends(get_db)):
    """Register a new user."""
    try:
        # Validate email format
        if not validate_email_format(user.email):
            raise create_error_response(
                f"'{user.email}' is not a valid email address. Please enter a valid email.",
                status.HTTP_400_BAD_REQUEST,
                error_code="INVALID_EMAIL"
            )
        
        # Validate username format
        if not validate_username(user.username):
            raise create_error_response(
                f"Username '{user.username}' is invalid. Username must be 3-50 characters and contain only letters, numbers, and underscores.",
                status.HTTP_400_BAD_REQUEST,
                error_code="INVALID_USERNAME"
            )
        
        # Validate password strength
        if not validate_password_strength(user.password):
            raise create_error_response(
                "Password does not meet security requirements. Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.",
                status.HTTP_400_BAD_REQUEST,
                error_code="WEAK_PASSWORD"
            )
        
        # Check if email already exists (case-insensitive)
        existing_email = db.query(User).filter(User.email.ilike(user.email)).first()
        if existing_email:
            raise create_error_response(
                f"An account with the email '{user.email}' already exists. Please use a different email or try logging in.",
                status.HTTP_400_BAD_REQUEST,
                error_code="EMAIL_EXISTS"
            )
        
        # Check if username already exists
        existing_username = db.query(User).filter(User.username == user.username).first()
        if existing_username:
            raise create_error_response(
                f"The username '{user.username}' is already taken. Please choose a different username.",
                status.HTTP_400_BAD_REQUEST,
                error_code="USERNAME_EXISTS"
            )
        
        # Create new user
        hashed_password = get_password_hash(user.password)
        db_user = User(
            email=user.email,
            username=user.username,
            first_name=user.first_name,
            last_name=user.last_name,
            hashed_password=hashed_password
        )
        
        db.add(db_user)
        db.commit()
        db.refresh(db_user)
        
        # Create JWT token for immediate login
        access_token = create_access_token(data={"sub": db_user.username})
        return {"access_token": access_token, "token_type": "bearer"}
        
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "register_user")
        raise create_error_response("Registration failed", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.post("/auth/login", response_model=Token)
async def login_user(login_data: LoginRequest, db: Session = Depends(get_db)):
    """Authenticate user and return JWT token."""
    try:
        # Try to find user by username or email (case-insensitive for email)
        user = db.query(User).filter(
            or_(
                User.username == login_data.username,
                User.email.ilike(login_data.username)
            )
        ).first()
        
        if not user or not verify_password(login_data.password, user.hashed_password):
            raise create_error_response("Invalid credentials", status.HTTP_401_UNAUTHORIZED)
        
        if not user.is_active:
            raise create_error_response("User account is disabled", status.HTTP_401_UNAUTHORIZED)
        
        access_token = create_access_token(data={"sub": user.username})
        return {"access_token": access_token, "token_type": "bearer"}
        
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "login_user")
        raise create_error_response("Login failed", status.HTTP_500_INTERNAL_SERVER_ERROR)

# User endpoints
@app.get("/users/me", response_model=UserResponse)
async def get_current_user_info(current_user: User = Depends(get_current_user)):
    """Get current user information."""
    return current_user

@app.put("/users/me", response_model=Union[UserResponse, UserUpdateResponse])
async def update_current_user(
    user_update: UserUpdate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Update current user information."""
    try:
        if user_update.email:
            # Check if email is already taken (case-insensitive)
            existing_user = db.query(User).filter(
                User.email.ilike(user_update.email),
                User.id != current_user.id
            ).first()
            if existing_user:
                raise create_error_response("Email already taken", status.HTTP_400_BAD_REQUEST)
            current_user.email = user_update.email
        
        username_updated = False
        if user_update.username:
            # Check if username is already taken
            existing_user = db.query(User).filter(
                User.username == user_update.username,
                User.id != current_user.id
            ).first()
            if existing_user:
                raise create_error_response("Username already taken", status.HTTP_400_BAD_REQUEST)
            current_user.username = user_update.username
            username_updated = True
        
        if user_update.first_name is not None:
            current_user.first_name = user_update.first_name
        
        if user_update.last_name is not None:
            current_user.last_name = user_update.last_name
        
        if user_update.password:
            current_user.hashed_password = get_password_hash(user_update.password)
        
        db.commit()
        db.refresh(current_user)
        
        # If username was updated, return a new token
        if username_updated:
            new_token = create_access_token(data={"sub": current_user.username})
            return UserUpdateResponse(
                user=current_user,
                new_token=new_token,
                token_type="bearer"
            )
        
        return current_user
        
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "update_current_user")
        raise create_error_response("Update failed", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.put("/users/me/password")
async def change_password(
    password_data: PasswordChangeRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Change user password."""
    try:
        # Verify current password
        if not verify_password(password_data.current_password, current_user.hashed_password):
            raise create_error_response(
                "Current password is incorrect",
                status.HTTP_400_BAD_REQUEST,
                error_code="INVALID_CURRENT_PASSWORD"
            )
        
        # Validate new password strength
        if not validate_password_strength(password_data.new_password):
            raise create_error_response(
                "Password does not meet security requirements. Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.",
                status.HTTP_400_BAD_REQUEST,
                error_code="WEAK_PASSWORD"
            )
        
        # Update password
        current_user.hashed_password = get_password_hash(password_data.new_password)
        db.commit()
        
        return {"message": "Password changed successfully"}
        
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "change_password")
        raise create_error_response("Password change failed", status.HTTP_500_INTERNAL_SERVER_ERROR)


# User preferences endpoints
@app.get("/users/me/preferences", response_model=UserPreferencesResponse)
async def get_user_preferences(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get user preferences."""
    preferences = db.query(UserPreferences).filter(
        UserPreferences.user_id == current_user.id
    ).first()
    
    if not preferences:
        # Create default preferences
        preferences = UserPreferences(user_id=current_user.id)
        db.add(preferences)
        db.commit()
        db.refresh(preferences)
    
    return preferences

@app.put("/users/me/preferences", response_model=UserPreferencesResponse)
async def update_user_preferences(
    preferences_update: UserPreferencesUpdate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Update user preferences."""
    try:
        preferences = db.query(UserPreferences).filter(
            UserPreferences.user_id == current_user.id
        ).first()
        
        if not preferences:
            preferences = UserPreferences(user_id=current_user.id)
            db.add(preferences)
        
        if preferences_update.roundup_enabled is not None:
            preferences.roundup_enabled = preferences_update.roundup_enabled
        
        if preferences_update.roundup_multiplier is not None:
            preferences.roundup_multiplier = preferences_update.roundup_multiplier
        
        if preferences_update.notifications_enabled is not None:
            preferences.notifications_enabled = preferences_update.notifications_enabled
        
        db.commit()
        db.refresh(preferences)
        
        return preferences
        
    except Exception as e:
        log_error(e, "update_user_preferences")
        raise create_error_response("Preferences update failed", status.HTTP_500_INTERNAL_SERVER_ERROR)

# Bank accounts endpoints
@app.get("/users/me/bank-accounts", response_model=List[BankAccountResponse])
async def get_user_bank_accounts(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get user's bank accounts."""
    accounts = db.query(BankAccount).filter(
        BankAccount.user_id == current_user.id,
        BankAccount.is_active == True
    ).all()
    
    return accounts

@app.post("/users/me/bank-accounts", response_model=BankAccountResponse)
async def add_bank_account(
    account_data: BankAccountCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Add a new bank account for the user."""
    try:
        # Verify the account belongs to the current user
        if account_data.user_id != current_user.id:
            raise create_error_response("Invalid user ID", status.HTTP_400_BAD_REQUEST)
        
        # Check if account already exists
        existing_account = db.query(BankAccount).filter(
            BankAccount.plaid_account_id == account_data.plaid_account_id,
            BankAccount.user_id == current_user.id
        ).first()
        
        if existing_account:
            raise create_error_response("Bank account already exists", status.HTTP_400_BAD_REQUEST)
        
        bank_account = BankAccount(**account_data.dict())
        db.add(bank_account)
        db.commit()
        db.refresh(bank_account)
        
        return bank_account
        
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "add_bank_account")
        raise create_error_response("Failed to add bank account", status.HTTP_500_INTERNAL_SERVER_ERROR)

# Health check
@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {"status": "healthy", "service": "user-service"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
