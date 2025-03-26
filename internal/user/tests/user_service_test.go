package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"untether/internal/proto"
	"untether/internal/service"
)

func TestUserService_CreateUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	db := setupTestDB(t)
	userService := service.NewUserService(db)

	// Test cases
	tests := []struct {
		name          string
		req           *proto.CreateUserRequest
		expectError   bool
		errorContains string
	}{
		{
			name: "successful user creation",
			req: &proto.CreateUserRequest{
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
			},
			expectError: false,
		},
		{
			name: "duplicate email",
			req: &proto.CreateUserRequest{
				Email:     "test@example.com",
				FirstName: "Another",
				LastName:  "User",
			},
			expectError:   true,
			errorContains: "already exists",
		},
		{
			name: "invalid email",
			req: &proto.CreateUserRequest{
				Email:     "invalid-email",
				FirstName: "Test",
				LastName:  "User",
			},
			expectError:   true,
			errorContains: "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create user
			createdUser, err := userService.CreateUser(ctx, tt.req)

			// Assertions
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, createdUser.Id)
			assert.Equal(t, tt.req.Email, createdUser.Email)
			assert.Equal(t, tt.req.FirstName, createdUser.FirstName)
			assert.Equal(t, tt.req.LastName, createdUser.LastName)
			assert.NotNil(t, createdUser.CreatedAt)
			assert.NotNil(t, createdUser.UpdatedAt)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	db := setupTestDB(t)
	userService := service.NewUserService(db)

	// Create a test user
	testUser, err := userService.CreateUser(ctx, &proto.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	})
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name          string
		req           *proto.GetUserRequest
		expectError   bool
		errorContains string
	}{
		{
			name: "successful user retrieval",
			req: &proto.GetUserRequest{
				Id: testUser.Id,
			},
			expectError: false,
		},
		{
			name: "non-existent user",
			req: &proto.GetUserRequest{
				Id: uuid.New().String(),
			},
			expectError:   true,
			errorContains: "not found",
		},
		{
			name: "invalid UUID",
			req: &proto.GetUserRequest{
				Id: "invalid-uuid",
			},
			expectError:   true,
			errorContains: "invalid UUID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get user
			foundUser, err := userService.GetUser(ctx, tt.req)

			// Assertions
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testUser.Id, foundUser.Id)
			assert.Equal(t, testUser.Email, foundUser.Email)
			assert.Equal(t, testUser.FirstName, foundUser.FirstName)
			assert.Equal(t, testUser.LastName, foundUser.LastName)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	db := setupTestDB(t)
	userService := service.NewUserService(db)

	// Create a test user
	testUser, err := userService.CreateUser(ctx, &proto.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	})
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name          string
		req           *proto.UpdateUserRequest
		expectError   bool
		errorContains string
	}{
		{
			name: "successful user update",
			req: &proto.UpdateUserRequest{
				Id:        testUser.Id,
				FirstName: "Updated",
				LastName:  "Name",
			},
			expectError: false,
		},
		{
			name: "non-existent user",
			req: &proto.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "Updated",
				LastName:  "Name",
			},
			expectError:   true,
			errorContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update user
			updatedUser, err := userService.UpdateUser(ctx, tt.req)

			// Assertions
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testUser.Id, updatedUser.Id)
			assert.Equal(t, testUser.Email, updatedUser.Email)
			assert.Equal(t, tt.req.FirstName, updatedUser.FirstName)
			assert.Equal(t, tt.req.LastName, updatedUser.LastName)
			assert.True(t, updatedUser.UpdatedAt.AsTime().After(testUser.UpdatedAt.AsTime()))
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	db := setupTestDB(t)
	userService := service.NewUserService(db)

	// Create a test user
	testUser, err := userService.CreateUser(ctx, &proto.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	})
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name          string
		req           *proto.DeleteUserRequest
		expectError   bool
		errorContains string
	}{
		{
			name: "successful user deletion",
			req: &proto.DeleteUserRequest{
				Id: testUser.Id,
			},
			expectError: false,
		},
		{
			name: "non-existent user",
			req: &proto.DeleteUserRequest{
				Id: uuid.New().String(),
			},
			expectError:   true,
			errorContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Delete user
			_, err := userService.DeleteUser(ctx, tt.req)

			// Assertions
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			require.NoError(t, err)

			// Verify user is deleted
			_, err = userService.GetUser(ctx, &proto.GetUserRequest{Id: tt.req.Id})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not found")
		})
	}
} 