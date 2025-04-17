package middleware

import (
	"context"
	"net/http"
	"strings"

	"untether/services/user/internal"
	ctx "untether/services/user/pkg/context"
)

type AuthMiddleware struct {
	userService *internal.UserService
}

func NewAuthMiddleware(userService *internal.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Skip authentication for public routes
		if isPublicRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// Validate token
		claims, err := m.userService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context
		newCtx := context.WithValue(r.Context(), ctx.UserIDKey, claims.UserID)
		r = r.WithContext(newCtx)

		next.ServeHTTP(w, r)
	})
}

// isPublicRoute returns true if the route does not require authentication
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/api/v1/auth/signup",
		"/api/v1/auth/signin",
		"/api/v1/auth/reset-password",
	}

	for _, route := range publicRoutes {
		if path == route {
			return true
		}
	}
	return false
}
