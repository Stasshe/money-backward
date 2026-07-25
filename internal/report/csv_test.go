package report

import (
	"strings"
	"testing"

	"money-backword/internal/ledger"
	"money-backword/internal/storage"
)

type MockStore struct {
	transactions   []*ledger.Transaction
	accounts       []*ledger.Account
	budgets        []*ledger.Budget
	monthlyReports map[string]*storage.MonthlyReport
}

func (m *MockStore) GetTransactions(accountID string, limit int) ([]*ledger.Transaction, error) {
	return m.transactions, nil
}

func (m *MockStore) GetAllAccounts() ([]*ledger.Account, error) {
	return m.accounts, nil
}

func (m *MockStore) GetAllBudgets() ([]*ledger.Budget, error) {
	return m.budgets, nil
}

func (m *MockStore) GenerateMonthlyReport(month string) (*storage.MonthlyReport, error) {
	if report, ok := m.monthlyReports[month]; ok {
		return report, nil
	}
	return &storage.MonthlyReport{
		Month:             month,
		CategoryBreakdown: make(map[string]float64),
		AccountBreakdown:  make(map[string]float64),
	}, nil
}

func (m *MockStore) AddTransaction(txn *ledger.Transaction) error          { return nil }
func (m *MockStore) GetTransaction(id string) (*ledger.Transaction, error) { return nil, nil }
func (m *MockStore) UpdateTransaction(txn *ledger.Transaction) error       { return nil }
func (m *MockStore) DeleteTransaction(id string) error                     { return nil }
func (m *MockStore) AddAccount(acc *ledger.Account) error                  { return nil }
func (m *MockStore) GetAccount(id string) (*ledger.Account, error)         { return nil, nil }
func (m *MockStore) UpdateAccount(acc *ledger.Account) error               { return nil }
func (m *MockStore) AddCategory(cat *ledger.Category) error                { return nil }
func (m *MockStore) GetCategory(name string) (*ledger.Category, error)     { return nil, nil }
func (m *MockStore) GetAllCategories() ([]ledger.Category, error)          { return nil, nil }
func (m *MockStore) SetBudget(budget *ledger.Budget) error                 { return nil }
func (m *MockStore) GetBudget(id string) (*ledger.Budget, error)           { return nil, nil }
func (m *MockStore) DeleteBudget(id string) error                          { return nil }
func (m *MockStore) Close() error                                          { return nil }

func TestExportMonthlyReport(t *testing.T) {
	mockStore := &MockStore{
		monthlyReports: map[string]*storage.MonthlyReport{
			"2024-01": {
				Month:            "2024-01",
				TotalIncome:      5000,
				TotalExpense:     3000,
				Net:              2000,
				TransactionCount: 25,
				CategoryBreakdown: map[string]float64{
					"groceries": 800,
					"utilities": 400,
					"rent":      1800,
				},
			},
		},
	}

	exporter := NewCSVExporter(mockStore)
	csv, err := exporter.ExportMonthlyReport("2024-01")
	if err != nil {
		t.Fatalf("ExportMonthlyReport failed: %v", err)
	}

	if !strings.Contains(csv, "Monthly Report") {
		t.Error("CSV missing 'Monthly Report' header")
	}
	if !strings.Contains(csv, "2024-01") {
		t.Error("CSV missing month in output")
	}
	if !strings.Contains(csv, "5000") {
		t.Error("CSV missing income amount")
	}
	if !strings.Contains(csv, "3000") {
		t.Error("CSV missing expense amount")
	}
	if !strings.Contains(csv, "groceries") {
		t.Error("CSV missing category breakdown")
	}
}

func TestExportAccountSummary(t *testing.T) {
	mockStore := &MockStore{
		accounts: []*ledger.Account{
			{
				ID:          "checking_01",
				Name:        "Main Checking",
				AccountType: "checking",
				Balance:     5000.00,
				Currency:    "USD",
				Active:      true,
			},
			{
				ID:          "savings_01",
				Name:        "Emergency Fund",
				AccountType: "savings",
				Balance:     15000.00,
				Currency:    "USD",
				Active:      true,
			},
		},
	}

	exporter := NewCSVExporter(mockStore)
	csv, err := exporter.ExportAccountSummary()
	if err != nil {
		t.Fatalf("ExportAccountSummary failed: %v", err)
	}

	if !strings.Contains(csv, "Account ID") {
		t.Error("CSV missing header")
	}
	if !strings.Contains(csv, "checking_01") {
		t.Error("CSV missing checking account")
	}
	if !strings.Contains(csv, "5000") {
		t.Error("CSV missing checking balance")
	}
}

func TestExportCategoryAnalysis(t *testing.T) {
	mockStore := &MockStore{
		monthlyReports: map[string]*storage.MonthlyReport{
			"2024-02": {
				Month:        "2024-02",
				TotalExpense: 2000,
				CategoryBreakdown: map[string]float64{
					"food":    600,
					"transit": 400,
					"misc":    1000,
				},
			},
		},
	}

	exporter := NewCSVExporter(mockStore)
	csv, err := exporter.ExportCategoryAnalysis("2024-02")
	if err != nil {
		t.Fatalf("ExportCategoryAnalysis failed: %v", err)
	}

	if !strings.Contains(csv, "Category") {
		t.Error("CSV missing category header")
	}
	if !strings.Contains(csv, "food") {
		t.Error("CSV missing food category")
	}
	if !strings.Contains(csv, "30.0%") {
		t.Error("CSV missing percentage calculation (food should be 30%)")
	}
}

func TestExportBudgetStatus(t *testing.T) {
	mockStore := &MockStore{
		budgets: []*ledger.Budget{
			{
				ID:           "b1",
				Category:     "groceries",
				MonthlyLimit: 500,
			},
			{
				ID:           "b2",
				Category:     "entertainment",
				MonthlyLimit: 200,
			},
		},
	}

	exporter := NewCSVExporter(mockStore)
	csv, err := exporter.ExportBudgetStatus()
	if err != nil {
		t.Fatalf("ExportBudgetStatus failed: %v", err)
	}

	if !strings.Contains(csv, "Category") {
		t.Error("CSV missing header")
	}
	if !strings.Contains(csv, "groceries") {
		t.Error("CSV missing groceries budget")
	}
	if !strings.Contains(csv, "500") {
		t.Error("CSV missing budget limit")
	}
}
