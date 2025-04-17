package internal

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// TokenExpiry is the duration for which a token is valid
	TokenExpiry = 24 * time.Hour
)

// Claims represents the claims in a JWT token
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// generateToken generates a new JWT token for a user
func (s *UserService) generateToken(userID string) (string, error) {
	// Create the claims
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// TODO: Use a proper secret key from configuration
	secretKey := []byte("your-256-bit-secret")

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *UserService) ValidateToken(tokenString string) (*Claims, error) {
	// TODO: Use a proper secret key from configuration
	secretKey := []byte("your-256-bit-secret")

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Get the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
