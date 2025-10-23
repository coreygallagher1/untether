# Pydantic schemas for API requests/responses
from pydantic import BaseModel, EmailStr, Field, validator
from typing import Optional, List
from datetime import datetime
from decimal import Decimal

# User schemas
class UserBase(BaseModel):
    email: EmailStr
    username: str = Field(..., min_length=3, max_length=50, pattern=r'^[a-zA-Z0-9_]+$')
    first_name: Optional[str] = Field(None, min_length=1, max_length=100)
    last_name: Optional[str] = Field(None, min_length=1, max_length=100)

class UserCreate(UserBase):
    password: str = Field(..., min_length=8, max_length=100)

class UserUpdate(BaseModel):
    email: Optional[EmailStr] = None
    username: Optional[str] = Field(None, min_length=3, max_length=50, pattern=r'^[a-zA-Z0-9_]+$')
    first_name: Optional[str] = Field(None, min_length=1, max_length=100)
    last_name: Optional[str] = Field(None, min_length=1, max_length=100)
    password: Optional[str] = Field(None, min_length=8, max_length=100)

class UserResponse(UserBase):
    id: int
    is_active: bool
    is_verified: bool
    created_at: datetime
    
    class Config:
        from_attributes = True

class UserUpdateResponse(BaseModel):
    user: UserResponse
    new_token: str
    token_type: str = "bearer"

# Auth schemas
class Token(BaseModel):
    access_token: str
    token_type: str = "bearer"

class TokenData(BaseModel):
    username: Optional[str] = None

class LoginRequest(BaseModel):
    username: str = Field(..., min_length=3, max_length=50)
    password: str = Field(..., min_length=1)

class PasswordChangeRequest(BaseModel):
    current_password: str = Field(..., min_length=1)
    new_password: str = Field(..., min_length=8, max_length=100)

# User preferences schemas
class UserPreferencesBase(BaseModel):
    roundup_enabled: bool = True
    roundup_multiplier: int = Field(1, ge=1, le=10)
    notifications_enabled: bool = True

class UserPreferencesCreate(UserPreferencesBase):
    user_id: int

class UserPreferencesUpdate(BaseModel):
    roundup_enabled: Optional[bool] = None
    roundup_multiplier: Optional[int] = Field(None, ge=1, le=10)
    notifications_enabled: Optional[bool] = None

class UserPreferencesResponse(UserPreferencesBase):
    id: int
    user_id: int
    created_at: datetime
    updated_at: Optional[datetime] = None
    
    class Config:
        from_attributes = True

# Bank account schemas
class BankAccountBase(BaseModel):
    account_name: str = Field(..., min_length=1, max_length=255)
    account_type: str = Field(..., min_length=1, max_length=50)
    account_subtype: str = Field(..., min_length=1, max_length=50)
    mask: Optional[str] = Field(None, max_length=10)

class BankAccountCreate(BankAccountBase):
    user_id: int
    plaid_item_id: int
    plaid_account_id: str = Field(..., min_length=1, max_length=255)

class BankAccountResponse(BankAccountBase):
    id: int
    user_id: int
    plaid_item_id: int
    plaid_account_id: str
    is_active: bool
    created_at: datetime
    updated_at: Optional[datetime] = None
    
    class Config:
        from_attributes = True

# Plaid schemas
class PlaidItemResponse(BaseModel):
    id: int
    user_id: int
    item_id: str
    institution_id: str
    institution_name: Optional[str]
    status: str
    created_at: datetime
    updated_at: Optional[datetime] = None
    
    class Config:
        from_attributes = True

class LinkTokenResponse(BaseModel):
    link_token: str

class ExchangeTokenRequest(BaseModel):
    public_token: str = Field(..., min_length=1)

class ExchangeTokenResponse(BaseModel):
    access_token: str
    success: bool

class PlaidAccount(BaseModel):
    account_id: str
    name: str
    type: str
    subtype: str
    mask: Optional[str] = None
    balance: float

class PlaidAccountsResponse(BaseModel):
    accounts: List[PlaidAccount]

class PlaidBalanceResponse(BaseModel):
    account_id: str
    balance: float

# Roundup calculation schemas
class RoundupCalculationBase(BaseModel):
    amount: Decimal = Field(..., gt=0)
    rounding_rule: str = Field(..., pattern=r'^(dollar|custom)$')
    custom_rounding_amount: Optional[Decimal] = Field(None, gt=0)
    transaction_id: Optional[str] = Field(None, max_length=255)

    @validator('custom_rounding_amount')
    def validate_custom_rounding_amount(cls, v, values):
        if values.get('rounding_rule') == 'custom' and v is None:
            raise ValueError('custom_rounding_amount is required when rounding_rule is custom')
        return v

class RoundupCalculationCreate(RoundupCalculationBase):
    pass

class RoundupCalculationResponse(RoundupCalculationBase):
    id: int
    user_id: int
    rounded_amount: Decimal
    roundup_amount: Decimal
    created_at: datetime
    
    class Config:
        from_attributes = True

class RoundupSummaryResponse(BaseModel):
    period_days: int
    total_roundup: float
    total_transactions: int
    average_roundup: float
    rounding_rule_stats: dict
    start_date: str
    end_date: str

class BatchCalculateRequest(BaseModel):
    transactions: List[dict] = Field(..., min_items=1)
    rounding_rule: str = Field("dollar", pattern=r'^(dollar|custom)$')
    custom_rounding_amount: Optional[Decimal] = Field(None, gt=0)

class BatchCalculateResponse(BaseModel):
    processed_transactions: int
    total_roundup: float
    results: List[dict]

# Error schemas
class ErrorResponse(BaseModel):
    detail: str
    error_code: Optional[str] = None
    timestamp: datetime = Field(default_factory=datetime.utcnow)
