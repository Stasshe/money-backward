package ledger

import "time"

func getCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func IsCurrentMonth(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month()
}

func GetMonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func GetMonthEnd(t time.Time) time.Time {
	return GetMonthStart(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}
