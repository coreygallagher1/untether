package plaid

import (
	"context"

	"github.com/plaid/plaid-go/v14/plaid"
)

// MockClient implements the PlaidClient interface for testing
type MockClient struct {
	accounts map[string]*plaid.AccountBase
}

// NewMockClient creates a new mock Plaid client
func NewMockClient() PlaidClient {
	return &MockClient{
		accounts: make(map[string]*plaid.AccountBase),
	}
}

// AddMockAccount adds a mock account for testing
func (m *MockClient) AddMockAccount(accountID string, account *plaid.AccountBase) {
	m.accounts[accountID] = account
}

// GetAccountDetails returns mock account details
func (m *MockClient) GetAccountDetails(ctx context.Context, accessToken, accountID string) (*plaid.AccountBase, error) {
	if account, exists := m.accounts[accountID]; exists {
		return account, nil
	}

	// Create a default mock account
	account := &plaid.AccountBase{}
	account.AccountId = accountID
	account.Name = "Mock Checking Account"
	account.Type = plaid.ACCOUNTTYPE_DEPOSITORY

	balances := &plaid.AccountBalance{}
	current := plaid.NewNullableFloat64(plaid.PtrFloat64(1000.00))
	currency := plaid.NewNullableString(plaid.PtrString("USD"))
	balances.Current = *current
	balances.IsoCurrencyCode = *currency
	account.Balances = *balances

	return account, nil
}

// CreateLinkToken returns a mock link token
func (m *MockClient) CreateLinkToken(ctx context.Context, clientUserID, clientName string) (string, error) {
	return "mock-link-token", nil
}

// ExchangePublicToken returns a mock access token
func (m *MockClient) ExchangePublicToken(ctx context.Context, publicToken string) (string, error) {
	return "mock-access-token", nil
}
