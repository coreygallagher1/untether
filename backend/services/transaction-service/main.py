# Transaction service - FastAPI application
from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy.orm import Session
from typing import List, Optional
import sys
import os
import math
from datetime import datetime, timedelta
from decimal import Decimal

# Add shared modules to path
sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', 'shared'))

from shared.database import get_db
from shared.models import User, RoundupCalculation
from shared.schemas import RoundupCalculationResponse, RoundupCalculationCreate, RoundupSummaryResponse, BatchCalculateRequest, BatchCalculateResponse
from shared.auth import get_current_user_token
from shared.config import settings
from shared.utils import create_error_response, log_error

# Create FastAPI app
app = FastAPI(
    title="Transaction Service",
    description="Transaction processing and roundup calculations",
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

# Transaction calculation endpoints
@app.post("/transactions/calculate-roundup", response_model=RoundupCalculationResponse)
async def calculate_roundup(
    calculation_data: RoundupCalculationCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Calculate roundup amount for a transaction."""
    try:
        # Normalize the rounding rule
        rounding_rule = calculation_data.rounding_rule.lower()
        
        # Calculate the rounded amount based on the rule
        amount = float(calculation_data.amount)
        rounded_amount = 0.0
        
        if rounding_rule == "dollar":
            rounded_amount = math.ceil(amount)
        elif rounding_rule == "custom":
            if calculation_data.custom_rounding_amount <= 0:
                raise create_error_response("Custom rounding amount must be positive", status.HTTP_400_BAD_REQUEST)
            rounded_amount = math.ceil(amount / calculation_data.custom_rounding_amount) * calculation_data.custom_rounding_amount
        else:
            raise create_error_response("Invalid rounding rule", status.HTTP_400_BAD_REQUEST)
        
        # Calculate the roundup amount
        roundup_amount = rounded_amount - amount
        
        # Store the calculation in the database
        calculation = RoundupCalculation(
            user_id=current_user.id,
            amount=calculation_data.amount,
            rounding_rule=rounding_rule,
            custom_rounding_amount=calculation_data.custom_rounding_amount,
            rounded_amount=rounded_amount,
            roundup_amount=roundup_amount
        )
        
        db.add(calculation)
        db.commit()
        db.refresh(calculation)
        
        return calculation
        
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "calculate_roundup")
        raise create_error_response("Roundup calculation failed", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.get("/transactions/roundup-history", response_model=List[RoundupCalculationResponse])
async def get_roundup_history(
    limit: int = 50,
    offset: int = 0,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get roundup calculation history for the current user."""
    try:
        calculations = db.query(RoundupCalculation).filter(
            RoundupCalculation.user_id == current_user.id
        ).order_by(RoundupCalculation.created_at.desc()).offset(offset).limit(limit).all()
        
        return calculations
        
    except Exception as e:
        log_error(e, "get_roundup_history")
        raise create_error_response("Failed to get roundup history", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.get("/transactions/roundup-summary")
async def get_roundup_summary(
    days: int = 30,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get roundup summary for the specified period."""
    try:
        # Calculate date range
        end_date = datetime.utcnow()
        start_date = end_date - timedelta(days=days)
        
        # Get calculations for the period
        calculations = db.query(RoundupCalculation).filter(
            RoundupCalculation.user_id == current_user.id,
            RoundupCalculation.created_at >= start_date,
            RoundupCalculation.created_at <= end_date
        ).all()
        
        # Calculate summary statistics
        total_roundup = sum(float(calc.roundup_amount) for calc in calculations)
        total_transactions = len(calculations)
        average_roundup = total_roundup / total_transactions if total_transactions > 0 else 0
        
        # Group by rounding rule
        rule_stats = {}
        for calc in calculations:
            rule = calc.rounding_rule
            if rule not in rule_stats:
                rule_stats[rule] = {"count": 0, "total": 0}
            rule_stats[rule]["count"] += 1
            rule_stats[rule]["total"] += float(calc.roundup_amount)
        
        return {
            "period_days": days,
            "total_roundup": total_roundup,
            "total_transactions": total_transactions,
            "average_roundup": average_roundup,
            "rounding_rule_stats": rule_stats,
            "start_date": start_date.isoformat(),
            "end_date": end_date.isoformat()
        }
        
    except Exception as e:
        log_error(e, "get_roundup_summary")
        raise create_error_response("Failed to get roundup summary", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.post("/transactions/batch-calculate")
async def batch_calculate_roundups(
    transactions: List[dict],
    rounding_rule: str = "dollar",
    custom_rounding_amount: Optional[float] = None,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Calculate roundups for multiple transactions at once."""
    try:
        results = []
        total_roundup = 0.0
        
        for transaction in transactions:
            amount = float(transaction.get("amount", 0))
            if amount <= 0:
                continue
            
            # Calculate roundup
            if rounding_rule.lower() == "dollar":
                rounded_amount = math.ceil(amount)
            elif rounding_rule.lower() == "custom":
                if not custom_rounding_amount or custom_rounding_amount <= 0:
                    raise create_error_response("Custom rounding amount must be positive", status.HTTP_400_BAD_REQUEST)
                rounded_amount = math.ceil(amount / custom_rounding_amount) * custom_rounding_amount
            else:
                raise create_error_response("Invalid rounding rule", status.HTTP_400_BAD_REQUEST)
            
            roundup_amount = rounded_amount - amount
            total_roundup += roundup_amount
            
            # Store calculation
            calculation = RoundupCalculation(
                user_id=current_user.id,
                amount=amount,
                rounding_rule=rounding_rule.lower(),
                custom_rounding_amount=custom_rounding_amount or 0,
                rounded_amount=rounded_amount,
                roundup_amount=roundup_amount
            )
            
            db.add(calculation)
            
            results.append({
                "transaction_id": transaction.get("id"),
                "original_amount": amount,
                "rounded_amount": rounded_amount,
                "roundup_amount": roundup_amount
            })
        
        db.commit()
        
        return {
            "processed_transactions": len(results),
            "total_roundup": total_roundup,
            "results": results
        }
        
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "batch_calculate_roundups")
        raise create_error_response("Batch calculation failed", status.HTTP_500_INTERNAL_SERVER_ERROR)

@app.delete("/transactions/roundup/{calculation_id}")
async def delete_roundup_calculation(
    calculation_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Delete a specific roundup calculation."""
    try:
        calculation = db.query(RoundupCalculation).filter(
            RoundupCalculation.id == calculation_id,
            RoundupCalculation.user_id == current_user.id
        ).first()
        
        if not calculation:
            raise create_error_response("Calculation not found", status.HTTP_404_NOT_FOUND)
        
        db.delete(calculation)
        db.commit()
        
        return {"success": True, "message": "Calculation deleted"}
        
    except HTTPException:
        raise
    except Exception as e:
        log_error(e, "delete_roundup_calculation")
        raise create_error_response("Failed to delete calculation", status.HTTP_500_INTERNAL_SERVER_ERROR)

# Health check
@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {"status": "healthy", "service": "transaction-service"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8003)
