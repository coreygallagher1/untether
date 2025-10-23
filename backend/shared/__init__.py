# Shared utilities package
from .config import settings, get_settings
from .database import get_db, engine, SessionLocal, Base
from .auth import (
    verify_password, 
    get_password_hash, 
    create_access_token, 
    verify_token,
    get_current_user_token,
    validate_password_strength
)
from .models import User, UserPreferences, PlaidItem, BankAccount, RoundupCalculation
from .schemas import *
from .utils import (
    create_error_response,
    log_error,
    validate_email_format,
    validate_username,
    validate_password_strength as utils_validate_password_strength,
    sanitize_input,
    format_currency,
    AppException
)

__all__ = [
    # Config
    "settings",
    "get_settings",
    
    # Database
    "get_db",
    "engine", 
    "SessionLocal",
    "Base",
    
    # Auth
    "verify_password",
    "get_password_hash",
    "create_access_token",
    "verify_token",
    "get_current_user_token",
    "validate_password_strength",
    
    # Models
    "User",
    "UserPreferences", 
    "PlaidItem",
    "BankAccount",
    "RoundupCalculation",
    
    # Utils
    "create_error_response",
    "log_error",
    "validate_email_format",
    "validate_username",
    "utils_validate_password_strength",
    "sanitize_input",
    "format_currency",
    "AppException",
]
