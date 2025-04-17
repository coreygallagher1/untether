package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	ctx "untether/services/user/pkg/context"
	pb "untether/services/user/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPHandler struct {
	userService *UserService
}

func NewHTTPHandler(userService *UserService) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
	}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := r.URL.Path

	if r.Method == http.MethodPost {
		switch path {
		case "/api/v1/auth/signup":
			h.handleSignUp(w, r)
			return
		case "/api/v1/auth/signin":
			h.handleSignIn(w, r)
			return
		case "/api/v1/auth/reset-password":
			h.handleResetPassword(w, r)
			return
		case "/api/v1/auth/change-password":
			h.handleChangePassword(w, r)
			return
		case "/api/v1/users":
			h.handleCreateUser(w, r)
			return
		case "/api/v1/users/preferences":
			h.handleCreateUserPreferences(w, r)
			return
		}
	} else if r.Method == http.MethodGet && strings.HasPrefix(path, "/api/v1/users/") {
		h.handleGetUser(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *HTTPHandler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.userService.SignUp(r.Context(), &pb.SignUpRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		log.Printf("Error in signup: %v", err)
		switch status.Code(err) {
		case codes.AlreadyExists:
			http.Error(w, "User already exists", http.StatusConflict)
		case codes.InvalidArgument:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *HTTPHandler) handleSignIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.userService.SignIn(r.Context(), &pb.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Error in signin: %v", err)
		switch status.Code(err) {
		case codes.NotFound, codes.Unauthenticated:
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *HTTPHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	if path == "" || path == "users" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	requestedUserID := path

	// Get authenticated user ID from context
	authenticatedUserID, ok := r.Context().Value(ctx.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Only allow users to access their own data
	if requestedUserID != authenticatedUserID {
		http.Error(w, "Forbidden: you can only access your own data", http.StatusForbidden)
		return
	}

	user, err := h.userService.GetUser(r.Context(), &pb.GetUserRequest{
		Id: requestedUserID,
	})
	if err != nil {
		log.Printf("Error getting user: %v", err)
		switch status.Code(err) {
		case codes.NotFound:
			http.Error(w, "User not found", http.StatusNotFound)
		case codes.InvalidArgument:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case codes.Internal:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *HTTPHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), &pb.CreateUserRequest{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		switch status.Code(err) {
		case codes.AlreadyExists:
			http.Error(w, "User already exists", http.StatusConflict)
		case codes.InvalidArgument:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case codes.Internal:
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *HTTPHandler) handleCreateUserPreferences(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserId   string  `json:"user_id"`
		Currency string  `json:"currency"`
		Timezone string  `json:"timezone"`
		Language string  `json:"language"`
		DarkMode bool    `json:"dark_mode"`
		Budget   float64 `json:"budget"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userPreferences, err := h.userService.CreateUserPreferences(r.Context(), &pb.CreateUserPreferencesRequest{
		UserId:   req.UserId,
		Currency: req.Currency,
		Timezone: req.Timezone,
		Language: req.Language,
		DarkMode: req.DarkMode,
		Budget:   req.Budget,
	})
	if err != nil {
		log.Printf("Error creating user preferences: %v", err)
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userPreferences); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *HTTPHandler) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.userService.ResetPassword(r.Context(), &pb.ResetPasswordRequest{
		Email: req.Email,
	})
	if err != nil {
		log.Printf("Error in reset password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *HTTPHandler) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserId      string `json:"userId"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.userService.ChangePassword(r.Context(), &pb.ChangePasswordRequest{
		UserId:      req.UserId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		log.Printf("Error in change password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
