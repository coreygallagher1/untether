package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cgallagher/Untether/pkg/plaid"
	pb "github.com/cgallagher/Untether/services/user/proto"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	db          *sql.DB
	plaidClient plaid.PlaidClient
}

func NewUserService(db *sql.DB, plaidClient plaid.PlaidClient) *UserService {
	return &UserService{
		db:          db,
		plaidClient: plaidClient,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	if req.Email == "" || req.FirstName == "" || req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "email, first_name, and last_name are required")
	}

	if !isValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("user with this email already exists")
	}

	now := time.Now()
	user := &pb.User{
		Id:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}

	_, err = s.db.ExecContext(ctx,
		"INSERT INTO users (id, email, first_name, last_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Id, user.Email, user.FirstName, user.LastName, user.CreatedAt.AsTime(), user.UpdatedAt.AsTime())
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, fmt.Errorf("invalid UUID")
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
		return nil, fmt.Errorf("failed to get user: %v", err)
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

	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, fmt.Errorf("invalid UUID")
	}

	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.Id).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	now := time.Now()
	_, err = s.db.ExecContext(ctx,
		"UPDATE users SET first_name = $1, last_name = $2, updated_at = $3 WHERE id = $4",
		req.FirstName, req.LastName, now, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return s.GetUser(ctx, &pb.GetUserRequest{Id: req.Id})
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, fmt.Errorf("invalid UUID")
	}

	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.Id).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	result, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %v", err)
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

	// Check if user exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.UserId).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check user existence: %v", err)
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Check if account is already linked
	err = s.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM bank_accounts WHERE user_id = $1 AND plaid_account_id = $2)",
		req.UserId, req.PlaidAccountId,
	).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check account existence: %v", err)
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "bank account already linked")
	}

	// Get account details from Plaid
	accounts, err := s.plaidClient.GetAccounts(ctx, req.PlaidAccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get accounts from Plaid: %v", err)
	}

	// Find the specific account
	var plaidAccount plaid.BankAccount
	found := false
	for _, acc := range accounts {
		if acc.ID == req.PlaidAccountId {
			plaidAccount = acc
			found = true
			break
		}
	}
	if !found {
		return nil, status.Error(codes.NotFound, "bank account not found in Plaid")
	}

	// Get current balance
	balance, err := s.plaidClient.GetBalance(ctx, req.PlaidAccessToken, req.PlaidAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get account balance from Plaid: %v", err)
	}

	account := &pb.BankAccount{
		Id:             uuid.New().String(),
		UserId:         req.UserId,
		PlaidAccountId: req.PlaidAccountId,
		Name:           plaidAccount.Name,
		Type:           strings.ToUpper(plaidAccount.Type),
		Balance:        balance,
		Currency:       plaidAccount.Currency,
		IsActive:       true,
	}

	var createdAt, updatedAt time.Time
	err = s.db.QueryRowContext(ctx,
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

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
