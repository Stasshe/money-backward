package ledger

type Account struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	AccountType   string  `json:"account_type"` // checking, savings, credit_card, investment
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
	Active        bool    `json:"active"`
	LastModified  string  `json:"last_modified"`
	InitialAmount float64 `json:"initial_amount"`
}

func NewAccount(id, name, accountType string, initialBalance float64) *Account {
	return &Account{
		ID:            id,
		Name:          name,
		AccountType:   accountType,
		Balance:       initialBalance,
		Currency:      "USD",
		Active:        true,
		InitialAmount: initialBalance,
		LastModified:  getCurrentTime(),
	}
}

func (a *Account) Credit(amount float64) error {
	if amount <= 0 {
		return NewValidationError("credit amount must be positive")
	}
	a.Balance += amount
	a.LastModified = getCurrentTime()
	return nil
}

func (a *Account) Debit(amount float64) error {
	if amount <= 0 {
		return NewValidationError("debit amount must be positive")
	}
	// TODO: implement overdraft protection / warning
	a.Balance -= amount
	a.LastModified = getCurrentTime()
	return nil
}

func (a *Account) GetBalance() float64 {
	return a.Balance
}

func (a *Account) SetBalance(balance float64) {
	a.Balance = balance
	a.LastModified = getCurrentTime()
}

func (a *Account) Deactivate() {
	a.Active = false
	a.LastModified = getCurrentTime()
}

func (a *Account) Validate() error {
	if a.ID == "" {
		return NewValidationError("account id is required")
	}
	if a.Name == "" {
		return NewValidationError("account name is required")
	}
	if a.AccountType == "" {
		return NewValidationError("account type is required")
	}
	return nil
}

type AccountType string

const (
	Checking   AccountType = "checking"
	Savings    AccountType = "savings"
	CreditCard AccountType = "credit_card"
	Investment AccountType = "investment"
	Cash       AccountType = "cash"
)

func IsValidAccountType(t string) bool {
	switch t {
	case "checking", "savings", "credit_card", "investment", "cash":
		return true
	}
	return false
}
