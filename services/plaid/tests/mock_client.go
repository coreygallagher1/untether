package tests

import (
	"context"

	"github.com/cgallagher/Untether/pkg/plaid"
	"github.com/stretchr/testify/mock"
)

type MockPlaidClient struct {
	mock.Mock
}

func (m *MockPlaidClient) CreateLinkToken(ctx context.Context, userId string) (string, error) {
	args := m.Called(ctx, userId)
	return args.String(0), args.Error(1)
}

func (m *MockPlaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (string, error) {
	args := m.Called(ctx, publicToken)
	return args.String(0), args.Error(1)
}

func (m *MockPlaidClient) GetAccounts(ctx context.Context, accessToken string) ([]plaid.BankAccount, error) {
	args := m.Called(ctx, accessToken)
	if accounts, ok := args.Get(0).([]plaid.BankAccount); ok {
		return accounts, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPlaidClient) GetBalance(ctx context.Context, accessToken string, accountId string) (float64, error) {
	args := m.Called(ctx, accessToken, accountId)
	return args.Get(0).(float64), args.Error(1)
}
