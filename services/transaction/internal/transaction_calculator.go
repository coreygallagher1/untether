package internal

import (
	"context"
	"math"
	"strings"

	pb "untether/services/transaction/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TransactionCalculator handles roundup calculations for transactions
type TransactionCalculator struct {
	pb.UnimplementedTransactionCalculatorServer
}

// NewTransactionCalculator creates a new TransactionCalculator instance
func NewTransactionCalculator() *TransactionCalculator {
	return &TransactionCalculator{}
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
	case "quarter":
		roundedAmount = math.Ceil(req.Amount*4) / 4
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

	return &pb.CalculateRoundupResponse{
		OriginalAmount:   req.Amount,
		RoundedAmount:    roundedAmount,
		RoundupAmount:    roundupAmount,
		RoundingRuleUsed: roundingRule,
	}, nil
}
