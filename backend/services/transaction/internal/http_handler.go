package internal

import (
	"encoding/json"
	"net/http"

	pb "untether/services/transaction/proto"
)

type HTTPHandler struct {
	calculator *TransactionCalculator
}

func NewHTTPHandler(calculator *TransactionCalculator) *HTTPHandler {
	return &HTTPHandler{
		calculator: calculator,
	}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCalculateRoundup(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandler) handleCalculateRoundup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount               float64 `json:"amount"`
		RoundingRule         string  `json:"rounding_rule"`
		CustomRoundingAmount float64 `json:"custom_rounding_amount,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create a gRPC request
	grpcReq := &pb.CalculateRoundupRequest{
		Amount:               req.Amount,
		RoundingRule:         req.RoundingRule,
		CustomRoundingAmount: req.CustomRoundingAmount,
	}

	// Call the calculator service
	resp, err := h.calculator.CalculateRoundup(r.Context(), grpcReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
