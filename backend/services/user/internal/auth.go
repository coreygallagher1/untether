package internal

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "untether/services/user/proto"
)

const (
	// Argon2 parameters
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
)

// generateSalt generates a random salt for password hashing
func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// hashPassword hashes a password using Argon2id
func hashPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
}

// verifyPassword verifies a password against a hash
func verifyPassword(password string, salt, hash []byte) bool {
	newHash := hashPassword(password, salt)
	return subtle.ConstantTimeCompare(hash, newHash) == 1
}

func (s *UserService) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.AuthResponse, error) {
	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "all fields are required")
	}

	if !isValidEmail(req.Email) {
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}

	// Check for existing user
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing user: %v", err)
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	}

	// Generate salt and hash password
	salt, err := generateSalt()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate salt: %v", err)
	}
	hash := hashPassword(req.Password, salt)

	// Create user
	now := time.Now()
	user := &pb.User{
		Id:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}

	// Insert user into database
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO users (id, email, first_name, last_name, password_hash, password_salt, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.Id, user.Email, user.FirstName, user.LastName,
		base64.StdEncoding.EncodeToString(hash),
		base64.StdEncoding.EncodeToString(salt),
		user.CreatedAt.AsTime(), user.UpdatedAt.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.AuthResponse{
		UserId: user.Id,
		Token:  token,
		User:   user,
	}, nil
}

func (s *UserService) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.AuthResponse, error) {
	// Validate required fields
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Get user from database
	var user pb.User
	var passwordHash, passwordSalt string
	var createdAt, updatedAt time.Time

	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, first_name, last_name, password_hash, password_salt, created_at, updated_at 
		 FROM users WHERE email = $1`,
		req.Email,
	).Scan(
		&user.Id, &user.Email, &user.FirstName, &user.LastName,
		&passwordHash, &passwordSalt,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "invalid email or password")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// Verify password
	hash, err := base64.StdEncoding.DecodeString(passwordHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decode hash: %v", err)
	}
	salt, err := base64.StdEncoding.DecodeString(passwordSalt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decode salt: %v", err)
	}

	if !verifyPassword(req.Password, salt, hash) {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// Set timestamps
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	// Generate JWT token
	token, err := s.generateToken(user.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.AuthResponse{
		UserId: user.Id,
		Token:  token,
		User:   &user,
	}, nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	// TODO: Implement password reset logic with email verification
	return nil, status.Error(codes.Unimplemented, "password reset not implemented")
}

func (s *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	// TODO: Implement password change logic
	return nil, status.Error(codes.Unimplemented, "password change not implemented")
}
