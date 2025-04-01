package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cgallagher/Untether/pkg/plaid"
	userinternal "github.com/cgallagher/Untether/services/user/internal"
	pb "github.com/cgallagher/Untether/services/user/proto"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory SQLite database and initializes the schema
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE bank_accounts (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			plaid_account_id TEXT NOT NULL,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			balance REAL NOT NULL,
			currency TEXT NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`)
	require.NoError(t, err)

	return db
}

// setupMockPlaid creates a mock Plaid client with common expectations
func setupMockPlaid() *plaid.MockPlaidClient {
	mockPlaid := plaid.NewMockPlaidClient()

	// Set up common mock expectations
	mockPlaid.On("GetAccounts", mock.Anything, "mock-access-token").Return([]plaid.BankAccount{
		{
			ID:       "mock-account-id",
			Name:     "Mock Checking Account",
			Type:     "DEPOSITORY",
			Balance:  1000.00,
			Currency: "USD",
		},
	}, nil)

	mockPlaid.On("GetBalance", mock.Anything, "mock-access-token", "mock-account-id").Return(1000.00, nil)
	mockPlaid.On("CreateLinkToken", mock.Anything, "test-user").Return("mock-link-token", nil)
	mockPlaid.On("ExchangePublicToken", mock.Anything, "mock-public-token").Return("mock-access-token", nil)

	return mockPlaid
}

// setupMockPlaidWithMultipleAccounts creates a mock Plaid client with expectations for multiple accounts
func setupMockPlaidWithMultipleAccounts() *plaid.MockPlaidClient {
	mockPlaid := plaid.NewMockPlaidClient()

	// Set up mock expectations for multiple accounts
	mockPlaid.On("GetAccounts", mock.Anything, "mock-access-token").Return([]plaid.BankAccount{
		{
			ID:       "account-1",
			Name:     "Mock Checking Account",
			Type:     "DEPOSITORY",
			Balance:  1000.00,
			Currency: "USD",
		},
		{
			ID:       "account-2",
			Name:     "Mock Checking Account",
			Type:     "DEPOSITORY",
			Balance:  1000.00,
			Currency: "USD",
		},
		{
			ID:       "account-3",
			Name:     "Mock Checking Account",
			Type:     "DEPOSITORY",
			Balance:  1000.00,
			Currency: "USD",
		},
	}, nil)

	mockPlaid.On("GetBalance", mock.Anything, "mock-access-token", "account-1").Return(1000.00, nil)
	mockPlaid.On("GetBalance", mock.Anything, "mock-access-token", "account-2").Return(1000.00, nil)
	mockPlaid.On("GetBalance", mock.Anything, "mock-access-token", "account-3").Return(1000.00, nil)

	return mockPlaid
}

// createTestUser creates a test user in the database
func createTestUser(ctx context.Context, userService *userinternal.UserService) (*pb.User, error) {
	return userService.CreateUser(ctx, &pb.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	})
}

func TestUserService_LinkBankAccount(t *testing.T) {
	db := setupTestDB(t)
	mockPlaid := setupMockPlaid()
	userService := userinternal.NewUserService(db, mockPlaid)

	ctx := context.Background()
	user, err := createTestUser(ctx, userService)
	require.NoError(t, err)
	require.NotNil(t, user)

	// Test linking a bank account
	bankAccount, err := userService.LinkBankAccount(ctx, &pb.LinkBankAccountRequest{
		UserId:           user.Id,
		PlaidAccessToken: "mock-access-token",
		PlaidAccountId:   "mock-account-id",
	})
	require.NoError(t, err)
	require.NotNil(t, bankAccount)

	// Verify bank account details
	assert.Equal(t, "Mock Checking Account", bankAccount.Name)
	assert.Equal(t, "DEPOSITORY", bankAccount.Type)
	assert.Equal(t, float64(1000.00), bankAccount.Balance)
	assert.Equal(t, "USD", bankAccount.Currency)
	assert.NotEmpty(t, bankAccount.Id)
	assert.Equal(t, user.Id, bankAccount.UserId)
	assert.Equal(t, "mock-account-id", bankAccount.PlaidAccountId)
	assert.True(t, bankAccount.IsActive)
	assert.NotZero(t, bankAccount.CreatedAt)
	assert.NotZero(t, bankAccount.UpdatedAt)
}

func TestPlaidTokenOperations(t *testing.T) {
	mockPlaid := setupMockPlaid()
	ctx := context.Background()

	// Test creating a link token
	linkToken, err := mockPlaid.CreateLinkToken(ctx, "test-user")
	require.NoError(t, err)
	assert.Equal(t, "mock-link-token", linkToken)

	// Test exchanging public token
	accessToken, err := mockPlaid.ExchangePublicToken(ctx, "mock-public-token")
	require.NoError(t, err)
	assert.Equal(t, "mock-access-token", accessToken)
}

func TestBankAccountLinking_ErrorCases(t *testing.T) {
	db := setupTestDB(t)
	mockPlaid := setupMockPlaid()
	userService := userinternal.NewUserService(db, mockPlaid)

	ctx := context.Background()
	user, err := createTestUser(ctx, userService)
	require.NoError(t, err)
	require.NotNil(t, user)

	// Test case 1: Missing required fields
	_, err = userService.LinkBankAccount(ctx, &pb.LinkBankAccountRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user_id, plaid_access_token, and plaid_account_id are required")

	// Test case 2: Invalid user ID
	_, err = userService.LinkBankAccount(ctx, &pb.LinkBankAccountRequest{
		UserId:           "invalid-user-id",
		PlaidAccessToken: "mock-access-token",
		PlaidAccountId:   "mock-account-id",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	// Test case 3: Link same account twice
	// First link should succeed
	account, err := userService.LinkBankAccount(ctx, &pb.LinkBankAccountRequest{
		UserId:           user.Id,
		PlaidAccessToken: "mock-access-token",
		PlaidAccountId:   "mock-account-id",
	})
	require.NoError(t, err)
	require.NotNil(t, account)

	// Second link should fail
	_, err = userService.LinkBankAccount(ctx, &pb.LinkBankAccountRequest{
		UserId:           user.Id,
		PlaidAccessToken: "mock-access-token",
		PlaidAccountId:   "mock-account-id",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bank account already linked")
}

func TestListBankAccounts(t *testing.T) {
	db := setupTestDB(t)
	mockPlaid := setupMockPlaidWithMultipleAccounts()
	userService := userinternal.NewUserService(db, mockPlaid)

	ctx := context.Background()
	user, err := createTestUser(ctx, userService)
	require.NoError(t, err)
	require.NotNil(t, user)

	// Link multiple bank accounts
	accountIDs := []string{"account-1", "account-2", "account-3"}
	for _, accountID := range accountIDs {
		account, err := userService.LinkBankAccount(ctx, &pb.LinkBankAccountRequest{
			UserId:           user.Id,
			PlaidAccessToken: "mock-access-token",
			PlaidAccountId:   accountID,
		})
		require.NoError(t, err)
		require.NotNil(t, account)
	}

	// Test listing bank accounts
	response, err := userService.ListBankAccounts(ctx, &pb.ListBankAccountsRequest{
		UserId: user.Id,
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Len(t, response.Accounts, len(accountIDs))

	// Verify each account
	for _, account := range response.Accounts {
		assert.NotEmpty(t, account.Id)
		assert.Equal(t, user.Id, account.UserId)
		assert.Contains(t, accountIDs, account.PlaidAccountId)
		assert.Equal(t, "Mock Checking Account", account.Name)
		assert.Equal(t, "DEPOSITORY", account.Type)
		assert.Equal(t, float64(1000.00), account.Balance)
		assert.Equal(t, "USD", account.Currency)
		assert.True(t, account.IsActive)
		assert.NotZero(t, account.CreatedAt)
		assert.NotZero(t, account.UpdatedAt)
	}

	// Test listing with invalid user ID
	invalidResponse, err := userService.ListBankAccounts(ctx, &pb.ListBankAccountsRequest{
		UserId: "invalid-user-id",
	})
	require.NoError(t, err)
	assert.Empty(t, invalidResponse.Accounts)

	// Test listing with missing user ID
	_, err = userService.ListBankAccounts(ctx, &pb.ListBankAccountsRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user_id is required")
}
