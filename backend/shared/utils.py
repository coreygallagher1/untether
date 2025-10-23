# Common utilities and helpers
import logging
from typing import Any, Dict, Optional
from fastapi import HTTPException, status
from datetime import datetime
import re

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class AppException(Exception):
    """Base exception for application-specific errors."""
    def __init__(self, message: str, error_code: Optional[str] = None):
        self.message = message
        self.error_code = error_code
        super().__init__(self.message)

def create_error_response(
    message: str, 
    status_code: int = status.HTTP_400_BAD_REQUEST,
    error_code: Optional[str] = None
) -> HTTPException:
    """Create a standardized error response."""
    return HTTPException(
        status_code=status_code, 
        detail={
            "message": message,
            "error_code": error_code,
            "timestamp": datetime.utcnow().isoformat()
        }
    )

def log_error(error: Exception, context: str = "", extra_data: Optional[Dict[str, Any]] = None):
    """Log errors with context and structured data."""
    log_data = {
        "error_type": type(error).__name__,
        "error_message": str(error),
        "context": context,
        "timestamp": datetime.utcnow().isoformat()
    }
    
    if extra_data:
        log_data.update(extra_data)
    
    logger.error(f"Error in {context}: {error}", extra=log_data)

def validate_email_format(email: str) -> bool:
    """Validate email format using regex."""
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return re.match(pattern, email) is not None

def sanitize_input(input_string: str) -> str:
    """Sanitize user input by trimming whitespace."""
    if not isinstance(input_string, str):
        return str(input_string)
    return input_string.strip()

def format_currency(amount: float, currency: str = "USD") -> str:
    """Format currency for display."""
    if currency == "USD":
        return f"${amount:.2f}"
    return f"{amount:.2f} {currency}"

def validate_username(username: str) -> bool:
    """Validate username format."""
    if not username or len(username) < 3 or len(username) > 50:
        return False
    
    # Only allow alphanumeric characters and underscores
    pattern = r'^[a-zA-Z0-9_]+$'
    return re.match(pattern, username) is not None

def validate_password_strength(password: str) -> tuple[bool, str]:
    """Validate password strength and return (is_valid, error_message)."""
    if len(password) < 8:
        return False, "Password must be at least 8 characters long"
    
    if not re.search(r'[A-Z]', password):
        return False, "Password must contain at least one uppercase letter"
    
    if not re.search(r'[a-z]', password):
        return False, "Password must contain at least one lowercase letter"
    
    if not re.search(r'\d', password):
        return False, "Password must contain at least one digit"
    
    return True, ""

def safe_float_conversion(value: Any, default: float = 0.0) -> float:
    """Safely convert value to float with default fallback."""
    try:
        return float(value)
    except (ValueError, TypeError):
        return default

def safe_int_conversion(value: Any, default: int = 0) -> int:
    """Safely convert value to int with default fallback."""
    try:
        return int(value)
    except (ValueError, TypeError):
        return default

def mask_sensitive_data(data: str, visible_chars: int = 4) -> str:
    """Mask sensitive data for logging."""
    if len(data) <= visible_chars:
        return "*" * len(data)
    
    return data[:visible_chars] + "*" * (len(data) - visible_chars)
