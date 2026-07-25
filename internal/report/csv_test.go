package report

import (
	"strings"
	"testing"
)

type MockStore struct {
	transactions   []*Transaction
	accounts       []*Account
	budgets        []*Budget
	monthlyReports map[string]*MonthlyReport
}

type Transaction struct {
	ID          string
	AccountID   string
	Category    string
	Amount      float64
	Description string
	Timestamp   string
}

type Account struct {
	ID            string
	Name          string
	AccountType   string
	Balance       float64
	Currency      string
	Active        bool
	LastModified  string
	InitialAmount float64
}

type Budget struct {
	ID            string
	Category      string
	MonthlyLimit  float64
	CurrentSpend  float64
	CreatedAt     string
}

type MonthlyReport struct {
	Month             string
	TotalIncome       float64
	TotalExpense      float64
	Net               float64
	TransactionCount  int
	AverageTransaction float64
	CategoryBreakdown map[string]float64
	AccountBreakdown  map[string]float64
}

func (m *MockStore) GetTransactions(accountID string, limit int) (interface{}, error) {
	return m.transactions, nil
}

func (m *MockStore) GetAllAccounts() (interface{}, error) {
	return m.accounts, nil
}

func (m *MockStore) GetAllBudgets() (interface{}, error) {
	return m.budgets, nil
}

func (m *MockStore) GenerateMonthlyReport(month string) (*MonthlyReport, error) {
	if report, ok := m.monthlyReports[month]; ok {
		return report, nil
	}
	return &MonthlyReport{
		Month:             month,
		CategoryBreakdown: make(map[string]float64),
		AccountBreakdown:  make(map[string]float64),
	}, nil
}

func (m *MockStore) AddTransaction(txn interface{}) error { return nil }
func (m *MockStore) GetTransaction(id string) (interface{}, error) { return nil, nil }
func (m *MockStore) UpdateTransaction(txn interface{}) error { return nil }
func (m *MockStore) DeleteTransaction(id string) error { return nil }
func (m *MockStore) AddAccount(acc interface{}) error { return nil }
func (m *MockStore) GetAccount(id string) (interface{}, error) { return nil, nil }
func (m *MockStore) UpdateAccount(acc interface{}) error { return nil }
func (m *MockStore) AddCategory(cat interface{}) error { return nil }
func (m *MockStore) GetCategory(name string) (interface{}, error) { return nil, nil }
func (m *MockStore) GetAllCategories() (interface{}, error) { return nil, nil }
func (m *MockStore) SetBudget(budget interface{}) error { return nil }
func (m *MockStore) GetBudget(id string) (interface{}, error) { return nil, nil }
func (m *MockStore) DeleteBudget(id string) error { return nil }
func (m *MockStore) Close() error { return nil }

func TestExportMonthlyReport(t *testing.T) {
	mockStore := &MockStore{
		monthlyReports: map[string]*MonthlyReport{
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

	// Verify CSV contains expected headers and data
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
		accounts: []*Account{
			{
				ID:            "checking_01",
				Name:          "Main Checking",
				AccountType:   "checking",
				Balance:       5000.00,
				Currency:      "USD",
				Active:        true,
				InitialAmount: 1000.00,
			},
			{
				ID:            "savings_01",
				Name:          "Emergency Fund",
				AccountType:   "savings",
				Balance:       15000.00,
				Currency:      "USD",
				Active:        true,
				InitialAmount: 10000.00,
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
		monthlyReports: map[string]*MonthlyReport{
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
		budgets: []*Budget{
			{
				ID:           "b1",
				Category:     "groceries",
				MonthlyLimit: 500,
				CurrentSpend: 350,
			},
			{
				ID:           "b2",
				Category:     "entertainment",
				MonthlyLimit: 200,
				CurrentSpend: 180,
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
