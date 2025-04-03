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
		req := &pb.GetUserRequest{
			Id: uuid.New().String(),
		}

		// Create a mock row that will return sql.ErrNoRows
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

	t.Run("database error", func(t *testing.T) {
		req := &pb.GetUserRequest{
			Id: uuid.New().String(),
		}

		// Create a mock row that will return a database error
		mockRow := &MockRow{
			err: sql.ErrConnDone,
		}

		mockDB.On("QueryRowContext", mock.Anything, mock.Anything, []interface{}{req.Id}).
			Return(mockRow).Once()

		resp, err := service.GetUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to get user")
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid UUID", func(t *testing.T) {
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
		userID := uuid.New().String()
		req := &pb.LinkBankAccountRequest{
			UserId:           userID,
			PlaidAccessToken: "access-token",
			PlaidAccountId:   "account-id",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{userID}).
			Return(mockRow).Once()

		// Mock bank account exists check
		mockRow = &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM bank_accounts WHERE user_id = $1 AND plaid_account_id = $2)", []interface{}{userID, req.PlaidAccountId}).
			Return(mockRow).Once()

		// Mock Plaid API calls
		mockPlaidClient.On("GetAccounts", mock.Anything, req.PlaidAccessToken).
			Return([]client.BankAccount{
				{
					AccountID: req.PlaidAccountId,
					Name:      "Test Account",
					Type:      "checking",
					Subtype:   "personal",
					Mask:      "1234",
				},
			}, nil).Once()

		mockPlaidClient.On("GetBalance", mock.Anything, req.PlaidAccessToken, req.PlaidAccountId).
			Return(1000.00, nil).Once()

		// Mock bank account insertion
		mockResult := new(MockSQLResult)
		mockResult.On("LastInsertId").Return(int64(1), nil)
		mockResult.On("RowsAffected").Return(int64(1), nil)

		mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockResult, nil).Once()

		resp, err := service.LinkBankAccount(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Id)
		mockDB.AssertExpectations(t)
		mockPlaidClient.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		req := &pb.LinkBankAccountRequest{
			UserId:           uuid.New().String(),
			PlaidAccessToken: "access-token",
			PlaidAccountId:   "account-id",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.UserId}).
			Return(mockRow).Once()

		resp, err := service.LinkBankAccount(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user not found")
		mockDB.AssertExpectations(t)
	})

	t.Run("bank account already exists", func(t *testing.T) {
		req := &pb.LinkBankAccountRequest{
			UserId:           uuid.New().String(),
			PlaidAccessToken: "access-token",
			PlaidAccountId:   "account-id",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.UserId}).
			Return(mockRow).Once()

		// Mock bank account exists check
		mockRow = &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM bank_accounts WHERE user_id = $1 AND plaid_account_id = $2)", []interface{}{req.UserId, req.PlaidAccountId}).
			Return(mockRow).Once()

		resp, err := service.LinkBankAccount(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "bank account already linked")
		mockDB.AssertExpectations(t)
	})

	t.Run("plaid API error", func(t *testing.T) {
		req := &pb.LinkBankAccountRequest{
			UserId:           uuid.New().String(),
			PlaidAccessToken: "access-token",
			PlaidAccountId:   "account-id",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.UserId}).
			Return(mockRow).Once()

		// Mock bank account exists check
		mockRow = &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM bank_accounts WHERE user_id = $1 AND plaid_account_id = $2)", []interface{}{req.UserId, req.PlaidAccountId}).
			Return(mockRow).Once()

		// Mock Plaid API error
		mockPlaidClient.On("GetAccounts", mock.Anything, req.PlaidAccessToken).
			Return([]client.BankAccount{}, fmt.Errorf("plaid API error")).Once()

		resp, err := service.LinkBankAccount(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "plaid API error")
		mockDB.AssertExpectations(t)
		mockPlaidClient.AssertExpectations(t)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		req := &pb.LinkBankAccountRequest{
			UserId:           "invalid-uuid",
			PlaidAccessToken: "access-token",
			PlaidAccountId:   "account-id",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.UserId}).
			Return(mockRow).Once()

		resp, err := service.LinkBankAccount(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user not found")
		mockDB.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockDB := new(MockDB)
	mockPlaidClient := new(MockPlaidClient)
	service := internal.NewUserService(mockDB, mockPlaidClient)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New().String()
		now := time.Now()
		req := &pb.UpdateUserRequest{
			Id:        userID,
			FirstName: "Updated",
			LastName:  "User",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.Id}).
			Return(mockRow).Once()

		// Mock update
		mockResult := new(MockSQLResult)
		mockResult.On("LastInsertId").Return(int64(1), nil)
		mockResult.On("RowsAffected").Return(int64(1), nil)

		mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockResult, nil).Once()

		// Mock get user after update
		mockRow = &MockRow{
			values: []interface{}{
				userID,
				"test@example.com",
				"Updated",
				"User",
				now,
				now,
			},
		}
		mockDB.On("QueryRowContext", mock.Anything, mock.Anything, []interface{}{req.Id}).
			Return(mockRow).Once()

		// Mock get bank accounts after update
		columns := []string{"id", "plaid_account_id", "name", "type", "balance", "currency", "is_active", "created_at", "updated_at"}
		values := [][]interface{}{} // Empty slice for no bank accounts
		mockRows, err := createMockRows(columns, values)
		assert.NoError(t, err)

		mockDB.On("QueryContext", mock.Anything, mock.Anything, []interface{}{req.Id}).
			Return(mockRows, nil).Once()

		resp, err := service.UpdateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, userID, resp.Id)
		assert.Equal(t, "Updated", resp.FirstName)
		assert.Equal(t, "User", resp.LastName)
		mockDB.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		req := &pb.UpdateUserRequest{
			Id:        uuid.New().String(),
			FirstName: "Updated",
			LastName:  "User",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.Id}).
			Return(mockRow).Once()

		resp, err := service.UpdateUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user not found")
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		req := &pb.UpdateUserRequest{
			Id:        "invalid-uuid",
			FirstName: "Updated",
			LastName:  "User",
		}

		resp, err := service.UpdateUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid UUID")
	})

	t.Run("database error", func(t *testing.T) {
		req := &pb.UpdateUserRequest{
			Id:        uuid.New().String(),
			FirstName: "Updated",
			LastName:  "User",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.Id}).
			Return(mockRow).Once()

		// Mock update error
		mockResult := new(MockSQLResult)
		mockResult.On("LastInsertId").Return(int64(0), sql.ErrConnDone)
		mockResult.On("RowsAffected").Return(int64(0), sql.ErrConnDone)

		mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockResult, sql.ErrConnDone).Once()

		resp, err := service.UpdateUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to update user")
		mockDB.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	mockDB := new(MockDB)
	mockPlaidClient := new(MockPlaidClient)
	service := internal.NewUserService(mockDB, mockPlaidClient)

	t.Run("success", func(t *testing.T) {
		req := &pb.DeleteUserRequest{
			Id: uuid.New().String(),
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.Id}).
			Return(mockRow).Once()

		// Mock delete
		mockResult := new(MockSQLResult)
		mockResult.On("LastInsertId").Return(int64(1), nil)
		mockResult.On("RowsAffected").Return(int64(1), nil)

		mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockResult, nil).Once()

		resp, err := service.DeleteUser(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)
		mockDB.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		req := &pb.DeleteUserRequest{
			Id: uuid.New().String(),
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.Id}).
			Return(mockRow).Once()

		resp, err := service.DeleteUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user not found")
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		req := &pb.DeleteUserRequest{
			Id: "invalid-uuid",
		}

		resp, err := service.DeleteUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid UUID")
	})

	t.Run("database error", func(t *testing.T) {
		req := &pb.DeleteUserRequest{
			Id: uuid.New().String(),
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.Id}).
			Return(mockRow).Once()

		// Mock delete error
		mockResult := new(MockSQLResult)
		mockResult.On("LastInsertId").Return(int64(0), sql.ErrConnDone)
		mockResult.On("RowsAffected").Return(int64(0), sql.ErrConnDone)

		mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockResult, sql.ErrConnDone).Once()

		resp, err := service.DeleteUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to delete user")
		mockDB.AssertExpectations(t)
	})
}

func TestGetUserWithBankAccounts(t *testing.T) {
	mockDB := new(MockDB)
	mockPlaidClient := new(MockPlaidClient)
	service := internal.NewUserService(mockDB, mockPlaidClient)

	t.Run("success with bank accounts", func(t *testing.T) {
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
		columns := []string{"id", "plaid_account_id", "name", "type", "balance", "currency", "is_active", "created_at", "updated_at"}
		values := [][]interface{}{
			{
				uuid.New().String(),
				"plaid-account-1",
				"Checking Account",
				"checking",
				1000.00,
				"USD",
				true,
				now,
				now,
			},
			{
				uuid.New().String(),
				"plaid-account-2",
				"Savings Account",
				"savings",
				5000.00,
				"USD",
				true,
				now,
				now,
			},
		}
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
		assert.Len(t, resp.BankAccounts, 2)
		assert.Equal(t, "Checking Account", resp.BankAccounts[0].Name)
		assert.Equal(t, "Savings Account", resp.BankAccounts[1].Name)
		mockDB.AssertExpectations(t)
	})

	t.Run("bank accounts query error", func(t *testing.T) {
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

		// Mock bank accounts query error
		columns := []string{"id", "plaid_account_id", "name", "type", "balance", "currency", "is_active", "created_at", "updated_at"}
		mockRows, err := createMockRows(columns, [][]interface{}{})
		assert.NoError(t, err)

		mockDB.On("QueryContext", mock.Anything, mock.Anything, []interface{}{req.Id}).
			Return(mockRows, sql.ErrConnDone).Once()

		resp, err := service.GetUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to get bank accounts")
		mockDB.AssertExpectations(t)
	})

	t.Run("bank account scan error", func(t *testing.T) {
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

		// Create mock rows with invalid data to cause scan error
		columns := []string{"id", "plaid_account_id", "name", "type", "balance", "currency", "is_active", "created_at", "updated_at"}
		values := [][]interface{}{
			{
				"invalid-uuid", // This should cause a scan error
				"plaid-account-1",
				"Checking Account",
				"checking",
				"invalid-balance", // This should cause a scan error
				"USD",
				true,
				"invalid-time", // This should cause a scan error
				now,
			},
		}
		mockRows, err := createMockRows(columns, values)
		assert.NoError(t, err)

		mockDB.On("QueryContext", mock.Anything, mock.Anything, []interface{}{req.Id}).
			Return(mockRows, nil).Once()

		resp, err := service.GetUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to scan bank account")
		mockDB.AssertExpectations(t)
	})
}

func TestLinkBankAccountAdditional(t *testing.T) {
	mockDB := new(MockDB)
	mockPlaidClient := new(MockPlaidClient)
	service := internal.NewUserService(mockDB, mockPlaidClient)

	t.Run("missing required fields", func(t *testing.T) {
		testCases := []struct {
			name string
			req  *pb.LinkBankAccountRequest
		}{
			{
				name: "missing user ID",
				req: &pb.LinkBankAccountRequest{
					PlaidAccessToken: "access-token",
					PlaidAccountId:   "account-id",
				},
			},
			{
				name: "missing access token",
				req: &pb.LinkBankAccountRequest{
					UserId:         uuid.New().String(),
					PlaidAccountId: "account-id",
				},
			},
			{
				name: "missing account ID",
				req: &pb.LinkBankAccountRequest{
					UserId:           uuid.New().String(),
					PlaidAccessToken: "access-token",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				resp, err := service.LinkBankAccount(context.Background(), tc.req)
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), "required")
			})
		}
	})

	t.Run("account not found in Plaid", func(t *testing.T) {
		req := &pb.LinkBankAccountRequest{
			UserId:           uuid.New().String(),
			PlaidAccessToken: "access-token",
			PlaidAccountId:   "non-existent-account",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.UserId}).
			Return(mockRow).Once()

		// Mock bank account exists check
		mockRow = &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM bank_accounts WHERE user_id = $1 AND plaid_account_id = $2)", []interface{}{req.UserId, req.PlaidAccountId}).
			Return(mockRow).Once()

		// Mock Plaid API call with different account
		mockPlaidClient.On("GetAccounts", mock.Anything, req.PlaidAccessToken).
			Return([]client.BankAccount{
				{
					AccountID: "different-account",
					Name:      "Different Account",
					Type:      "checking",
				},
			}, nil).Once()

		resp, err := service.LinkBankAccount(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "bank account not found in Plaid")
		mockDB.AssertExpectations(t)
		mockPlaidClient.AssertExpectations(t)
	})

	t.Run("Plaid balance error", func(t *testing.T) {
		req := &pb.LinkBankAccountRequest{
			UserId:           uuid.New().String(),
			PlaidAccessToken: "access-token",
			PlaidAccountId:   "account-id",
		}

		// Mock user exists check
		mockRow := &MockRow{
			values: []interface{}{true},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", []interface{}{req.UserId}).
			Return(mockRow).Once()

		// Mock bank account exists check
		mockRow = &MockRow{
			values: []interface{}{false},
		}
		mockDB.On("QueryRowContext", mock.Anything, "SELECT EXISTS(SELECT 1 FROM bank_accounts WHERE user_id = $1 AND plaid_account_id = $2)", []interface{}{req.UserId, req.PlaidAccountId}).
			Return(mockRow).Once()

		// Mock Plaid API calls
		mockPlaidClient.On("GetAccounts", mock.Anything, req.PlaidAccessToken).
			Return([]client.BankAccount{
				{
					AccountID: req.PlaidAccountId,
					Name:      "Test Account",
					Type:      "checking",
				},
			}, nil).Once()

		mockPlaidClient.On("GetBalance", mock.Anything, req.PlaidAccessToken, req.PlaidAccountId).
			Return(0.0, fmt.Errorf("failed to get balance from Plaid")).Once()

		resp, err := service.LinkBankAccount(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to get account balance from Plaid")
		mockDB.AssertExpectations(t)
		mockPlaidClient.AssertExpectations(t)
	})
}

func TestListBankAccounts(t *testing.T) {
	mockDB := new(MockDB)
	mockPlaidClient := new(MockPlaidClient)
	service := internal.NewUserService(mockDB, mockPlaidClient)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New().String()
		req := &pb.ListBankAccountsRequest{
			UserId: userID,
		}

		// Mock Plaid API calls
		mockPlaidClient.On("GetAccounts", mock.Anything, "access-sandbox-123").
			Return([]client.BankAccount{
				{
					AccountID: "plaid-account-1",
					Name:      "Checking Account",
					Type:      "checking",
				},
				{
					AccountID: "plaid-account-2",
					Name:      "Savings Account",
					Type:      "savings",
				},
			}, nil).Once()

		// Mock GetBalance calls for each account
		mockPlaidClient.On("GetBalance", mock.Anything, "access-sandbox-123", "plaid-account-1").
			Return(1000.00, nil).Once()
		mockPlaidClient.On("GetBalance", mock.Anything, "access-sandbox-123", "plaid-account-2").
			Return(5000.00, nil).Once()

		resp, err := service.ListBankAccounts(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Accounts, 2)
		assert.Equal(t, "Checking Account", resp.Accounts[0].Name)
		assert.Equal(t, "Savings Account", resp.Accounts[1].Name)
		assert.Equal(t, float64(1000.00), resp.Accounts[0].Balance)
		assert.Equal(t, float64(5000.00), resp.Accounts[1].Balance)
		mockDB.AssertExpectations(t)
		mockPlaidClient.AssertExpectations(t)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		req := &pb.ListBankAccountsRequest{
			UserId: "invalid-uuid",
		}

		resp, err := service.ListBankAccounts(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid UUID")
	})

	t.Run("Plaid GetAccounts error", func(t *testing.T) {
		req := &pb.ListBankAccountsRequest{
			UserId: uuid.New().String(),
		}

		// Mock Plaid API error
		mockPlaidClient.On("GetAccounts", mock.Anything, "access-sandbox-123").
			Return([]client.BankAccount{}, fmt.Errorf("Plaid API error")).Once()

		resp, err := service.ListBankAccounts(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to get accounts")
		mockDB.AssertExpectations(t)
		mockPlaidClient.AssertExpectations(t)
	})

	t.Run("Plaid GetBalance error", func(t *testing.T) {
		req := &pb.ListBankAccountsRequest{
			UserId: uuid.New().String(),
		}

		// Mock Plaid API calls
		mockPlaidClient.On("GetAccounts", mock.Anything, "access-sandbox-123").
			Return([]client.BankAccount{
				{
					AccountID: "plaid-account-1",
					Name:      "Checking Account",
					Type:      "checking",
				},
			}, nil).Once()

		// Mock GetBalance error
		mockPlaidClient.On("GetBalance", mock.Anything, "access-sandbox-123", "plaid-account-1").
			Return(0.0, fmt.Errorf("Plaid API error")).Once()

		resp, err := service.ListBankAccounts(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to get balance")
		mockDB.AssertExpectations(t)
		mockPlaidClient.AssertExpectations(t)
	})
}
