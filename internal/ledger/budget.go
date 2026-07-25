package ledger

import "time"

type Budget struct {
	ID             string     `json:"id"`
	Category       string     `json:"category"`
	MonthlyLimit   float64    `json:"monthly_limit"`
	AlertThreshold float64    `json:"alert_threshold"` // percentage of limit (e.g., 0.8 = 80%)
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	Active         bool       `json:"active"`
}

func NewBudget(category string, monthlyLimit float64) *Budget {
	return &Budget{
		ID:             generateID(),
		Category:       category,
		MonthlyLimit:   monthlyLimit,
		AlertThreshold: 0.8,
		StartDate:      time.Now(),
		Active:         true,
	}
}

func (b *Budget) Validate() error {
	if b.Category == "" {
		return NewValidationError("category is required")
	}
	if b.MonthlyLimit <= 0 {
		return NewValidationError("monthly limit must be positive")
	}
	if b.AlertThreshold < 0 || b.AlertThreshold > 1 {
		return NewValidationError("alert threshold must be between 0 and 1")
	}
	return nil
}

func (b *Budget) IsActive() bool {
	now := time.Now()
	if !b.Active {
		return false
	}
	if b.StartDate.After(now) {
		return false
	}
	if b.EndDate != nil && b.EndDate.Before(now) {
		return false
	}
	return true
}

func (b *Budget) CheckAlert(spent float64) bool {
	if b.MonthlyLimit <= 0 {
		return false
	}
	percentSpent := spent / b.MonthlyLimit
	return percentSpent >= b.AlertThreshold
}

func (b *Budget) RemainingBudget(spent float64) float64 {
	return b.MonthlyLimit - spent
}

// BudgetPeriod represents a time period for budget calculations
type BudgetPeriod struct {
	Start time.Time
	End   time.Time
}

func GetCurrentMonthPeriod() BudgetPeriod {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return BudgetPeriod{start, end}
}

func GetMonthPeriod(year int, month time.Month) BudgetPeriod {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return BudgetPeriod{start, end}
}
