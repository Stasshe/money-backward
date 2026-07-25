package ledger

import (
	"testing"
	"time"
)

func TestNewRecurringTransaction(t *testing.T) {
	rt := NewRecurringTransaction("acc1", "utilities", -100.0, "monthly bill", "monthly")
	if rt.AccountID != "acc1" {
		t.Errorf("expected account_id 'acc1', got %s", rt.AccountID)
	}
	if rt.Frequency != "monthly" {
		t.Errorf("expected frequency 'monthly', got %s", rt.Frequency)
	}
	if !rt.Active {
		t.Error("expected recurring transaction to be active")
	}
}

func TestRecurringTransactionNextOccurrence(t *testing.T) {
	tests := []struct {
		name      string
		frequency string
		daysAdded int
	}{
		{"daily", "daily", 1},
		{"weekly", "weekly", 7},
		{"biweekly", "biweekly", 14},
		{"monthly", "monthly", 30},     // approximate
		{"quarterly", "quarterly", 90}, // approximate
		{"yearly", "yearly", 365},      // approximate
	}

	for _, tt := range tests {
		rt := NewRecurringTransaction("acc1", "test", 100, "test", tt.frequency)
		next, err := rt.NextOccurrence()
		if err != nil {
			t.Errorf("%s: NextOccurrence() error = %v", tt.name, err)
			continue
		}

		diff := next.Sub(rt.StartDate)
		expectedMin := time.Duration(tt.daysAdded-3) * 24 * time.Hour
		expectedMax := time.Duration(tt.daysAdded+3) * 24 * time.Hour

		if diff < expectedMin || diff > expectedMax {
			t.Errorf("%s: expected diff ~%d days, got %v", tt.name, tt.daysAdded, diff.Hours()/24)
		}
	}
}

func TestRecurringTransactionApply(t *testing.T) {
	rt := NewRecurringTransaction("acc1", "groceries", -50.0, "shopping", "weekly")
	txn := rt.Apply()

	if txn.AccountID != rt.AccountID {
		t.Errorf("expected account_id %s, got %s", rt.AccountID, txn.AccountID)
	}
	if txn.Amount != rt.Amount {
		t.Errorf("expected amount %f, got %f", rt.Amount, txn.Amount)
	}
	if rt.LastApplied == nil {
		t.Error("expected LastApplied to be set after Apply()")
	}
}

func TestRecurringTransactionValidate(t *testing.T) {
	tests := []struct {
		name     string
		rt       *RecurringTransaction
		hasError bool
	}{
		{
			"valid recurring transaction",
			NewRecurringTransaction("acc1", "utilities", -100, "bill", "monthly"),
			false,
		},
		{
			"missing account_id",
			&RecurringTransaction{ID: "r1", Category: "test", Amount: 100, Frequency: "monthly"},
			true,
		},
		{
			"missing category",
			&RecurringTransaction{ID: "r1", AccountID: "acc1", Amount: 100, Frequency: "monthly"},
			true,
		},
		{
			"zero amount",
			&RecurringTransaction{ID: "r1", AccountID: "acc1", Category: "test", Amount: 0, Frequency: "monthly"},
			true,
		},
		{
			"invalid frequency",
			&RecurringTransaction{ID: "r1", AccountID: "acc1", Category: "test", Amount: 100, Frequency: "never"},
			true,
		},
	}

	for _, tt := range tests {
		err := tt.rt.Validate()
		if (err != nil) != tt.hasError {
			t.Errorf("%s: Validate() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestRecurringTransactionDeactivate(t *testing.T) {
	rt := NewRecurringTransaction("acc1", "test", 100, "test", "monthly")
	if !rt.IsActive() {
		t.Error("expected active recurring transaction")
	}

	rt.Deactivate()
	if rt.IsActive() {
		t.Error("expected inactive recurring transaction after deactivate")
	}
}

func TestRecurringTransactionIsActive(t *testing.T) {
	now := time.Now()
	past := now.AddDate(0, 0, -10)
	future := now.AddDate(0, 0, 10)

	rt := NewRecurringTransaction("acc1", "test", 100, "test", "monthly")
	rt.EndDate = &past
	if rt.IsActive() {
		t.Error("expected inactive when end date is in past")
	}

	rt.EndDate = &future
	if !rt.IsActive() {
		t.Error("expected active when end date is in future")
	}
}
