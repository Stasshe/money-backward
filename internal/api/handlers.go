package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"money-backword/internal/ledger"
	"money-backword/internal/storage"
)

type TransactionRequest struct {
	AccountID   string  `json:"account_id"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type AccountRequest struct {
	Name          string  `json:"name"`
	AccountType   string  `json:"account_type"`
	InitialAmount float64 `json:"initial_amount"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (h *Handler) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var req TransactionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, "invalid request format", http.StatusBadRequest)
		return
	}

	txn := ledger.NewTransaction(req.AccountID, req.Category, req.Amount, req.Description)
	if err := h.store.AddTransaction(txn); err != nil {
		writeError(w, fmt.Sprintf("failed to add transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    txn,
		Message: "transaction created",
	})
}

func (h *Handler) handleListTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountID := r.URL.Query().Get("account_id")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	txns, err := h.store.GetTransactions(accountID, limit)
	if err != nil {
		writeError(w, fmt.Sprintf("failed to fetch transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    txns,
	})
}

func (h *Handler) handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/transactions/"):]
	txn, err := h.store.GetTransaction(id)
	if err != nil {
		writeError(w, "transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    txn,
	})
}

func (h *Handler) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var req AccountRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, "invalid request format", http.StatusBadRequest)
		return
	}

	acc := ledger.NewAccount(generateID(), req.Name, req.AccountType, req.InitialAmount)
	if err := h.store.AddAccount(acc); err != nil {
		writeError(w, fmt.Sprintf("failed to create account: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    acc,
		Message: "account created",
	})
}

func (h *Handler) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accounts, err := h.store.GetAllAccounts()
	if err != nil {
		writeError(w, fmt.Sprintf("failed to fetch accounts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    accounts,
	})
}

func (h *Handler) handleGetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	month := r.URL.Query().Get("month")
	report, err := h.store.GenerateMonthlyReport(month)
	if err != nil {
		writeError(w, fmt.Sprintf("failed to generate report: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    report,
	})
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func generateID() string {
	// TODO: use proper UUID generation
	return fmt.Sprintf("acc_%d", len([]byte{}))
}

func (h *Handler) handleListCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories, err := h.store.GetAllCategories()
	if err != nil {
		writeError(w, fmt.Sprintf("failed to fetch categories: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    categories,
	})
}

func (h *Handler) handleListBudgets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	budgets, err := h.store.GetAllBudgets()
	if err != nil {
		writeError(w, fmt.Sprintf("failed to fetch budgets: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    budgets,
	})
}

func (h *Handler) handleNotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, "endpoint not found", http.StatusNotFound)
}
