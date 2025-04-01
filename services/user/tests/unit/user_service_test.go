package unit

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"untether/services/plaid/client"
	"untether/services/user/internal"
	pb "untether/services/user/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database
type MockDB struct {
	mock.Mock
}

// MockRow is a mock implementation of sql.Row
type MockRow struct {
	values []interface{}
	err    error
}

func createMockRow(values []interface{}, err error) *sql.Row {
	db, mock, _ := sqlmock.New()
	if err != nil {
		mock.ExpectQuery("").WillReturnError(err)
	} else {
		// Create columns based on the number of values
		columns := make([]string, len(values))
		for i := range columns {
			columns[i] = fmt.Sprintf("column%d", i)
		}
		rows := sqlmock.NewRows(columns)
		driverValues := make([]driver.Value, len(values))
		for i, v := range values {
			driverValues[i] = v
		}
		rows.AddRow(driverValues...)
		mock.ExpectQuery("").WillReturnRows(rows)
	}
	return db.QueryRow("")
}

func createMockRows(columns []string, values [][]interface{}) (*sql.Rows, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	rows := sqlmock.NewRows(columns)
	for _, rowValues := range values {
		driverValues := make([]driver.Value, len(rowValues))
		for i, v := range rowValues {
			driverValues[i] = v
		}
		rows.AddRow(driverValues...)
	}
	mock.ExpectQuery("").WillReturnRows(rows)
	return db.Query("")
}

func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	callArgs := m.Called(ctx, query, args)
	mockRow := callArgs.Get(0).(*MockRow)
	return createMockRow(mockRow.values, mockRow.err)
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(sql.Result), callArgs.Error(1)
}

func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(*sql.Rows), callArgs.Error(1)
}

// MockSQLResult implements sql.Result
type MockSQLResult struct {
	mock.Mock
}

func (m *MockSQLResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSQLResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// MockPlaidClient is a mock implementation of the PlaidClient interface
type MockPlaidClient struct {
	mock.Mock
}

func (m *MockPlaidClient) CreateLinkToken(ctx context.Context, userID string) (string, error) {
	callArgs := m.Called(ctx, userID)
	return callArgs.String(0), callArgs.Error(1)
}

func (m *MockPlaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (string, error) {
	callArgs := m.Called(ctx, publicToken)
	return callArgs.String(0), callArgs.Error(1)
}

func (m *MockPlaidClient) GetAccounts(ctx context.Context, accessToken string) ([]client.BankAccount, error) {
	callArgs := m.Called(ctx, accessToken)
	return callArgs.Get(0).([]client.BankAccount), callArgs.Error(1)
}

func (m *MockPlaidClient) GetBalance(ctx context.Context, accessToken string, accountID string) (float64, error) {
	callArgs := m.Called(ctx, accessToken, accountID)
	return callArgs.Get(0).(float64), callArgs.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockDB := new(MockDB)
	mockPlaidClient := new(MockPlaidClient)
	service := internal.NewUserService(mockDB, mockPlaidClient)

	t.Run("success", func(t *testing.T) {
		req := &pb.CreateUserRequest{
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		}

		// Create a mock row that will return false for the EXISTS query
		mockRow := &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", []interface{}{req.Email}).
			Return(mockRow).Once()

		mockResult := new(MockSQLResult)
		mockResult.On("LastInsertId").Return(int64(1), nil)
		mockResult.On("RowsAffected").Return(int64(1), nil)

		mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockResult, nil).Once()

		resp, err := service.CreateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Id)
		mockDB.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		req := &pb.CreateUserRequest{
			Email:     "existing@example.com",
			FirstName: "Existing",
			LastName:  "User",
		}

		// Create a mock row that will return true for the EXISTS query
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", []interface{}{req.Email}).
			Return(mockRow).Once()

		resp, err := service.CreateUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user with this email already exists")
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		req := &pb.CreateUserRequest{
			Email:     "invalid-email",
			FirstName: "Test",
			LastName:  "User",
		}

		resp, err := service.CreateUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("missing required fields", func(t *testing.T) {
		req := &pb.CreateUserRequest{
			Email: "test@example.com",
		}

		resp, err := service.CreateUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "email, first_name, and last_name are required")
	})
}

func TestGetUser(t *testing.T) {
	mockDB := new(MockDB)
	mockPlaidClient := new(MockPlaidClient)
	service := internal.NewUserService(mockDB, mockPlaidClient)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New().String()
		now := time.Now()
		req := &pb.GetUserRequest{
			Id: userID,
		}

		// Create a mock row with user data
		mockRow := &MockRow{
			values: []interface{}{
				userID,
				"test@example.com",
				"Test",
				"User",
				now,
				now,
			},
		}

		mockDB.On("QueryRowContext", mock.Anything, mock.Anything, []interface{}{req.Id}).
			Return(mockRow).Once()

		// Create mock rows for bank accounts
		columns := []string{"id", "user_id", "plaid_account_id", "name", "type", "subtype", "mask", "created_at", "updated_at"}
		values := [][]interface{}{} // Empty slice for no bank accounts
		mockRows, err := createMockRows(columns, values)
		assert.NoError(t, err)

		mockDB.On("QueryContext", mock.Anything, mock.Anything, []interface{}{req.Id}).
			Return(mockRows, nil).Once()

		resp, err := service.GetUser(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, userID, resp.Id)
		assert.Equal(t, "test@example.com", resp.Email)
		assert.Equal(t, "Test", resp.FirstName)
		assert.Equal(t, "User", resp.LastName)
		assert.NotNil(t, resp.CreatedAt)
		assert.NotNil(t, resp.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uuid.New().String()
		req := &pb.GetUserRequest{
			Id: userID,
		}

		// Create a mock row that returns ErrNoRows
		mockRow := &MockRow{
			err: sql.ErrNoRows,
		}

		mockDB.On("QueryRowContext", mock.Anything, mock.Anything, []interface{}{req.Id}).
			Return(mockRow).Once()

		resp, err := service.GetUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user not found")
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid user id", func(t *testing.T) {
		req := &pb.GetUserRequest{
			Id: "invalid-uuid",
		}

		resp, err := service.GetUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid UUID")
	})
}

func TestLinkBankAccount(t *testing.T) {
	mockDB := new(MockDB)
	mockPlaidClient := new(MockPlaidClient)
	service := internal.NewUserService(mockDB, mockPlaidClient)

	t.Run("success", func(t *testing.T) {
		req := &pb.LinkBankAccountRequest{
			UserId:           "test-user-id",
			PlaidAccessToken: "access-token",
			PlaidAccountId:   "test-account-id",
		}

		mockPlaidClient.On("GetAccounts", mock.Anything, req.PlaidAccessToken).
			Return([]client.BankAccount{
				{
					AccountID: req.PlaidAccountId,
					Name:      "Test Account",
					Type:      "checking",
					Subtype:   "personal",
					Mask:      "1234",
				},
			}, nil)

		mockPlaidClient.On("GetBalance", mock.Anything, req.PlaidAccessToken, req.PlaidAccountId).
			Return(float64(1000.00), nil)

		// Create a mock row that will return true for the EXISTS query (user exists)
		mockUserRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.UserId}).
			Return(mockUserRow).Once()

		// Create a mock row that will return false for the EXISTS query (bank account doesn't exist)
		mockBankAccountRow := &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM bank_accounts WHERE user_id = $1 AND plaid_account_id = $2)", []interface{}{req.UserId, req.PlaidAccountId}).
			Return(mockBankAccountRow).Once()

		mockResult := new(MockSQLResult)
		mockResult.On("LastInsertId").Return(int64(1), nil)
		mockResult.On("RowsAffected").Return(int64(1), nil)

		mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(mockResult, nil).Once()

		resp, err := service.LinkBankAccount(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Id)
		mockDB.AssertExpectations(t)
		mockPlaidClient.AssertExpectations(t)
	})

	t.Run("account already linked", func(t *testing.T) {
		req := &pb.LinkBankAccountRequest{
			UserId:           "test-user-id",
			PlaidAccessToken: "access-token",
			PlaidAccountId:   "existing-account-id",
		}

		mockPlaidClient.On("GetAccounts", mock.Anything, req.PlaidAccessToken).
			Return([]client.BankAccount{
				{
					AccountID: req.PlaidAccountId,
					Name:      "Existing Account",
					Type:      "checking",
					Subtype:   "personal",
					Mask:      "1234",
				},
			}, nil)

		// Create a mock row that will return true for the EXISTS query (user exists)
		mockUserRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.UserId}).
			Return(mockUserRow).Once()

		// Create a mock row that will return true for the EXISTS query (bank account already exists)
		mockBankAccountRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM bank_accounts WHERE user_id = $1 AND plaid_account_id = $2)", []interface{}{req.UserId, req.PlaidAccountId}).
			Return(mockBankAccountRow).Once()

		resp, err := service.LinkBankAccount(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "bank account already linked")
		mockDB.AssertExpectations(t)
		mockPlaidClient.AssertExpectations(t)
	})
}
