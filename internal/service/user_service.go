package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "untether/internal/proto"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	if req.Email == "" || req.FirstName == "" || req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "email, first_name, and last_name are required")
	}

	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing user: %v", err)
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	}

	var user pb.User
	var createdAt, updatedAt time.Time

	err = s.db.QueryRowContext(ctx,
		`INSERT INTO users (email, first_name, last_name) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, email, first_name, last_name, created_at, updated_at`,
		req.Email, req.FirstName, req.LastName,
	).Scan(
		&user.Id, &user.Email, &user.FirstName, &user.LastName,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	return &user, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	var user pb.User
	var createdAt, updatedAt time.Time

	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, first_name, last_name, created_at, updated_at 
		 FROM users WHERE id = $1`,
		req.Id,
	).Scan(
		&user.Id, &user.Email, &user.FirstName, &user.LastName,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, plaid_account_id, name, type, balance, currency, is_active, created_at, updated_at 
		 FROM bank_accounts WHERE user_id = $1`,
		user.Id,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get bank accounts: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account pb.BankAccount
		var accCreatedAt, accUpdatedAt time.Time
		err := rows.Scan(
			&account.Id, &account.PlaidAccountId, &account.Name, &account.Type,
			&account.Balance, &account.Currency, &account.IsActive,
			&accCreatedAt, &accUpdatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan bank account: %v", err)
		}
		account.UserId = user.Id
		account.CreatedAt = timestamppb.New(accCreatedAt)
		account.UpdatedAt = timestamppb.New(accUpdatedAt)
		user.BankAccounts = append(user.BankAccounts, &account)
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	var user pb.User
	var createdAt, updatedAt time.Time

	err := s.db.QueryRowContext(ctx,
		`UPDATE users 
		 SET first_name = $1, last_name = $2, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $3 
		 RETURNING id, email, first_name, last_name, created_at, updated_at`,
		req.FirstName, req.LastName, req.Id,
	).Scan(
		&user.Id, &user.Email, &user.FirstName, &user.LastName,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	return &user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	result, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get rows affected: %v", err)
	}

	return &pb.DeleteUserResponse{
		Success: rowsAffected > 0,
	}, nil
}

func (s *UserService) LinkBankAccount(ctx context.Context, req *pb.LinkBankAccountRequest) (*pb.BankAccount, error) {
	if req.UserId == "" || req.PlaidAccessToken == "" || req.PlaidAccountId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, plaid_access_token, and plaid_account_id are required")
	}

	// TODO: Use Plaid API to get account details
	// For now, we'll create a placeholder account
	account := &pb.BankAccount{
		Id:             uuid.New().String(),
		UserId:         req.UserId,
		PlaidAccountId: req.PlaidAccountId,
		Name:           "Checking Account", // This should come from Plaid
		Type:           "checking",        // This should come from Plaid
		Balance:        0.0,              // This should come from Plaid
		Currency:       "USD",
		IsActive:       true,
	}

	var createdAt, updatedAt time.Time
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO bank_accounts 
		 (id, user_id, plaid_account_id, name, type, balance, currency, is_active) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		 RETURNING created_at, updated_at`,
		account.Id, account.UserId, account.PlaidAccountId, account.Name,
		account.Type, account.Balance, account.Currency, account.IsActive,
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create bank account: %v", err)
	}

	account.CreatedAt = timestamppb.New(createdAt)
	account.UpdatedAt = timestamppb.New(updatedAt)

	return account, nil
}

func (s *UserService) ListBankAccounts(ctx context.Context, req *pb.ListBankAccountsRequest) (*pb.ListBankAccountsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, plaid_account_id, name, type, balance, currency, is_active, created_at, updated_at 
		 FROM bank_accounts WHERE user_id = $1`,
		req.UserId,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get bank accounts: %v", err)
	}
	defer rows.Close()

	var accounts []*pb.BankAccount
	for rows.Next() {
		var account pb.BankAccount
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&account.Id, &account.PlaidAccountId, &account.Name, &account.Type,
			&account.Balance, &account.Currency, &account.IsActive,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan bank account: %v", err)
		}
		account.UserId = req.UserId
		account.CreatedAt = timestamppb.New(createdAt)
		account.UpdatedAt = timestamppb.New(updatedAt)
		accounts = append(accounts, &account)
	}

	return &pb.ListBankAccountsResponse{
		Accounts: accounts,
	}, nil
} 