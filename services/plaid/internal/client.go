package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cgallagher/Untether/pkg/plaid"
	plaidSDK "github.com/plaid/plaid-go/v14/plaid"
)

type plaidClient struct {
	clientID     string
	clientSecret string
	environment  string
	httpClient   *http.Client
	sdkClient    *plaidSDK.APIClient
}

func NewPlaidClient(clientID, clientSecret, environment string) plaid.PlaidClient {
	configuration := plaidSDK.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	configuration.AddDefaultHeader("PLAID-SECRET", clientSecret)

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

	client := &plaidClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		environment:  environment,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		sdkClient: plaidSDK.NewAPIClient(configuration),
	}
	return client
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
	request.SetProducts([]plaidSDK.Products{plaidSDK.PRODUCTS_AUTH})

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
func (c *plaidClient) GetAccounts(ctx context.Context, accessToken string) ([]plaid.BankAccount, error) {
	request := plaidSDK.NewAccountsGetRequest(accessToken)
	response, _, err := c.sdkClient.PlaidApi.AccountsGet(ctx).AccountsGetRequest(*request).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %v", err)
	}

	var accounts []plaid.BankAccount
	for _, acc := range response.GetAccounts() {
		balance := 0.0
		if acc.Balances.Current.IsSet() {
			balance = *acc.Balances.Current.Get()
		}

		currency := "USD"
		if acc.Balances.IsoCurrencyCode.IsSet() {
			currency = *acc.Balances.IsoCurrencyCode.Get()
		}

		accounts = append(accounts, plaid.BankAccount{
			ID:          acc.GetAccountId(),
			Name:        acc.GetName(),
			Type:        string(acc.GetType()),
			Subtype:     string(acc.GetSubtype()),
			Balance:     balance,
			Currency:    currency,
			LastUpdated: time.Now(),
		})
	}

	return accounts, nil
}

// GetBalance retrieves the current balance for a specific account
func (c *plaidClient) GetBalance(ctx context.Context, accessToken string, accountId string) (float64, error) {
	request := plaidSDK.NewAccountsGetRequest(accessToken)
	response, _, err := c.sdkClient.PlaidApi.AccountsGet(ctx).AccountsGetRequest(*request).Execute()
	if err != nil {
		return 0, fmt.Errorf("failed to get account balance: %v", err)
	}

	for _, acc := range response.GetAccounts() {
		if acc.GetAccountId() == accountId {
			if acc.Balances.Current.IsSet() {
				return *acc.Balances.Current.Get(), nil
			}
			return 0, nil
		}
	}

	return 0, fmt.Errorf("account not found: %s", accountId)
}
