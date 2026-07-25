package ledger

type Category struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // income, expense, transfer
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

func NewCategory(name, categoryType string) *Category {
	return &Category{
		Name:   name,
		Type:   categoryType,
		Color:  "#808080",
		Icon:   "default",
		Active: true,
	}
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return NewValidationError("category name is required")
	}
	if !IsValidCategoryType(c.Type) {
		return NewValidationError("invalid category type")
	}
	return nil
}

func (c *Category) SetColor(hex string) error {
	// Basic hex validation
	if len(hex) != 7 || hex[0] != '#' {
		return NewValidationError("color must be hex format (#RRGGBB)")
	}
	c.Color = hex
	return nil
}

func IsValidCategoryType(t string) bool {
	switch t {
	case "income", "expense", "transfer":
		return true
	}
	return false
}

// DefaultCategories provides a standard set of expense categories
var DefaultCategories = []Category{
	{Name: "Groceries", Type: "expense", Color: "#2ecc71", Icon: "shopping"},
	{Name: "Utilities", Type: "expense", Color: "#3498db", Icon: "lightbulb"},
	{Name: "Transport", Type: "expense", Color: "#e74c3c", Icon: "car"},
	{Name: "Entertainment", Type: "expense", Color: "#9b59b6", Icon: "film"},
	{Name: "Healthcare", Type: "expense", Color: "#e67e22", Icon: "heart"},
	{Name: "Salary", Type: "income", Color: "#27ae60", Icon: "briefcase"},
	{Name: "Bonus", Type: "income", Color: "#2ecc71", Icon: "gift"},
	{Name: "Investment Return", Type: "income", Color: "#16a085", Icon: "chart"},
	{Name: "Freelance", Type: "income", Color: "#8e44ad", Icon: "laptop"},
	{Name: "Transfer", Type: "transfer", Color: "#95a5a6", Icon: "arrow"},
}
