package internal

import (
	"context"
	"database/sql"
	"math"
	"strings"
	"time"

	pb "untether/services/transaction/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DB interface defines the database operations needed by the TransactionCalculator
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// TransactionCalculator handles roundup calculations for transactions
type TransactionCalculator struct {
	pb.UnimplementedTransactionCalculatorServer
	db DB
}

// NewTransactionCalculator creates a new TransactionCalculator instance
func NewTransactionCalculator(db DB) *TransactionCalculator {
	return &TransactionCalculator{
		db: db,
	}
}

// CalculateRoundup calculates the roundup amount for a given transaction
func (c *TransactionCalculator) CalculateRoundup(ctx context.Context, req *pb.CalculateRoundupRequest) (*pb.CalculateRoundupResponse, error) {
	// Normalize the rounding rule
	roundingRule := strings.ToLower(req.RoundingRule)

	// Calculate the rounded amount based on the rule
	var roundedAmount float64
	switch roundingRule {
	case "dollar":
		roundedAmount = math.Ceil(req.Amount)
	case "custom":
		if req.CustomRoundingAmount <= 0 {
			return nil, status.Error(codes.InvalidArgument, "custom rounding amount must be positive")
		}
		roundedAmount = math.Ceil(req.Amount/req.CustomRoundingAmount) * req.CustomRoundingAmount
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid rounding rule")
	}

	// Calculate the roundup amount
	roundupAmount := roundedAmount - req.Amount

	// Store the calculation in the database
	_, err := c.db.ExecContext(ctx,
		`INSERT INTO roundup_calculations 
		(amount, rounding_rule, custom_rounding_amount, rounded_amount, roundup_amount, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		req.Amount, roundingRule, req.CustomRoundingAmount, roundedAmount, roundupAmount, time.Now())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to store calculation")
	}

	return &pb.CalculateRoundupResponse{
		OriginalAmount:   req.Amount,
		RoundedAmount:    roundedAmount,
		RoundupAmount:    roundupAmount,
		RoundingRuleUsed: roundingRule,
	}, nil
}
