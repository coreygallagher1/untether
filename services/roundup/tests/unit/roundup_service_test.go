package unit

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"untether/services/roundup/internal"
	pb "untether/services/roundup/proto"
)

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

// MockDB is a mock implementation of the database
type MockDB struct {
	mock.Mock
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(sql.Result), callArgs.Error(1)
}

func TestRoundupTransaction(t *testing.T) {
	mockDB := new(MockDB)
	service := internal.NewRoundupService(mockDB)

	ctx := context.Background()
	userID := uuid.New().String()
	req := &pb.RoundupRequest{
		UserId:            userID,
		TransactionAmount: 10.75,
	}

	t.Run("success", func(t *testing.T) {
		mockResult := new(MockSQLResult)
		mockResult.On("LastInsertId").Return(int64(1), nil)
		mockResult.On("RowsAffected").Return(int64(1), nil)

		mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockResult, nil).Once()

		response, err := service.RoundupTransaction(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.RoundupId)
		assert.Equal(t, 0.25, response.AmountRounded) // 11.00 - 10.75 = 0.25
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid user id", func(t *testing.T) {
		req := &pb.RoundupRequest{
			UserId:            "invalid-uuid",
			TransactionAmount: 10.75,
		}

		_, err := service.RoundupTransaction(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UUID")
	})

	t.Run("database error", func(t *testing.T) {
		mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(&MockSQLResult{}, sql.ErrConnDone).Once()

		_, err := service.RoundupTransaction(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		mockDB.AssertExpectations(t)
	})
}
