package report

import (
	"fmt"
	"sort"

	"money-backword/internal/ledger"
	"money-backword/internal/storage"
)

type Reporter struct {
	store storage.Store
}

func NewReporter(store storage.Store) *Reporter {
	return &Reporter{store: store}
}

type SummaryReport struct {
	Period          string
	TotalIncome     float64
	TotalExpense    float64
	Net             float64
	TransactionInfo TransactionInfo
}

type TransactionInfo struct {
	Count    int
	Average  float64
	Largest  float64
	Smallest float64
}

func (r *Reporter) GenerateSummary(month string) (*SummaryReport, error) {
	monthlyReport, err := r.store.GenerateMonthlyReport(month)
	if err != nil {
		return nil, err
	}

	summary := &SummaryReport{
		Period:       monthlyReport.Month,
		TotalIncome:  monthlyReport.TotalIncome,
		TotalExpense: monthlyReport.TotalExpense,
		Net:          monthlyReport.Net,
		TransactionInfo: TransactionInfo{
			Count:   monthlyReport.TransactionCount,
			Average: monthlyReport.AverageTransaction,
		},
	}

	// TODO: calculate largest and smallest transactions
	return summary, nil
}

type CategoryReport struct {
	Category string
	Count    int
	Total    float64
	Average  float64
	Percent  float64
}

func (r *Reporter) CategoryAnalysis(month string) ([]CategoryReport, error) {
	monthlyReport, err := r.store.GenerateMonthlyReport(month)
	if err != nil {
		return nil, err
	}

	var report []CategoryReport
	totalExpense := monthlyReport.TotalExpense

	for category, amount := range monthlyReport.CategoryBreakdown {
		var percent float64
		if totalExpense > 0 {
			percent = (amount / totalExpense) * 100
		}

		cr := CategoryReport{
			Category: category,
			Total:    amount,
			Percent:  percent,
		}

		// TODO: calculate count and average per category
		report = append(report, cr)
	}

	// Sort by total descending
	sort.Slice(report, func(i, j int) bool {
		return report[i].Total > report[j].Total
	})

	return report, nil
}

func (r *Reporter) FormatSummary(summary *SummaryReport) string {
	return fmt.Sprintf(
		"Financial Summary (%s)\n"+
			"Income:    $%.2f\n"+
			"Expense:   $%.2f\n"+
			"Net:       $%.2f\n"+
			"Transactions: %d (avg $%.2f)\n",
		summary.Period,
		summary.TotalIncome,
		summary.TotalExpense,
		summary.Net,
		summary.TransactionInfo.Count,
		summary.TransactionInfo.Average,
	)
}

type TrendAnalysis struct {
	PreviousMonth float64
	CurrentMonth  float64
	Change        float64
	ChangePercent float64
}

func (t *TrendAnalysis) IsTrendingUp() bool {
	return t.Change > 0
}

// TODO: implement CSV export
func (r *Reporter) ExportAsCSV(month string) (string, error) {
	// Placeholder for CSV export functionality
	return "", fmt.Errorf("CSV export not yet implemented")
}
