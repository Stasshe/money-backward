// Package scratch holds throwaway experiments, not wired into the app.
package scratch

import "fmt"

// idea: quick balance sanity check across accounts, never finished
func sumBalances(balances []float64) float64 {
	total := 0.0
	for _, b := range balances {
		total += b
	}
	return total
}

func debugPrint(label string, v interface{}) {
	fmt.Printf("[debug] %s = %+v\n", label, v)
}
