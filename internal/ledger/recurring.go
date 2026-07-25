package ledger

import (
	"fmt"
	"time"
)

type RecurringTransaction struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Frequency   string    `json:"frequency"` // daily, weekly, biweekly, monthly, quarterly, yearly
	LastApplied *time.Time `json:"last_applied,omitempty"`
	Active      bool      `json:"active"`
}

func NewRecurringTransaction(accountID, category string, amount float64, description string, frequency string) *RecurringTransaction {
	return &RecurringTransaction{
		ID:          generateID(),
		AccountID:   accountID,
		Category:    category,
		Amount:      amount,
		Description: description,
		StartDate:   time.Now(),
		Frequency:   frequency,
		Active:      true,
	}
}

func (r *RecurringTransaction) NextOccurrence() (time.Time, error) {
	base := r.StartDate
	if r.LastApplied != nil {
		base = *r.LastApplied
	}

	var next time.Time
	switch r.Frequency {
	case "daily":
		next = base.AddDate(0, 0, 1)
	case "weekly":
		next = base.AddDate(0, 0, 7)
	case "biweekly":
		next = base.AddDate(0, 0, 14)
	case "monthly":
		next = base.AddDate(0, 1, 0)
	case "quarterly":
		next = base.AddDate(0, 3, 0)
	case "yearly":
		next = base.AddDate(1, 0, 0)
	default:
		return time.Time{}, fmt.Errorf("invalid frequency: %s", r.Frequency)
	}

	if r.EndDate != nil && next.After(*r.EndDate) {
		return time.Time{}, fmt.Errorf("next occurrence after end date")
	}

	return next, nil
}

func (r *RecurringTransaction) Apply() *Transaction {
	txn := &Transaction{
		ID:          generateID(),
		AccountID:   r.AccountID,
		Category:    r.Category,
		Amount:      r.Amount,
		Description: r.Description,
		Timestamp:   time.Now(),
		Tags:        []string{fmt.Sprintf("recurring:%s", r.ID)},
	}
	now := time.Now()
	r.LastApplied = &now
	return txn
}

func (r *RecurringTransaction) Validate() error {
	if r.AccountID == "" {
		return NewValidationError("account_id is required")
	}
	if r.Category == "" {
		return NewValidationError("category is required")
	}
	if r.Amount == 0 {
		return NewValidationError("amount cannot be zero")
	}
	if r.Frequency == "" {
		return NewValidationError("frequency is required")
	}
	return validateFrequency(r.Frequency)
}

func validateFrequency(freq string) error {
	switch freq {
	case "daily", "weekly", "biweekly", "monthly", "quarterly", "yearly":
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid frequency: %s", freq))
	}
}

func (r *RecurringTransaction) Deactivate() {
	r.Active = false
}

func (r *RecurringTransaction) IsActive() bool {
	if !r.Active {
		return false
	}
	if r.EndDate != nil && time.Now().After(*r.EndDate) {
		return false
	}
	return true
}
