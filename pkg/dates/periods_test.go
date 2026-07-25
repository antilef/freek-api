package dates

import (
	"testing"
	"time"
)

func TestGetPeriodForDate_SameMonth(t *testing.T) {
	pc := PeriodConfig{StartDay: 1, EndDay: 30}
	date := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)

	start, end := pc.GetPeriodForDate(date)

	if start.Day() != 1 || start.Month() != 7 {
		t.Errorf("start = %v, want day 1 month 7", start)
	}
	if end.Day() != 30 || end.Month() != 7 {
		t.Errorf("end = %v, want day 30 month 7", end)
	}
}

func TestGetPeriodForDate_CrossMonth(t *testing.T) {
	pc := PeriodConfig{StartDay: 5, EndDay: 4}
	date := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)

	start, end := pc.GetPeriodForDate(date)

	if start.Day() != 5 || start.Month() != 7 {
		t.Errorf("start = %v, want day 5 month 7", start)
	}
	if end.Day() != 4 || end.Month() != 8 {
		t.Errorf("end = %v, want day 4 month 8", end)
	}
}

func TestGetPeriodForDate_CrossMonthBeforeStart(t *testing.T) {
	pc := PeriodConfig{StartDay: 5, EndDay: 4}
	date := time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC)

	start, end := pc.GetPeriodForDate(date)

	if start.Day() != 5 || start.Month() != 6 {
		t.Errorf("start = %v, want day 5 month 6", start)
	}
	if end.Day() != 4 || end.Month() != 7 {
		t.Errorf("end = %v, want day 4 month 7", end)
	}
}

func TestGetCurrentPeriod(t *testing.T) {
	pc := PeriodConfig{StartDay: 1, EndDay: 31}
	now := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)

	start, end := pc.GetCurrentPeriod(now)

	if start.Day() != 1 || end.Day() != 31 {
		t.Errorf("period = %v to %v, want 1 to 31", start.Day(), end.Day())
	}
}

func TestGetNextPeriodStart(t *testing.T) {
	pc := PeriodConfig{StartDay: 1, EndDay: 30}
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)

	next := pc.GetNextPeriodStart(now)

	// Next period starts at 30 23:59:59 + 1 second = 31 00:00:00
	if next.Day() != 31 || next.Month() != 7 {
		t.Errorf("next period start = %v, want day 31 month 7", next)
	}
}

func TestDateRange(t *testing.T) {
	from := time.Date(2026, 7, 19, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC)

	start, end := DateRange(from, to)

	if !start.Equal(to) || !end.Equal(from) {
		t.Errorf("DateRange reversed: start=%v, end=%v, want start=%v, end=%v", start, end, to, from)
	}
}

func TestDateRange_Normal(t *testing.T) {
	from := time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 7, 19, 0, 0, 0, 0, time.UTC)

	start, end := DateRange(from, to)

	if !start.Equal(from) || !end.Equal(to) {
		t.Errorf("DateRange: start=%v, end=%v, want start=%v, end=%v", start, end, from, to)
	}
}
