package plaid

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockPlaidClient implements the PlaidClient interface for testing
type MockPlaidClient struct {
	mock.Mock
}

// NewMockPlaidClient creates a new mock Plaid client
func NewMockPlaidClient() *MockPlaidClient {
	return &MockPlaidClient{}
}

// CreateLinkToken creates a mock link token
func (m *MockPlaidClient) CreateLinkToken(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

// ExchangePublicToken exchanges a mock public token for an access token
func (m *MockPlaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (string, error) {
	args := m.Called(ctx, publicToken)
	return args.String(0), args.Error(1)
}

// GetAccounts retrieves mock account information
func (m *MockPlaidClient) GetAccounts(ctx context.Context, accessToken string) ([]BankAccount, error) {
	args := m.Called(ctx, accessToken)
	if accounts, ok := args.Get(0).([]BankAccount); ok {
		return accounts, args.Error(1)
	}
	return nil, args.Error(1)
}

// GetBalance retrieves mock balance information
func (m *MockPlaidClient) GetBalance(ctx context.Context, accessToken string, accountID string) (float64, error) {
	args := m.Called(ctx, accessToken, accountID)
	return args.Get(0).(float64), args.Error(1)
}
