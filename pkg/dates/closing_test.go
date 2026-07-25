package dates

import (
	"testing"
	"time"
)

func TestNextClosingDate_BeforeClosing(t *testing.T) {
	purchase := time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC)
	closingDay := 15

	result := NextClosingDate(purchase, closingDay)

	if result.Day() != 15 || result.Month() != 7 {
		t.Errorf("NextClosingDate = %v, want day 15 month 7", result)
	}
}

func TestNextClosingDate_AfterClosing(t *testing.T) {
	purchase := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	closingDay := 15

	result := NextClosingDate(purchase, closingDay)

	if result.Day() != 15 || result.Month() != 8 {
		t.Errorf("NextClosingDate = %v, want day 15 month 8", result)
	}
}

func TestNextClosingDate_SameAsClosing(t *testing.T) {
	purchase := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	closingDay := 15

	result := NextClosingDate(purchase, closingDay)

	if result.Day() != 15 || result.Month() != 7 {
		t.Errorf("NextClosingDate = %v, want day 15 month 7", result)
	}
}

func TestClosingDates(t *testing.T) {
	purchase := time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC)
	closingDay := 15
	months := 6

	dates := ClosingDates(purchase, closingDay, months)

	if len(dates) != months {
		t.Errorf("len(dates) = %d, want %d", len(dates), months)
	}

	for i, d := range dates {
		expectedMonth := time.Month(7 + i)
		if d.Day() != 15 {
			t.Errorf("dates[%d].Day() = %d, want 15", i, d.Day())
		}
		if d.Month() != expectedMonth {
			t.Errorf("dates[%d].Month() = %v, want %v", i, d.Month(), expectedMonth)
		}
	}
}

func TestSettlementForMonth(t *testing.T) {
	closingDay := 15
	year := 2026
	month := time.August

	start, end := SettlementForMonth(closingDay, year, month)

	if start.Day() != 16 || start.Month() != 7 {
		t.Errorf("start = %v, want day 16 month 7", start)
	}
	if end.Day() != 15 || end.Month() != 8 {
		t.Errorf("end = %v, want day 15 month 8", end)
	}
}

func TestSettlementForMonth_January(t *testing.T) {
	closingDay := 15
	year := 2026
	month := time.January

	start, end := SettlementForMonth(closingDay, year, month)

	if start.Day() != 16 || start.Month() != 12 || start.Year() != 2025 {
		t.Errorf("start = %v, want day 16 month 12 year 2025", start)
	}
	if end.Day() != 15 || end.Month() != 1 {
		t.Errorf("end = %v, want day 15 month 1", end)
	}
}

func TestSettlementForMonth_Day28(t *testing.T) {
	closingDay := 28
	year := 2026
	month := time.August

	start, end := SettlementForMonth(closingDay, year, month)

	// When closing_day is 28, start_day is 29 (28+1)
	// But 29 > 28, so it wraps to next month: start = Aug 1
	if start.Month() != time.August || start.Day() != 1 {
		t.Errorf("start = %v, want day 1 month 8", start)
	}
	if end.Day() != 28 || end.Month() != 8 {
		t.Errorf("end = %v, want day 28 month 8", end)
	}
}
