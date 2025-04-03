package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
	switch r.Method {
	case http.MethodPost:
		h.handleCreateUser(w, r)
	case http.MethodGet:
		h.handleGetUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	if path == "" || path == "users" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userID := path

	user, err := h.userService.GetUser(r.Context(), &pb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		log.Printf("Error getting user: %v", err)
		switch status.Code(err) {
		case codes.NotFound:
			http.Error(w, "User not found", http.StatusNotFound)
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
