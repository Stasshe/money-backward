package report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"sort"

	"money-backword/internal/storage"
)

type CSVExporter struct {
	store storage.Store
}

func NewCSVExporter(store storage.Store) *CSVExporter {
	return &CSVExporter{store: store}
}

func (e *CSVExporter) ExportTransactions(accountID string, limit int) (string, error) {
	txns, err := e.store.GetTransactions(accountID, limit)
	if err != nil {
		return "", fmt.Errorf("failed to fetch transactions: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Account", "Category", "Amount", "Description", "Timestamp"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Write transaction rows
	for _, txn := range txns {
		row := []string{
			txn.ID,
			txn.AccountID,
			txn.Category,
			fmt.Sprintf("%.2f", txn.Amount),
			txn.Description,
			txn.Timestamp.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	return buf.String(), nil
}

func (e *CSVExporter) ExportMonthlyReport(month string) (string, error) {
	monthlyReport, err := e.store.GenerateMonthlyReport(month)
	if err != nil {
		return "", fmt.Errorf("failed to generate monthly report: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Metadata section
	writer.Write([]string{"Monthly Report", month})
	writer.Write([]string{})
	writer.Write([]string{"Summary"})
	writer.Write([]string{"Total Income", fmt.Sprintf("%.2f", monthlyReport.TotalIncome)})
	writer.Write([]string{"Total Expense", fmt.Sprintf("%.2f", monthlyReport.TotalExpense)})
	writer.Write([]string{"Net", fmt.Sprintf("%.2f", monthlyReport.Net)})
	writer.Write([]string{"Transaction Count", fmt.Sprintf("%d", monthlyReport.TransactionCount)})
	writer.Write([]string{})
	writer.Write([]string{"Category Breakdown"})
	writer.Write([]string{"Category", "Amount"})

	// Sort categories by name for consistent output
	categories := make([]string, 0, len(monthlyReport.CategoryBreakdown))
	for cat := range monthlyReport.CategoryBreakdown {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, cat := range categories {
		writer.Write([]string{cat, fmt.Sprintf("%.2f", monthlyReport.CategoryBreakdown[cat])})
	}

	writer.Flush()
	return buf.String(), nil
}

func (e *CSVExporter) ExportAccountSummary() (string, error) {
	accounts, err := e.store.GetAllAccounts()
	if err != nil {
		return "", fmt.Errorf("failed to fetch accounts: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	writer.Write([]string{"Account ID", "Name", "Type", "Balance", "Currency", "Active"})

	for _, acc := range accounts {
		row := []string{
			acc.ID,
			acc.Name,
			acc.AccountType,
			fmt.Sprintf("%.2f", acc.Balance),
			acc.Currency,
			fmt.Sprintf("%v", acc.Active),
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	return buf.String(), nil
}

func (e *CSVExporter) ExportCategoryAnalysis(month string) (string, error) {
	monthlyReport, err := e.store.GenerateMonthlyReport(month)
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	writer.Write([]string{"Category", "Amount", "Percentage"})

	totalExpense := monthlyReport.TotalExpense
	if totalExpense == 0 {
		return buf.String(), nil
	}

	categories := make([]string, 0, len(monthlyReport.CategoryBreakdown))
	for cat := range monthlyReport.CategoryBreakdown {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, cat := range categories {
		amount := monthlyReport.CategoryBreakdown[cat]
		percent := (amount / totalExpense) * 100
		row := []string{cat, fmt.Sprintf("%.2f", amount), fmt.Sprintf("%.1f%%", percent)}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	return buf.String(), nil
}

func (e *CSVExporter) ExportBudgetStatus() (string, error) {
	budgets, err := e.store.GetAllBudgets()
	if err != nil {
		return "", fmt.Errorf("failed to fetch budgets: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	writer.Write([]string{"Category", "Monthly Limit", "Status"})

	for _, budget := range budgets {
		// TODO: calculate actual spending this month
		status := "OK"
		row := []string{budget.Category, fmt.Sprintf("%.2f", budget.MonthlyLimit), status}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	return buf.String(), nil
}
