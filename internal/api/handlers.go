package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alvinatthariq/billing-engine/internal/loan"
	"github.com/gorilla/mux"
)

type LoanHandler struct {
	loans map[string]*loan.Loan
}

func NewLoanHandler() *LoanHandler {
	return &LoanHandler{
		loans: make(map[string]*loan.Loan),
	}
}

type CreateLoanRequest struct {
	ID           string  `json:"id"`
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interest_rate"`
}

type PaymentRequest struct {
	Amount      float64 `json:"amount"`
	PaymentDate string  `json:"payment_date"`
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var req CreateLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, exists := h.loans[req.ID]; exists {
		http.Error(w, "Loan ID already exists", http.StatusConflict)
		return
	}

	newLoan, err := loan.NewLoan(req.ID, req.Amount, req.InterestRate, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.loans[req.ID] = newLoan

	response := map[string]interface{}{
		"id":              newLoan.ID,
		"amount":          newLoan.Amount,
		"interest_rate":   newLoan.InterestRate,
		"weekly_payment":  newLoan.WeeklyPayment,
		"total_repayment": newLoan.TotalRepayment,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *LoanHandler) GetOutstanding(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loanID := vars["id"]

	l, exists := h.loans[loanID]
	if !exists {
		http.Error(w, "Loan not found", http.StatusNotFound)
		return
	}

	response := map[string]float64{
		"outstanding": l.GetOutstanding(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *LoanHandler) IsDelinquent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loanID := vars["id"]

	l, exists := h.loans[loanID]
	if !exists {
		http.Error(w, "Loan not found", http.StatusNotFound)
		return
	}

	currentTime := time.Now()
	response := map[string]bool{
		"isDelinquent": l.IsDelinquent(currentTime),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *LoanHandler) MakePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loanID := vars["id"]

	l, exists := h.loans[loanID]
	if !exists {
		http.Error(w, "Loan not found", http.StatusNotFound)
		return
	}

	var paymentReq PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	paymentDate, err := time.Parse("2006-01-02", paymentReq.PaymentDate)
	if err != nil {
		http.Error(w, "Invalid payment date format", http.StatusBadRequest)
		return
	}

	if err := l.MakePayment(paymentReq.Amount, paymentDate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status": "Payment processed successfully",
	}
	json.NewEncoder(w).Encode(response)
}
