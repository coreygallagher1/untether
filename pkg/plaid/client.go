package plaid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BankAccount represents a bank account from Plaid
type BankAccount struct {
	ID          string
	Name        string
	Type        string
	Subtype     string
	Balance     float64
	Currency    string
	LastUpdated time.Time
}

// PlaidClient defines the interface for Plaid API interactions
type PlaidClient interface {
	// CreateLinkToken creates a new Plaid Link token for account linking
	CreateLinkToken(ctx context.Context, userId string) (string, error)

	// ExchangePublicToken exchanges a public token for an access token
	ExchangePublicToken(ctx context.Context, publicToken string) (string, error)

	// GetAccounts retrieves all accounts associated with an access token
	GetAccounts(ctx context.Context, accessToken string) ([]BankAccount, error)

	// GetBalance retrieves the current balance for a specific account
	GetBalance(ctx context.Context, accessToken string, accountId string) (float64, error)
}

type plaidClient struct {
	clientID     string
	clientSecret string
	environment  string
	httpClient   *http.Client
}

func NewPlaidClient(clientID, clientSecret, environment string) PlaidClient {
	return &plaidClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		environment:  environment,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *plaidClient) baseURL() string {
	switch c.environment {
	case "sandbox":
		return "https://sandbox.plaid.com"
	case "development":
		return "https://development.plaid.com"
	case "production":
		return "https://production.plaid.com"
	default:
		return "https://sandbox.plaid.com"
	}
}

func (c *plaidClient) makeRequest(ctx context.Context, endpoint string, payload interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL()+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PLAID-CLIENT-ID", c.clientID)
	req.Header.Set("PLAID-SECRET", c.clientSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plaid API error: %s", string(body))
	}

	return body, nil
}

func (c *plaidClient) CreateLinkToken(ctx context.Context, userId string) (string, error) {
	payload := map[string]interface{}{
		"client_id":     c.clientID,
		"secret":        c.clientSecret,
		"user":          map[string]string{"client_user_id": userId},
		"client_name":   "Untether",
		"products":      []string{"auth", "transactions"},
		"country_codes": []string{"US"},
		"language":      "en",
	}

	body, err := c.makeRequest(ctx, "/link/token/create", payload)
	if err != nil {
		return "", err
	}

	var response struct {
		LinkToken string `json:"link_token"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response.LinkToken, nil
}

func (c *plaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (string, error) {
	payload := map[string]interface{}{
		"client_id":    c.clientID,
		"secret":       c.clientSecret,
		"public_token": publicToken,
	}

	body, err := c.makeRequest(ctx, "/item/public_token/exchange", payload)
	if err != nil {
		return "", err
	}

	var response struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response.AccessToken, nil
}

func (c *plaidClient) GetAccounts(ctx context.Context, accessToken string) ([]BankAccount, error) {
	payload := map[string]interface{}{
		"client_id":    c.clientID,
		"secret":       c.clientSecret,
		"access_token": accessToken,
	}

	body, err := c.makeRequest(ctx, "/accounts/get", payload)
	if err != nil {
		return nil, err
	}

	var response struct {
		Accounts []struct {
			AccountID string `json:"account_id"`
			Name      string `json:"name"`
			Type      string `json:"type"`
			Subtype   string `json:"subtype"`
		} `json:"accounts"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	var accounts []BankAccount
	for _, acc := range response.Accounts {
		accounts = append(accounts, BankAccount{
			ID:      acc.AccountID,
			Name:    acc.Name,
			Type:    acc.Type,
			Subtype: acc.Subtype,
		})
	}

	return accounts, nil
}

func (c *plaidClient) GetBalance(ctx context.Context, accessToken string, accountId string) (float64, error) {
	payload := map[string]interface{}{
		"client_id":    c.clientID,
		"secret":       c.clientSecret,
		"access_token": accessToken,
	}

	body, err := c.makeRequest(ctx, "/accounts/balance/get", payload)
	if err != nil {
		return 0, err
	}

	var response struct {
		Accounts []struct {
			AccountID string `json:"account_id"`
			Balances  struct {
				Current float64 `json:"current"`
			} `json:"balances"`
		} `json:"accounts"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	for _, acc := range response.Accounts {
		if acc.AccountID == accountId {
			return acc.Balances.Current, nil
		}
	}

	return 0, fmt.Errorf("account not found: %s", accountId)
}
