package tests

import (
	"context"
	"testing"

	plaidinternal "github.com/cgallagher/Untether/services/plaid/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlaidClient_CreateLinkToken(t *testing.T) {
	// Create a new Plaid client with test configuration
	client := plaidinternal.NewPlaidClient("test-client-id", "test-secret", "sandbox")
	require.NotNil(t, client)

	// Test creating a link token
	ctx := context.Background()
	userID := "test-user-id"
	token, err := client.CreateLinkToken(ctx, userID)

	// In sandbox mode, we expect an error because we don't have real credentials
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestPlaidClient_ExchangePublicToken(t *testing.T) {
	// Create a new Plaid client with test configuration
	client := plaidinternal.NewPlaidClient("test-client-id", "test-secret", "sandbox")
	require.NotNil(t, client)

	// Test exchanging a public token
	ctx := context.Background()
	publicToken := "test-public-token"
	accessToken, err := client.ExchangePublicToken(ctx, publicToken)

	// In sandbox mode, we expect an error because we don't have a real public token
	assert.Error(t, err)
	assert.Empty(t, accessToken)
}

func TestPlaidClient_GetAccounts(t *testing.T) {
	// Create a new Plaid client with test configuration
	client := plaidinternal.NewPlaidClient("test-client-id", "test-secret", "sandbox")
	require.NotNil(t, client)

	// Test getting accounts
	ctx := context.Background()
	accessToken := "test-access-token"
	accounts, err := client.GetAccounts(ctx, accessToken)

	// In sandbox mode, we expect an error because we don't have a real access token
	assert.Error(t, err)
	assert.Empty(t, accounts)
}

func TestPlaidClient_GetBalance(t *testing.T) {
	// Create a new Plaid client with test configuration
	client := plaidinternal.NewPlaidClient("test-client-id", "test-secret", "sandbox")
	require.NotNil(t, client)

	// Test getting balance
	ctx := context.Background()
	accessToken := "test-access-token"
	accountID := "test-account-id"
	balance, err := client.GetBalance(ctx, accessToken, accountID)

	// In sandbox mode, we expect an error because we don't have a real access token
	assert.Error(t, err)
	assert.Zero(t, balance)
}

func TestPlaidClient_InvalidEnvironment(t *testing.T) {
	// Test creating a client with an invalid environment
	client := plaidinternal.NewPlaidClient("test-client-id", "test-secret", "invalid")
	assert.Nil(t, client)
}
