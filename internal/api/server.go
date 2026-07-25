package api

import (
	"net/http"

	"money-backword/internal/storage"
)

type Handler struct {
	store  *storage.JSONStore
	router *http.ServeMux
}

func NewServer(addr string, store *storage.JSONStore) *http.Server {
	handler := &Handler{
		store:  store,
		router: http.NewServeMux(),
	}

	// Transaction endpoints
	handler.router.HandleFunc("/api/transactions", handler.handleListTransactions)
	handler.router.HandleFunc("/api/transactions/create", handler.handleCreateTransaction)
	handler.router.HandleFunc("/api/transaction/", handler.handleGetTransaction)

	// Account endpoints
	handler.router.HandleFunc("/api/accounts", handler.handleListAccounts)
	handler.router.HandleFunc("/api/accounts/create", handler.handleCreateAccount)

	// Category endpoints
	handler.router.HandleFunc("/api/categories", handler.handleListCategories)

	// Budget endpoints
	handler.router.HandleFunc("/api/budgets", handler.handleListBudgets)

	// Report endpoints
	handler.router.HandleFunc("/api/report", handler.handleGetReport)

	// Health check
	handler.router.HandleFunc("/health", handler.handleHealth)

	// Default 404
	handler.router.HandleFunc("/", handler.handleNotFound)

	return &http.Server{
		Addr:    addr,
		Handler: handler.router,
	}
}
