package storage

import "money-backword/internal/ledger"

type Store interface {
	AddTransaction(txn *ledger.Transaction) error
	GetTransaction(id string) (*ledger.Transaction, error)
	GetTransactions(accountID string, limit int) ([]*ledger.Transaction, error)
	UpdateTransaction(txn *ledger.Transaction) error
	DeleteTransaction(id string) error

	AddAccount(acc *ledger.Account) error
	GetAccount(id string) (*ledger.Account, error)
	GetAllAccounts() ([]*ledger.Account, error)
	UpdateAccount(acc *ledger.Account) error

	AddCategory(cat *ledger.Category) error
	GetCategory(name string) (*ledger.Category, error)
	GetAllCategories() ([]ledger.Category, error)

	SetBudget(budget *ledger.Budget) error
	GetBudget(id string) (*ledger.Budget, error)
	GetAllBudgets() ([]*ledger.Budget, error)
	DeleteBudget(id string) error

	GenerateMonthlyReport(month string) (*MonthlyReport, error)
	Close() error
}

type MonthlyReport struct {
	Month              string             `json:"month"`
	TotalIncome        float64            `json:"total_income"`
	TotalExpense       float64            `json:"total_expense"`
	Net                float64            `json:"net"`
	CategoryBreakdown  map[string]float64 `json:"category_breakdown"`
	AccountBreakdown   map[string]float64 `json:"account_breakdown"`
	TransactionCount   int                `json:"transaction_count"`
	AverageTransaction float64            `json:"average_transaction"`
}
