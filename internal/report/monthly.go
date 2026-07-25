package report

import (
	"fmt"
	"time"
)

type MonthlyStats struct {
	Year           int
	Month          time.Month
	Income         float64
	Expenses       float64
	NetCashFlow    float64
	AccountCount   int
	TransactionCnt int
}

func NewMonthlyStats(year int, month time.Month) *MonthlyStats {
	return &MonthlyStats{
		Year:  year,
		Month: month,
	}
}

func (m *MonthlyStats) SavingsRate() float64 {
	if m.Income == 0 {
		return 0
	}
	return (m.NetCashFlow / m.Income) * 100
}

func (m *MonthlyStats) IsDeficit() bool {
	return m.NetCashFlow < 0
}

func (m *MonthlyStats) FormattedMonthYear() string {
	return fmt.Sprintf("%s %d", m.Month.String(), m.Year)
}

type MonthlyComparison struct {
	CurrentMonth  MonthlyStats
	PreviousMonth MonthlyStats
}

func (mc *MonthlyComparison) IncomeChange() float64 {
	return mc.CurrentMonth.Income - mc.PreviousMonth.Income
}

func (mc *MonthlyComparison) ExpenseChange() float64 {
	return mc.CurrentMonth.Expenses - mc.PreviousMonth.Expenses
}

func (mc *MonthlyComparison) NetChange() float64 {
	return mc.CurrentMonth.NetCashFlow - mc.PreviousMonth.NetCashFlow
}

func (mc *MonthlyComparison) IncomeChangePercent() float64 {
	if mc.PreviousMonth.Income == 0 {
		return 0
	}
	return (mc.IncomeChange() / mc.PreviousMonth.Income) * 100
}

type QuarterlyAnalysis struct {
	Q1   MonthlyStats
	Q2   MonthlyStats
	Q3   MonthlyStats
	Q4   MonthlyStats
	Year int
}

func (qa *QuarterlyAnalysis) TotalIncome() float64 {
	return qa.Q1.Income + qa.Q2.Income + qa.Q3.Income + qa.Q4.Income
}

func (qa *QuarterlyAnalysis) TotalExpenses() float64 {
	return qa.Q1.Expenses + qa.Q2.Expenses + qa.Q3.Expenses + qa.Q4.Expenses
}

func (qa *QuarterlyAnalysis) AverageSavingsRate() float64 {
	totalRate := qa.Q1.SavingsRate() + qa.Q2.SavingsRate() + qa.Q3.SavingsRate() + qa.Q4.SavingsRate()
	return totalRate / 4
}

type MonthlyForecast struct {
	BasedOnMonth     MonthlyStats
	ProjectedIncome  float64
	ProjectedExpense float64
}

func (mf *MonthlyForecast) ProjectedSavings() float64 {
	return mf.ProjectedIncome - mf.ProjectedExpense
}

// TODO: implement actual forecasting based on historical trends
func GenerateForecast(history []MonthlyStats) *MonthlyForecast {
	if len(history) == 0 {
		return nil
	}

	lastMonth := history[len(history)-1]
	forecast := &MonthlyForecast{
		BasedOnMonth: lastMonth,
	}

	// Simple placeholder: use last month as forecast
	forecast.ProjectedIncome = lastMonth.Income
	forecast.ProjectedExpense = lastMonth.Expenses

	return forecast
}
