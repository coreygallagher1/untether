package unit

import (
	"context"
	"testing"

	"untether/services/plaid/client"

	"github.com/stretchr/testify/assert"
)

// MockPlaidClient implements the PlaidClient interface for testing
type MockPlaidClient struct{}

func (m *MockPlaidClient) CreateLinkToken(ctx context.Context, userID string) (string, error) {
	return "mock-link-token", nil
}

func (m *MockPlaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (string, error) {
	return "mock-access-token", nil
}

func (m *MockPlaidClient) GetAccounts(ctx context.Context, accessToken string) ([]client.BankAccount, error) {
	return []client.BankAccount{
		{
			AccountID: "mock-account-id",
			Name:      "Mock Account",
			Type:      "checking",
			Subtype:   "personal",
			Mask:      "1234",
		},
	}, nil
}

func (m *MockPlaidClient) GetBalance(ctx context.Context, accessToken string, accountID string) (float64, error) {
	return 1000.00, nil
}

func TestPlaidClient(t *testing.T) {
	mockClient := &MockPlaidClient{}

	t.Run("CreateLinkToken", func(t *testing.T) {
		token, err := mockClient.CreateLinkToken(context.Background(), "test-user")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("ExchangePublicToken", func(t *testing.T) {
		accessToken, err := mockClient.ExchangePublicToken(context.Background(), "public-token")
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
	})

	t.Run("GetAccounts", func(t *testing.T) {
		accounts, err := mockClient.GetAccounts(context.Background(), "access-token")
		assert.NoError(t, err)
		assert.NotEmpty(t, accounts)
	})

	t.Run("GetBalance", func(t *testing.T) {
		balance, err := mockClient.GetBalance(context.Background(), "access-token", "account-id")
		assert.NoError(t, err)
		assert.NotZero(t, balance)
	})
}
