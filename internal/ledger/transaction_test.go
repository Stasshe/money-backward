package ledger

import (
	"testing"
)

func TestNewTransaction(t *testing.T) {
	txn := NewTransaction("checking", "groceries", 50.0, "weekly shopping")
	if txn.AccountID != "checking" {
		t.Errorf("expected account_id 'checking', got %s", txn.AccountID)
	}
	if txn.Category != "groceries" {
		t.Errorf("expected category 'groceries', got %s", txn.Category)
	}
}

func TestTransactionIsExpense(t *testing.T) {
	tests := []struct {
		amount   float64
		expected bool
	}{
		{-50.0, true},
		{-100.0, true},
		{50.0, false},
		{100.0, false},
	}

	for _, tt := range tests {
		txn := NewTransaction("acc", "cat", tt.amount, "test")
		if txn.IsExpense() != tt.expected {
			t.Errorf("IsExpense(%v) = %v, expected %v", tt.amount, txn.IsExpense(), tt.expected)
		}
	}
}

func TestTransactionValidation(t *testing.T) {
	tests := []struct {
		name     string
		txn      *Transaction
		hasError bool
	}{
		{
			"valid transaction",
			&Transaction{ID: "1", AccountID: "checking", Category: "food", Amount: 50},
			false,
		},
		{
			"missing account_id",
			&Transaction{ID: "1", AccountID: "", Category: "food", Amount: 50},
			true,
		},
		{
			"missing category",
			&Transaction{ID: "1", AccountID: "checking", Category: "", Amount: 50},
			true,
		},
		{
			"zero amount",
			&Transaction{ID: "1", AccountID: "checking", Category: "food", Amount: 0},
			true,
		},
	}

	for _, tt := range tests {
		err := tt.txn.Validate()
		if (err != nil) != tt.hasError {
			t.Errorf("%s: Validate() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestTransactionAbsAmount(t *testing.T) {
	tests := []struct {
		amount   float64
		expected float64
	}{
		{50.0, 50.0},
		{-50.0, 50.0},
		{0.01, 0.01},
		{-0.01, 0.01},
	}

	for _, tt := range tests {
		txn := NewTransaction("acc", "cat", tt.amount, "test")
		if txn.AbsAmount() != tt.expected {
			t.Errorf("AbsAmount(%v) = %v, expected %v", tt.amount, txn.AbsAmount(), tt.expected)
		}
	}
}

func TestTransactionTags(t *testing.T) {
	txn := NewTransaction("acc", "cat", 50, "test")
	txn.AddTag("urgent")
	txn.AddTag("work")

	if len(txn.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(txn.Tags))
	}

	// Adding duplicate tag should not create duplicate
	txn.AddTag("urgent")
	if len(txn.Tags) != 2 {
		t.Errorf("expected 2 tags after adding duplicate, got %d", len(txn.Tags))
	}
}
