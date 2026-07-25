package ledger

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	fv := NewFieldValidator()

	tests := []struct {
		name     string
		email    string
		hasError bool
	}{
		{"valid email", "user@example.com", false},
		{"empty email (optional)", "", false},
		{"invalid email no domain", "user@", true},
		{"invalid email no @", "userexample.com", true},
		{"invalid email spaces", "user @example.com", true},
		{"valid complex email", "first.last+tag@example.co.uk", false},
	}

	for _, tt := range tests {
		err := fv.ValidateEmail(tt.email)
		if (err != nil) != tt.hasError {
			t.Errorf("%s: ValidateEmail() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestValidateAmountRange(t *testing.T) {
	fv := NewFieldValidator()

	tests := []struct {
		name     string
		amount   float64
		min      float64
		max      float64
		hasError bool
	}{
		{"within range", 50.0, 0, 100, false},
		{"at min boundary", 0, 0, 100, false},
		{"at max boundary", 100, 0, 100, false},
		{"below min", -10, 0, 100, true},
		{"above max", 150, 0, 100, true},
	}

	for _, tt := range tests {
		err := fv.ValidateAmountRange(tt.amount, tt.min, tt.max)
		if (err != nil) != tt.hasError {
			t.Errorf("%s: ValidateAmountRange() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestValidatePositiveAmount(t *testing.T) {
	fv := NewFieldValidator()

	tests := []struct {
		name     string
		amount   float64
		hasError bool
	}{
		{"positive", 100.0, false},
		{"small positive", 0.01, false},
		{"zero", 0, true},
		{"negative", -50.0, true},
	}

	for _, tt := range tests {
		err := fv.ValidatePositiveAmount(tt.amount)
		if (err != nil) != tt.hasError {
			t.Errorf("%s: ValidatePositiveAmount() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestValidateStringLength(t *testing.T) {
	fv := NewFieldValidator()

	tests := []struct {
		name     string
		str      string
		min      int
		max      int
		hasError bool
	}{
		{"within length", "hello", 1, 10, false},
		{"exact min", "a", 1, 10, false},
		{"exact max", "0123456789", 1, 10, false},
		{"too short", "a", 2, 10, true},
		{"too long", "0123456789a", 1, 10, true},
	}

	for _, tt := range tests {
		err := fv.ValidateStringLength(tt.str, tt.min, tt.max)
		if (err != nil) != tt.hasError {
			t.Errorf("%s: ValidateStringLength() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestValidateNotEmpty(t *testing.T) {
	fv := NewFieldValidator()

	tests := []struct {
		name     string
		str      string
		hasError bool
	}{
		{"valid string", "hello", false},
		{"empty string", "", true},
		{"whitespace only", "   ", true},
	}

	for _, tt := range tests {
		err := fv.ValidateNotEmpty(tt.str)
		if (err != nil) != tt.hasError {
			t.Errorf("%s: ValidateNotEmpty() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestValidateOwner(t *testing.T) {
	aov := NewAccountOwnerValidator()

	tests := []struct {
		name     string
		owner    string
		hasError bool
	}{
		{"valid owner", "John Doe", false},
		{"empty owner (optional)", "", false},
		{"too short", "A", true},
		{"very long", string(make([]byte, 101)), true},
		{"normal name", "Jane Smith", false},
	}

	for _, tt := range tests {
		err := aov.ValidateOwner(tt.owner)
		if (err != nil) != tt.hasError {
			t.Errorf("%s: ValidateOwner() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestValidateOwnerEmail(t *testing.T) {
	aov := NewAccountOwnerValidator()

	tests := []struct {
		name     string
		email    string
		hasError bool
	}{
		{"valid email", "owner@company.com", false},
		{"invalid email", "not-an-email", true},
	}

	for _, tt := range tests {
		err := aov.ValidateOwnerEmail(tt.email)
		if (err != nil) != tt.hasError {
			t.Errorf("%s: ValidateOwnerEmail() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestTransactionValidator(t *testing.T) {
	tv := NewTransactionValidator()

	tests := []struct {
		name     string
		txn      *Transaction
		hasError bool
	}{
		{
			"valid transaction",
			NewTransaction("acc1", "groceries", -50.0, "shopping"),
			false,
		},
		{
			"amount too large",
			NewTransaction("acc1", "invalid", 1000000.0, "test"),
			true,
		},
		{
			"description too long",
			&Transaction{
				ID:          "t1",
				AccountID:   "acc1",
				Category:    "test",
				Amount:      100,
				Description: string(make([]byte, 501)),
			},
			true,
		},
	}

	for _, tt := range tests {
		err := tv.ValidateTransaction(tt.txn)
		if (err != nil) != tt.hasError {
			t.Errorf("%s: ValidateTransaction() error = %v, expected error = %v", tt.name, err, tt.hasError)
		}
	}
}

func TestSetAmountRange(t *testing.T) {
	tv := NewTransactionValidator()
	tv.SetAmountRange(-100, 100)

	txn := NewTransaction("acc1", "test", 150, "test")
	err := tv.ValidateTransaction(txn)
	if err == nil {
		t.Error("expected validation error for amount outside range")
	}
}
