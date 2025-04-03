package unit

import (
	"context"
	"database/sql"
	"testing"

	"untether/services/transaction/internal"
	pb "untether/services/transaction/proto"

	"github.com/stretchr/testify/assert"
)

// mockDB implements the DB interface for testing
type mockDB struct{}

func (m *mockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

func (m *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return &mockResult{}, nil
}

func (m *mockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

// mockResult implements sql.Result for testing
type mockResult struct{}

func (m *mockResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (m *mockResult) RowsAffected() (int64, error) {
	return 1, nil
}

func TestCalculateRoundup(t *testing.T) {
	calculator := internal.NewTransactionCalculator(&mockDB{})
	ctx := context.Background()

	tests := []struct {
		name            string
		amount          float64
		roundingRule    string
		customAmount    float64
		expectedRoundup float64
		expectError     bool
	}{
		{
			name:            "Round up to dollar",
			amount:          4.55,
			roundingRule:    "dollar",
			expectedRoundup: 0.45,
			expectError:     false,
		},
		{
			name:            "Custom rounding",
			amount:          4.55,
			roundingRule:    "custom",
			customAmount:    0.50,
			expectedRoundup: 0.45,
			expectError:     false,
		},
		{
			name:         "Invalid rounding rule",
			amount:       4.55,
			roundingRule: "invalid",
			expectError:  true,
		},
		{
			name:         "Invalid custom amount",
			amount:       4.55,
			roundingRule: "custom",
			customAmount: 0,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.CalculateRoundupRequest{
				Amount:               tt.amount,
				RoundingRule:         tt.roundingRule,
				CustomRoundingAmount: tt.customAmount,
			}

			resp, err := calculator.CalculateRoundup(ctx, req)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.amount, resp.OriginalAmount)
			assert.InDelta(t, tt.expectedRoundup, resp.RoundupAmount, 0.0001)
			assert.Equal(t, tt.roundingRule, resp.RoundingRuleUsed)
		})
	}
}
