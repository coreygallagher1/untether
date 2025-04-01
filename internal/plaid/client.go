package plaid

import (
	"context"
	"fmt"

	"github.com/plaid/plaid-go/v14/plaid"
)

// PlaidClient defines the interface for interacting with Plaid
type PlaidClient interface {
	GetAccountDetails(ctx context.Context, accessToken, accountID string) (*plaid.AccountBase, error)
	CreateLinkToken(ctx context.Context, clientUserID, clientName string) (string, error)
	ExchangePublicToken(ctx context.Context, publicToken string) (string, error)
}

// Client wraps the Plaid client with our custom functionality
type Client struct {
	plaidClient *plaid.APIClient
	environment plaid.Environment
}

// NewClient creates a new Plaid client with the provided credentials
func NewClient(clientID, secret, environment string) (PlaidClient, error) {
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)

	// Set the environment
	var env plaid.Environment
	switch environment {
	case "sandbox":
		env = plaid.Sandbox
	case "development":
		env = plaid.Development
	case "production":
		env = plaid.Production
	default:
		return nil, fmt.Errorf("invalid environment: %s", environment)
	}
	configuration.UseEnvironment(env)

	client := plaid.NewAPIClient(configuration)
	return &Client{
		plaidClient: client,
		environment: env,
	}, nil
}

// GetAccountDetails retrieves account details from Plaid
func (c *Client) GetAccountDetails(ctx context.Context, accessToken, accountID string) (*plaid.AccountBase, error) {
	request := plaid.NewAccountsGetRequest(accessToken)
	response, _, err := c.plaidClient.PlaidApi.AccountsGet(ctx).AccountsGetRequest(*request).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get account details: %v", err)
	}

	// Find the specific account
	for _, account := range response.GetAccounts() {
		if account.GetAccountId() == accountID {
			return &account, nil
		}
	}

	return nil, fmt.Errorf("account not found: %s", accountID)
}

// CreateLinkToken creates a link token for initializing Plaid Link
func (c *Client) CreateLinkToken(ctx context.Context, clientUserID, clientName string) (string, error) {
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: clientUserID,
	}

	request := plaid.NewLinkTokenCreateRequest(
		clientName,
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		user,
	)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH})

	response, _, err := c.plaidClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to create link token: %v", err)
	}

	return response.GetLinkToken(), nil
}

// ExchangePublicToken exchanges a public token for an access token
func (c *Client) ExchangePublicToken(ctx context.Context, publicToken string) (string, error) {
	request := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	response, _, err := c.plaidClient.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(*request).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to exchange public token: %v", err)
	}

	return response.GetAccessToken(), nil
}
