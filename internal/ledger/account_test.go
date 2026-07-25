package ledger

import (
	"testing"
)

func TestNewAccount(t *testing.T) {
	acc := NewAccount("checking_01", "Main Checking", "checking", 1000.0)
	if acc.ID != "checking_01" {
		t.Errorf("expected id 'checking_01', got %s", acc.ID)
	}
	if acc.Balance != 1000.0 {
		t.Errorf("expected balance 1000.0, got %f", acc.Balance)
	}
	if !acc.Active {
		t.Error("expected account to be active")
	}
}

func TestAccountCredit(t *testing.T) {
	acc := NewAccount("acc1", "Test", "checking", 100.0)
	err := acc.Credit(50.0)
	if err != nil {
		t.Errorf("Credit failed: %v", err)
	}
	if acc.Balance != 150.0 {
		t.Errorf("expected balance 150.0, got %f", acc.Balance)
	}
}

func TestAccountDebit(t *testing.T) {
	acc := NewAccount("acc1", "Test", "checking", 100.0)
	err := acc.Debit(30.0)
	if err != nil {
		t.Errorf("Debit failed: %v", err)
	}
	if acc.Balance != 70.0 {
		t.Errorf("expected balance 70.0, got %f", acc.Balance)
	}
}

func TestAccountCreditNegativeAmount(t *testing.T) {
	acc := NewAccount("acc1", "Test", "checking", 100.0)
	err := acc.Credit(-10.0)
	if err == nil {
		t.Error("expected error when crediting negative amount")
	}
}

func TestAccountDebitNegativeAmount(t *testing.T) {
	acc := NewAccount("acc1", "Test", "checking", 100.0)
	err := acc.Debit(-10.0)
	if err == nil {
		t.Error("expected error when debiting negative amount")
	}
}

func TestAccountDeactivate(t *testing.T) {
	acc := NewAccount("acc1", "Test", "checking", 100.0)
	if !acc.Active {
		t.Error("account should be active on creation")
	}
	acc.Deactivate()
	if acc.Active {
		t.Error("account should be inactive after deactivation")
	}
}

func TestAccountValidation(t *testing.T) {
	tests := []struct {
		name     string
		acc      *Account
		hasError bool
	}{
		{
			"valid account",
			&Account{ID: "acc1", Name: "Checking", AccountType: "checking", Balance: 100},
			false,
		},
		{
			"missing id",
			&Account{ID: "", Name: "Checking", AccountType: "checking", Balance: 100},
			true,
		},
		{
			"missing name",
			&Account{ID: "acc1", Name: "", AccountType: "checking", Balance: 100},
			true,
		},
		{
			"missing type",
			&Account{ID: "acc1", Name: "Checking", AccountType: "", Balance: 100},
			true,
		},
	}

	for _, tt := range tests {
		err := tt.acc.Validate()
		if (err != nil) != tt.hasError {
			t.Errorf("%s: Validate() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}
