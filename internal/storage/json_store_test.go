package storage

import (
	"os"
	"testing"

	"money-backword/internal/ledger"
)

func TestNewJSONStore(t *testing.T) {
	// Use temp file
	tmpfile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	store, err := NewJSONStore(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewJSONStore failed: %v", err)
	}

	if store == nil {
		t.Error("expected non-nil store")
	}
}

func TestAddAndGetTransaction(t *testing.T) {
	tmpfile, _ := os.CreateTemp("", "test-*.json")
	defer os.Remove(tmpfile.Name())

	store, _ := NewJSONStore(tmpfile.Name())

	txn := ledger.NewTransaction("checking", "groceries", -50.0, "weekly shopping")
	err := store.AddTransaction(txn)
	if err != nil {
		t.Fatalf("AddTransaction failed: %v", err)
	}

	retrieved, err := store.GetTransaction(txn.ID)
	if err != nil {
		t.Fatalf("GetTransaction failed: %v", err)
	}

	if retrieved.ID != txn.ID {
		t.Errorf("expected id %s, got %s", txn.ID, retrieved.ID)
	}
	if retrieved.Amount != -50.0 {
		t.Errorf("expected amount -50.0, got %f", retrieved.Amount)
	}
}

func TestGetTransactions(t *testing.T) {
	tmpfile, _ := os.CreateTemp("", "test-*.json")
	defer os.Remove(tmpfile.Name())

	store, _ := NewJSONStore(tmpfile.Name())

	// Add multiple transactions
	txn1 := ledger.NewTransaction("checking", "groceries", -50.0, "shopping")
	txn2 := ledger.NewTransaction("checking", "utilities", -80.0, "electric")
	txn3 := ledger.NewTransaction("savings", "salary", 2000.0, "monthly")

	store.AddTransaction(txn1)
	store.AddTransaction(txn2)
	store.AddTransaction(txn3)

	// Get transactions for checking account
	txns, err := store.GetTransactions("checking", 10)
	if err != nil {
		t.Fatalf("GetTransactions failed: %v", err)
	}

	if len(txns) != 2 {
		t.Errorf("expected 2 transactions for checking, got %d", len(txns))
	}
}

func TestAddAndGetAccount(t *testing.T) {
	tmpfile, _ := os.CreateTemp("", "test-*.json")
	defer os.Remove(tmpfile.Name())

	store, _ := NewJSONStore(tmpfile.Name())

	acc := ledger.NewAccount("checking_01", "Main Checking", "checking", 1000.0)
	err := store.AddAccount(acc)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	retrieved, err := store.GetAccount(acc.ID)
	if err != nil {
		t.Fatalf("GetAccount failed: %v", err)
	}

	if retrieved.ID != acc.ID {
		t.Errorf("expected id %s, got %s", acc.ID, retrieved.ID)
	}
}

func TestAddCategory(t *testing.T) {
	tmpfile, _ := os.CreateTemp("", "test-*.json")
	defer os.Remove(tmpfile.Name())

	store, _ := NewJSONStore(tmpfile.Name())

	cat := ledger.NewCategory("Groceries", "expense")
	err := store.AddCategory(cat)
	if err != nil {
		t.Fatalf("AddCategory failed: %v", err)
	}

	retrieved, err := store.GetCategory("Groceries")
	if err != nil {
		t.Fatalf("GetCategory failed: %v", err)
	}

	if retrieved.Name != "Groceries" {
		t.Errorf("expected name Groceries, got %s", retrieved.Name)
	}
}

func TestSetAndGetBudget(t *testing.T) {
	tmpfile, _ := os.CreateTemp("", "test-*.json")
	defer os.Remove(tmpfile.Name())

	store, _ := NewJSONStore(tmpfile.Name())

	budget := ledger.NewBudget("groceries", 500.0)
	err := store.SetBudget(budget)
	if err != nil {
		t.Fatalf("SetBudget failed: %v", err)
	}

	retrieved, err := store.GetBudget(budget.ID)
	if err != nil {
		t.Fatalf("GetBudget failed: %v", err)
	}

	if retrieved.MonthlyLimit != 500.0 {
		t.Errorf("expected limit 500.0, got %f", retrieved.MonthlyLimit)
	}
}

func TestDeleteTransaction(t *testing.T) {
	tmpfile, _ := os.CreateTemp("", "test-*.json")
	defer os.Remove(tmpfile.Name())

	store, _ := NewJSONStore(tmpfile.Name())

	txn := ledger.NewTransaction("checking", "groceries", -50.0, "shopping")
	store.AddTransaction(txn)

	err := store.DeleteTransaction(txn.ID)
	if err != nil {
		t.Fatalf("DeleteTransaction failed: %v", err)
	}

	_, err = store.GetTransaction(txn.ID)
	if err == nil {
		t.Error("expected error when retrieving deleted transaction")
	}
}

func TestGenerateMonthlyReport(t *testing.T) {
	tmpfile, _ := os.CreateTemp("", "test-*.json")
	defer os.Remove(tmpfile.Name())

	store, _ := NewJSONStore(tmpfile.Name())

	// Add transactions
	store.AddTransaction(ledger.NewTransaction("checking", "salary", 3000.0, "monthly"))
	store.AddTransaction(ledger.NewTransaction("checking", "groceries", -200.0, "shopping"))
	store.AddTransaction(ledger.NewTransaction("checking", "utilities", -100.0, "electric"))

	report, err := store.GenerateMonthlyReport("")
	if err != nil {
		t.Fatalf("GenerateMonthlyReport failed: %v", err)
	}

	if report.TotalIncome != 3000.0 {
		t.Errorf("expected income 3000.0, got %f", report.TotalIncome)
	}
	if report.TotalExpense != 300.0 {
		t.Errorf("expected expense 300.0, got %f", report.TotalExpense)
	}
	if report.TransactionCount != 3 {
		t.Errorf("expected 3 transactions, got %d", report.TransactionCount)
	}
}
