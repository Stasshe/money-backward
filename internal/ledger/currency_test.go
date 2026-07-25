package ledger

import (
	"math"
	"testing"
)

func TestIsValidCurrency(t *testing.T) {
	tests := []struct {
		currency string
		valid    bool
	}{
		{"USD", true},
		{"EUR", true},
		{"JPY", true},
		{"GBP", true},
		{"CAD", true},
		{"XXX", false},
		{"", false},
		{"usd", false},
	}

	for _, tt := range tests {
		result := IsValidCurrency(tt.currency)
		if result != tt.valid {
			t.Errorf("IsValidCurrency(%s) = %v, expected %v", tt.currency, result, tt.valid)
		}
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		from        string
		to          string
		expectedMin float64
		expectedMax float64
		shouldError bool
	}{
		{"USD to EUR", 100, "USD", "EUR", 91, 93, false},
		{"USD to JPY", 100, "USD", "JPY", 14900, 15000, false},
		{"EUR to GBP", 50, "EUR", "GBP", 42, 44, false},
		{"same currency", 100, "USD", "USD", 99.9, 100.1, false},
		{"invalid from currency", 100, "ZZZ", "EUR", 0, 0, true},
		{"invalid to currency", 100, "USD", "ZZZ", 0, 0, true},
	}

	for _, tt := range tests {
		result, err := Convert(tt.amount, tt.from, tt.to)
		if (err != nil) != tt.shouldError {
			t.Errorf("%s: Convert() error = %v, expected error = %v", tt.name, err, tt.shouldError)
			continue
		}

		if !tt.shouldError {
			if result < tt.expectedMin || result > tt.expectedMax {
				t.Errorf("%s: Convert() = %f, expected between %f and %f", tt.name, result, tt.expectedMin, tt.expectedMax)
			}
		}
	}
}

func TestGetExchangeRate(t *testing.T) {
	tests := []struct {
		name              string
		from              string
		to                string
		expectedRateApprox float64
		shouldError       bool
	}{
		{"USD to EUR", "USD", "EUR", 0.92, false},
		{"EUR to USD", "EUR", "USD", 1.09, false},
		{"JPY to GBP", "JPY", "GBP", 0.0053, false},
		{"USD to USD", "USD", "USD", 1.0, false},
		{"invalid from", "ZZZ", "USD", 0, true},
		{"invalid to", "USD", "ZZZ", 0, true},
	}

	for _, tt := range tests {
		result, err := GetExchangeRate(tt.from, tt.to)
		if (err != nil) != tt.shouldError {
			t.Errorf("%s: GetExchangeRate() error = %v, expected error = %v", tt.name, err, tt.shouldError)
			continue
		}

		if !tt.shouldError {
			// Use tolerance for float comparison
			if math.Abs(result-tt.expectedRateApprox) > 0.01 {
				t.Errorf("%s: GetExchangeRate() = %f, expected ~%f", tt.name, result, tt.expectedRateApprox)
			}
		}
	}
}

func TestConvertRoundTrip(t *testing.T) {
	amount := 100.0
	original, _ := Convert(amount, "USD", "EUR")
	result, _ := Convert(original, "EUR", "USD")

	// Round trip should be close to original (within 1% due to rounding)
	tolerance := amount * 0.01
	if math.Abs(result-amount) > tolerance {
		t.Errorf("round-trip conversion: %f USD -> EUR -> %f USD, expected close to %f", amount, result, amount)
	}
}
