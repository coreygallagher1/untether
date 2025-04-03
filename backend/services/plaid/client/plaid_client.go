package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	plaidSDK "github.com/plaid/plaid-go/v14/plaid"
)

// PlaidClient defines the interface for Plaid API interactions
type PlaidClient interface {
	CreateLinkToken(ctx context.Context, userId string) (string, error)
	ExchangePublicToken(ctx context.Context, publicToken string) (string, error)
	GetAccounts(ctx context.Context, accessToken string) ([]BankAccount, error)
	GetBalance(ctx context.Context, accessToken string, accountId string) (float64, error)
}

// BankAccount represents a bank account from Plaid
type BankAccount struct {
	AccountID string
	Name      string
	Type      string
	Subtype   string
	Mask      string
}

type plaidClient struct {
	clientID    string
	secret      string
	environment string
	httpClient  *http.Client
	sdkClient   *plaidSDK.APIClient
}

func NewPlaidClient(clientID, secret, environment string) PlaidClient {
	configuration := plaidSDK.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)

	// Set the environment
	var env plaidSDK.Environment
	switch environment {
	case "sandbox":
		env = plaidSDK.Sandbox
	case "development":
		env = plaidSDK.Development
	case "production":
		env = plaidSDK.Production
	default:
		return nil
	}
	configuration.UseEnvironment(env)

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	// Create Plaid client
	sdkClient := plaidSDK.NewAPIClient(configuration)

	return &plaidClient{
		clientID:    clientID,
		secret:      secret,
		environment: environment,
		httpClient:  httpClient,
		sdkClient:   sdkClient,
	}
}

// CreateLinkToken creates a link token for initializing Plaid Link
func (c *plaidClient) CreateLinkToken(ctx context.Context, userId string) (string, error) {
	user := plaidSDK.LinkTokenCreateRequestUser{
		ClientUserId: userId,
	}

	request := plaidSDK.NewLinkTokenCreateRequest(
		"Untether",
		"en",
		[]plaidSDK.CountryCode{plaidSDK.COUNTRYCODE_US},
		user,
	)
	request.SetProducts([]plaidSDK.Products{plaidSDK.PRODUCTS_AUTH, plaidSDK.PRODUCTS_TRANSACTIONS})

	response, _, err := c.sdkClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to create link token: %v", err)
	}

	return response.GetLinkToken(), nil
}

// ExchangePublicToken exchanges a public token for an access token
func (c *plaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (string, error) {
	request := plaidSDK.NewItemPublicTokenExchangeRequest(publicToken)

	response, _, err := c.sdkClient.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(*request).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to exchange public token: %v", err)
	}

	return response.GetAccessToken(), nil
}

// GetAccounts retrieves all accounts associated with an access token
func (c *plaidClient) GetAccounts(ctx context.Context, accessToken string) ([]BankAccount, error) {
	request := plaidSDK.NewAccountsGetRequest(accessToken)

	response, _, err := c.sdkClient.PlaidApi.AccountsGet(ctx).AccountsGetRequest(*request).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %v", err)
	}

	accounts := make([]BankAccount, len(response.GetAccounts()))
	for i, account := range response.GetAccounts() {
		accounts[i] = BankAccount{
			AccountID: account.GetAccountId(),
			Name:      account.GetName(),
			Type:      string(account.GetType()),
			Subtype:   string(account.GetSubtype()),
			Mask:      account.GetMask(),
		}
	}

	return accounts, nil
}

// GetBalance retrieves the current balance for a specific account
func (c *plaidClient) GetBalance(ctx context.Context, accessToken string, accountId string) (float64, error) {
	request := plaidSDK.NewAccountsBalanceGetRequest(accessToken)

	response, _, err := c.sdkClient.PlaidApi.AccountsBalanceGet(ctx).AccountsBalanceGetRequest(*request).Execute()
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %v", err)
	}

	for _, account := range response.GetAccounts() {
		if account.GetAccountId() == accountId {
			if account.Balances.Current.IsSet() {
				return *account.Balances.Current.Get(), nil
			}
			return 0, nil
		}
	}

	return 0, fmt.Errorf("account not found")
}
