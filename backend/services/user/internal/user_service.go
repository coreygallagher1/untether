package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"untether/services/plaid/client"
	pb "untether/services/user/proto"
)

// DB interface defines the database operations needed by the UserService
type DB interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// UserService implements the UserServiceServer interface
type UserService struct {
	pb.UnimplementedUserServiceServer
	db          DB
	plaidClient client.PlaidClient
}

func (s *UserService) mustEmbedUnimplementedUserServiceServer() {}

// NewUserService creates a new UserService
func NewUserService(db DB, plaidClient client.PlaidClient) *UserService {
	return &UserService{
		db:          db,
		plaidClient: plaidClient,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	// Validate required fields
	if req.Email == "" || req.FirstName == "" || req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "email, first_name, and last_name are required")
	}

	if !isValidEmail(req.Email) {
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}

	// Check for existing user
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing user: %v", err)
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	}

	// Create user
	now := time.Now()
	user := &pb.User{
		Id:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}

	// Insert user into database
	_, err = s.db.ExecContext(ctx,
		"INSERT INTO users (id, email, first_name, last_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Id, user.Email, user.FirstName, user.LastName, user.CreatedAt.AsTime(), user.UpdatedAt.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	// Validate required fields
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Validate UUID format
	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, fmt.Errorf("invalid UUID")
	}

	var user pb.User
	var createdAt, updatedAt time.Time

	// Query user from database
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

	// Convert timestamps to protobuf format
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	// Query bank accounts from database
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, plaid_account_id, name, type, balance, currency, is_active, created_at, updated_at 
		 FROM bank_accounts WHERE user_id = $1`,
		user.Id,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get bank accounts: %v", err)
	}
	defer rows.Close()

	// Iterate over bank accounts
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
	// Validate required fields
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, fmt.Errorf("invalid UUID")
	}

	// Check if user exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.Id).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Update user fields
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
	// Validate required fields
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Validate UUID format
	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, fmt.Errorf("invalid UUID")
	}

	// Check if user exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.Id).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Delete user from database
	result, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %v", err)
	}

	// Get rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.DeleteUserResponse{
		Success: rowsAffected > 0,
	}, nil
}

func (s *UserService) CreateUserPreferences(ctx context.Context, req *pb.CreateUserPreferencesRequest) (*pb.UserPreferences, error) {
	// Validate required fields
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	validCurrencies := []string{"USD"}
	if !slices.Contains(validCurrencies, req.Currency) {
		return nil, status.Error(codes.InvalidArgument, "Invalid currency")
	}

	validTimezones := []string{"UTC"}
	if !slices.Contains(validTimezones, req.Timezone) {
		return nil, status.Error(codes.InvalidArgument, "Invalid timezone")
	}

	validLanguages := []string{"en"}
	if !slices.Contains(validLanguages, req.Language) {
		return nil, status.Error(codes.InvalidArgument, "Invalid language")
	}

	if req.Budget <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Budget must be greater than 0")
	}

	// Check if user exists
	var userExists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.UserId).Scan(&userExists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check user existence: %v", err)
	}
	if !userExists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	// Check if preferences already exist
	var preferencesExists bool
	err = s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM user_preferences WHERE user_id = $1)", req.UserId).Scan(&preferencesExists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check user preferences existence: %v", err)
	}
	if preferencesExists {
		return nil, status.Error(codes.AlreadyExists, "User preferences already exist")
	}

	// Create user preferences
	now := time.Now()
	userPreferences := &pb.UserPreferences{
		Currency: req.Currency,
		Timezone: req.Timezone,
		Language: req.Language,
		DarkMode: req.DarkMode,
		Budget:   req.Budget,
	}

	// Insert user preferences into database
	_, err = s.db.ExecContext(ctx,
		"INSERT INTO user_preferences (user_id, currency, timezone, language, dark_mode, budget, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		req.UserId, userPreferences.Currency, userPreferences.Timezone, userPreferences.Language, userPreferences.DarkMode, userPreferences.Budget, now, now,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user preferences: %v", err)
	}

	return userPreferences, nil
}

func (s *UserService) LinkBankAccount(ctx context.Context, req *pb.LinkBankAccountRequest) (*pb.BankAccount, error) {
	// Validate required fields
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

	// Check if Plaid client is available
	if s.plaidClient == nil {
		return nil, status.Error(codes.Unavailable, "Plaid service is not available")
	}

	// Get account details from Plaid
	accounts, err := s.plaidClient.GetAccounts(ctx, req.PlaidAccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get accounts from Plaid: %v", err)
	}

	// Find the specific account
	var plaidAccount client.BankAccount
	found := false
	for _, acc := range accounts {
		if acc.AccountID == req.PlaidAccountId {
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

	// Create bank account
	account := &pb.BankAccount{
		Id:             uuid.New().String(),
		UserId:         req.UserId,
		PlaidAccountId: req.PlaidAccountId,
		Name:           plaidAccount.Name,
		Type:           plaidAccount.Type,
		Balance:        balance,
		Currency:       "USD", // Default to USD for now
		IsActive:       true,
		CreatedAt:      timestamppb.Now(),
		UpdatedAt:      timestamppb.Now(),
	}

	// Store account in database
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO bank_accounts (id, user_id, plaid_account_id, name, type, balance, currency, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		account.Id, account.UserId, account.PlaidAccountId, account.Name, account.Type,
		account.Balance, account.Currency, account.IsActive,
		account.CreatedAt.AsTime(), account.UpdatedAt.AsTime(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store bank account: %v", err)
	}

	return account, nil
}

func (s *UserService) ListBankAccounts(ctx context.Context, req *pb.ListBankAccountsRequest) (*pb.ListBankAccountsResponse, error) {
	// Validate required fields
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Validate UUID format
	if _, err := uuid.Parse(req.UserId); err != nil {
		return nil, fmt.Errorf("invalid UUID")
	}

	// Get user's access token from database
	accessToken, err := s.getUserAccessToken(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// Get accounts from Plaid
	accounts, err := s.plaidClient.GetAccounts(ctx, accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get accounts: %v", err)
	}

	// Create bank account list
	var pbAccounts []*pb.BankAccount
	for _, acc := range accounts {
		balance, err := s.plaidClient.GetBalance(ctx, accessToken, acc.AccountID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get balance: %v", err)
		}

		// Create bank account
		pbAccounts = append(pbAccounts, &pb.BankAccount{
			Id:             acc.AccountID,
			UserId:         req.UserId,
			PlaidAccountId: acc.AccountID,
			Name:           acc.Name,
			Type:           acc.Type,
			Balance:        balance,
			Currency:       "USD", // Default to USD for now
			IsActive:       true,
			CreatedAt:      timestamppb.Now(),
			UpdatedAt:      timestamppb.Now(),
		})
	}

	return &pb.ListBankAccountsResponse{
		Accounts: pbAccounts,
	}, nil
}

// getUserAccessToken retrieves the Plaid access token for a user from the database
func (s *UserService) getUserAccessToken(ctx context.Context, userID string) (string, error) {
	// Validate UUID
	if _, err := uuid.Parse(userID); err != nil {
		return "", status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	// TODO: Implement database lookup for user's Plaid access token
	// For now, return a mock token
	return "access-sandbox-123", nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
