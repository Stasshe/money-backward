package ledger

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Transaction struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Tags        []string  `json:"tags,omitempty"`
}

func NewTransaction(accountID, category string, amount float64, description string) *Transaction {
	return &Transaction{
		ID:          generateID(),
		AccountID:   accountID,
		Category:    category,
		Amount:      amount,
		Description: description,
		Timestamp:   time.Now(),
		Tags:        []string{},
	}
}

func (t *Transaction) IsExpense() bool {
	return t.Amount < 0
}

func (t *Transaction) IsIncome() bool {
	return t.Amount > 0
}

func (t *Transaction) AbsAmount() float64 {
	if t.Amount < 0 {
		return -t.Amount
	}
	return t.Amount
}

func (t *Transaction) AddTag(tag string) {
	for _, existing := range t.Tags {
		if existing == tag {
			return
		}
	}
	t.Tags = append(t.Tags, tag)
}

func (t *Transaction) Validate() error {
	if t.AccountID == "" {
		return NewValidationError("account_id is required")
	}
	if t.Category == "" {
		return NewValidationError("category is required")
	}
	if t.Amount == 0 {
		return NewValidationError("amount cannot be zero")
	}
	return nil
}

func generateID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// fallback: use timestamp
		return time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(b)
}

type ValidationError struct {
	msg string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{msg}
}

func (e *ValidationError) Error() string {
	return "validation error: " + e.msg
}
