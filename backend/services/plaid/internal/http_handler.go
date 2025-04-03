package internal

import (
	"encoding/json"
	"net/http"

	"untether/services/plaid/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPHandler struct {
	service *PlaidService
}

func NewHTTPHandler(service *PlaidService) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/link-token":
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleCreateLinkToken(w, r)
	case "/exchange-token":
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleExchangePublicToken(w, r)
	case "/accounts":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleGetAccounts(w, r)
	case "/balance":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleGetBalance(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *HTTPHandler) handleCreateLinkToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateLinkToken(r.Context(), &proto.CreateLinkTokenRequest{
		UserId: req.UserID,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *HTTPHandler) handleExchangePublicToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PublicToken string `json:"public_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.ExchangePublicToken(r.Context(), &proto.ExchangePublicTokenRequest{
		PublicToken: req.PublicToken,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *HTTPHandler) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	if accessToken == "" {
		http.Error(w, "access_token is required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetAccounts(r.Context(), &proto.GetAccountsRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *HTTPHandler) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	if accessToken == "" {
		http.Error(w, "access_token is required", http.StatusBadRequest)
		return
	}

	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetBalance(r.Context(), &proto.GetBalanceRequest{
		AccessToken: accessToken,
		AccountId:   accountID,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
