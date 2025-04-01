package internal

import (
	"context"
	"database/sql"
	"math"

	pb "untether/services/roundup/proto"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// DB is an interface that matches the methods we need from sql.DB
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type RoundupService struct {
	pb.UnimplementedRoundupServiceServer
	db DB
}

func NewRoundupService(db DB) *RoundupService {
	return &RoundupService{
		db: db,
	}
}

func (s *RoundupService) RoundupTransaction(ctx context.Context, req *pb.RoundupRequest) (*pb.RoundupResponse, error) {
	// Calculate roundup amount (round up to nearest dollar)
	amount := req.TransactionAmount
	rounded := math.Ceil(amount)
	roundupAmount := rounded - amount

	// Create roundup record
	roundupID := uuid.New().String()
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	// Store roundup in database
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO roundups (id, user_id, transaction_amount, roundup_amount, status)
		VALUES ($1, $2, $3, $4, $5)
	`, roundupID, userID, amount, roundupAmount, "pending")
	if err != nil {
		return nil, err
	}

	return &pb.RoundupResponse{
		RoundupId:     roundupID,
		AmountRounded: roundupAmount,
	}, nil
}
