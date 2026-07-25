package ledger

import "fmt"

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	JPY Currency = "JPY"
	GBP Currency = "GBP"
	CAD Currency = "CAD"
)

var exchangeRates = map[string]map[string]float64{
	"USD": {
		"EUR": 0.92,
		"JPY": 149.50,
		"GBP": 0.79,
		"CAD": 1.36,
		"USD": 1.0,
	},
	"EUR": {
		"USD": 1.09,
		"JPY": 162.50,
		"GBP": 0.86,
		"CAD": 1.48,
		"EUR": 1.0,
	},
	"JPY": {
		"USD": 0.0067,
		"EUR": 0.0062,
		"GBP": 0.0053,
		"CAD": 0.0091,
		"JPY": 1.0,
	},
	"GBP": {
		"USD": 1.27,
		"EUR": 1.16,
		"JPY": 188.70,
		"CAD": 1.72,
		"GBP": 1.0,
	},
	"CAD": {
		"USD": 0.74,
		"EUR": 0.68,
		"JPY": 110.00,
		"GBP": 0.58,
		"CAD": 1.0,
	},
}

func IsValidCurrency(c string) bool {
	_, ok := exchangeRates[c]
	return ok
}

func Convert(amount float64, from, to string) (float64, error) {
	if !IsValidCurrency(from) {
		return 0, fmt.Errorf("invalid source currency: %s", from)
	}
	if !IsValidCurrency(to) {
		return 0, fmt.Errorf("invalid target currency: %s", to)
	}

	rate, ok := exchangeRates[from][to]
	if !ok {
		return 0, fmt.Errorf("no exchange rate found from %s to %s", from, to)
	}

	return amount * rate, nil
}

func GetExchangeRate(from, to string) (float64, error) {
	if !IsValidCurrency(from) {
		return 0, fmt.Errorf("invalid source currency: %s", from)
	}
	if !IsValidCurrency(to) {
		return 0, fmt.Errorf("invalid target currency: %s", to)
	}

	rate, ok := exchangeRates[from][to]
	if !ok {
		return 0, fmt.Errorf("no exchange rate found from %s to %s", from, to)
	}

	return rate, nil
}
