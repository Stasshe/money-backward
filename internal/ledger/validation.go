package ledger

import (
	"fmt"
	"regexp"
	"strings"
)

type FieldValidator struct {
	// Stub field validator for extended validation rules
}

func NewFieldValidator() *FieldValidator {
	return &FieldValidator{}
}

func (fv *FieldValidator) ValidateEmail(email string) error {
	if email == "" {
		return nil // Email is optional
	}

	// Simple email-like validation (not comprehensive RFC 5322)
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return NewValidationError("email validation error")
	}

	if !matched {
		return NewValidationError("invalid email format")
	}

	return nil
}

func (fv *FieldValidator) ValidateAmountRange(amount float64, minInclusive, maxInclusive float64) error {
	if amount < minInclusive || amount > maxInclusive {
		return NewValidationError(
			fmt.Sprintf("amount %.2f outside range [%.2f, %.2f]", amount, minInclusive, maxInclusive),
		)
	}
	return nil
}

func (fv *FieldValidator) ValidatePositiveAmount(amount float64) error {
	if amount <= 0 {
		return NewValidationError("amount must be positive")
	}
	return nil
}

func (fv *FieldValidator) ValidateNonNegativeAmount(amount float64) error {
	if amount < 0 {
		return NewValidationError("amount cannot be negative")
	}
	return nil
}

func (fv *FieldValidator) ValidateStringLength(s string, minLen, maxLen int) error {
	length := len(s)
	if length < minLen || length > maxLen {
		return NewValidationError(
			fmt.Sprintf("string length %d outside range [%d, %d]", length, minLen, maxLen),
		)
	}
	return nil
}

func (fv *FieldValidator) ValidateNotEmpty(s string) error {
	if strings.TrimSpace(s) == "" {
		return NewValidationError("field cannot be empty")
	}
	return nil
}

type AccountOwnerValidator struct {
	// Validator for account owner fields (e.g., email, phone)
}

func NewAccountOwnerValidator() *AccountOwnerValidator {
	return &AccountOwnerValidator{}
}

func (aov *AccountOwnerValidator) ValidateOwner(owner string) error {
	if owner == "" {
		return nil // Owner is optional
	}

	if len(owner) < 2 {
		return NewValidationError("owner name too short")
	}

	if len(owner) > 100 {
		return NewValidationError("owner name too long")
	}

	return nil
}

func (aov *AccountOwnerValidator) ValidateOwnerEmail(email string) error {
	fv := NewFieldValidator()
	if err := fv.ValidateEmail(email); err != nil {
		return err
	}
	return nil
}

type TransactionValidator struct {
	minAmount float64
	maxAmount float64
}

func NewTransactionValidator() *TransactionValidator {
	return &TransactionValidator{
		minAmount: -999999.99,
		maxAmount: 999999.99,
	}
}

func (tv *TransactionValidator) ValidateTransaction(txn *Transaction) error {
	if err := txn.Validate(); err != nil {
		return err
	}

	// Additional range validation
	if txn.Amount < tv.minAmount || txn.Amount > tv.maxAmount {
		return NewValidationError(
			fmt.Sprintf("amount %.2f outside valid range [%.2f, %.2f]", txn.Amount, tv.minAmount, tv.maxAmount),
		)
	}

	// Description can't be too long
	if len(txn.Description) > 500 {
		return NewValidationError("description too long (max 500 characters)")
	}

	return nil
}

func (tv *TransactionValidator) SetAmountRange(min, max float64) {
	tv.minAmount = min
	tv.maxAmount = max
}
