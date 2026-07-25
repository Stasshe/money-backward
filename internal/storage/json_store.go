package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"money-backword/internal/ledger"
)

type JSONStore struct {
	path         string
	transactions []*ledger.Transaction
	accounts     []*ledger.Account
	categories   []ledger.Category
	budgets      []*ledger.Budget
	mu           sync.RWMutex
}

type JSONData struct {
	Transactions []*ledger.Transaction `json:"transactions"`
	Accounts     []*ledger.Account     `json:"accounts"`
	Categories   []ledger.Category     `json:"categories"`
	Budgets      []*ledger.Budget      `json:"budgets"`
}

func NewJSONStore(path string) (*JSONStore, error) {
	store := &JSONStore{
		path:         path,
		transactions: make([]*ledger.Transaction, 0),
		accounts:     make([]*ledger.Account, 0),
		categories:   make([]ledger.Category, 0),
		budgets:      make([]*ledger.Budget, 0),
	}

	// Try to load existing data
	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load store: %w", err)
	}

	// Initialize with default categories if empty
	if len(store.categories) == 0 {
		store.categories = ledger.DefaultCategories
	}

	return store, nil
}

func (s *JSONStore) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	var jd JSONData
	if err := json.Unmarshal(data, &jd); err != nil {
		return err
	}

	s.transactions = jd.Transactions
	s.accounts = jd.Accounts
	s.categories = jd.Categories
	s.budgets = jd.Budgets

	return nil
}

func (s *JSONStore) save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := JSONData{
		Transactions: s.transactions,
		Accounts:     s.accounts,
		Categories:   s.categories,
		Budgets:      s.budgets,
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, b, 0644)
}

// Transaction methods
func (s *JSONStore) AddTransaction(txn *ledger.Transaction) error {
	if err := txn.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	s.transactions = append(s.transactions, txn)
	s.mu.Unlock()

	return s.save()
}

func (s *JSONStore) GetTransaction(id string) (*ledger.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, txn := range s.transactions {
		if txn.ID == id {
			return txn, nil
		}
	}
	return nil, fmt.Errorf("transaction not found: %s", id)
}

func (s *JSONStore) GetTransactions(accountID string, limit int) ([]*ledger.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*ledger.Transaction

	for _, txn := range s.transactions {
		if accountID == "" || txn.AccountID == accountID {
			results = append(results, txn)
		}
	}

	// Sort by timestamp descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.After(results[j].Timestamp)
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func (s *JSONStore) UpdateTransaction(txn *ledger.Transaction) error {
	if err := txn.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	for i, existing := range s.transactions {
		if existing.ID == txn.ID {
			s.transactions[i] = txn
			s.mu.Unlock()
			return s.save()
		}
	}
	s.mu.Unlock()

	return fmt.Errorf("transaction not found: %s", txn.ID)
}

func (s *JSONStore) DeleteTransaction(id string) error {
	s.mu.Lock()
	for i, txn := range s.transactions {
		if txn.ID == id {
			s.transactions = append(s.transactions[:i], s.transactions[i+1:]...)
			s.mu.Unlock()
			return s.save()
		}
	}
	s.mu.Unlock()

	return fmt.Errorf("transaction not found: %s", id)
}

// Account methods
func (s *JSONStore) AddAccount(acc *ledger.Account) error {
	if err := acc.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	s.accounts = append(s.accounts, acc)
	s.mu.Unlock()

	return s.save()
}

func (s *JSONStore) GetAccount(id string) (*ledger.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, acc := range s.accounts {
		if acc.ID == id {
			return acc, nil
		}
	}
	return nil, fmt.Errorf("account not found: %s", id)
}

func (s *JSONStore) GetAllAccounts() ([]*ledger.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return active accounts only
	var active []*ledger.Account
	for _, acc := range s.accounts {
		if acc.Active {
			active = append(active, acc)
		}
	}
	return active, nil
}

func (s *JSONStore) UpdateAccount(acc *ledger.Account) error {
	if err := acc.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	for i, existing := range s.accounts {
		if existing.ID == acc.ID {
			s.accounts[i] = acc
			s.mu.Unlock()
			return s.save()
		}
	}
	s.mu.Unlock()

	return fmt.Errorf("account not found: %s", acc.ID)
}

// Category methods
func (s *JSONStore) AddCategory(cat *ledger.Category) error {
	if err := cat.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	s.categories = append(s.categories, *cat)
	s.mu.Unlock()

	return s.save()
}

func (s *JSONStore) GetCategory(name string) (*ledger.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, cat := range s.categories {
		if cat.Name == name {
			return &s.categories[i], nil
		}
	}
	return nil, fmt.Errorf("category not found: %s", name)
}

func (s *JSONStore) GetAllCategories() ([]ledger.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.categories, nil
}

// Budget methods
func (s *JSONStore) SetBudget(budget *ledger.Budget) error {
	if err := budget.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	// Check if budget for this category already exists
	for i, existing := range s.budgets {
		if existing.Category == budget.Category {
			s.budgets[i] = budget
			s.mu.Unlock()
			return s.save()
		}
	}
	s.budgets = append(s.budgets, budget)
	s.mu.Unlock()

	return s.save()
}

func (s *JSONStore) GetBudget(id string) (*ledger.Budget, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, budget := range s.budgets {
		if budget.ID == id {
			return budget, nil
		}
	}
	return nil, fmt.Errorf("budget not found: %s", id)
}

func (s *JSONStore) GetAllBudgets() ([]*ledger.Budget, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.budgets, nil
}

func (s *JSONStore) DeleteBudget(id string) error {
	s.mu.Lock()
	for i, budget := range s.budgets {
		if budget.ID == id {
			s.budgets = append(s.budgets[:i], s.budgets[i+1:]...)
			s.mu.Unlock()
			return s.save()
		}
	}
	s.mu.Unlock()

	return fmt.Errorf("budget not found: %s", id)
}

// Reporting
func (s *JSONStore) GenerateMonthlyReport(month string) (*MonthlyReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	report := &MonthlyReport{
		Month:             month,
		CategoryBreakdown: make(map[string]float64),
		AccountBreakdown:  make(map[string]float64),
	}

	now := time.Now()
	var monthTime time.Time
	if month == "" {
		monthTime = now
	} else {
		// TODO: parse month string (YYYY-MM format)
		t, err := time.Parse("2006-01", month)
		if err != nil {
			return nil, fmt.Errorf("invalid month format: %s", month)
		}
		monthTime = t
	}

	targetMonth := monthTime.Month()
	targetYear := monthTime.Year()

	for _, txn := range s.transactions {
		if txn.Timestamp.Month() == targetMonth && txn.Timestamp.Year() == targetYear {
			report.TransactionCount++

			if txn.IsIncome() {
				report.TotalIncome += txn.Amount
			} else {
				report.TotalExpense += txn.AbsAmount()
			}

			report.CategoryBreakdown[txn.Category] += txn.AbsAmount()
			report.AccountBreakdown[txn.AccountID] += txn.Amount
		}
	}

	report.Net = report.TotalIncome - report.TotalExpense
	if report.TransactionCount > 0 {
		report.AverageTransaction = (report.TotalIncome - report.TotalExpense) / float64(report.TransactionCount)
	}

	return report, nil
}

func (s *JSONStore) Close() error {
	return nil
}
