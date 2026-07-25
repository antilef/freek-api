package dates

import "time"

type PeriodConfig struct {
	StartDay int
	EndDay   int
}

func (pc PeriodConfig) GetPeriodForDate(date time.Time) (start, end time.Time) {
	year, month, _ := date.Date()

	startDay := pc.StartDay
	endDay := pc.EndDay

	if startDay <= endDay {
		start = time.Date(year, month, startDay, 0, 0, 0, 0, date.Location())
		end = time.Date(year, month, endDay, 23, 59, 59, 0, date.Location())
	} else {
		if date.Day() >= startDay {
			start = time.Date(year, month, startDay, 0, 0, 0, 0, date.Location())
			end = time.Date(year, month+1, endDay, 23, 59, 59, 0, date.Location())
		} else {
			start = time.Date(year, month-1, startDay, 0, 0, 0, 0, date.Location())
			end = time.Date(year, month, endDay, 23, 59, 59, 0, date.Location())
		}
	}

	return start, end
}

func (pc PeriodConfig) GetCurrentPeriod(now time.Time) (start, end time.Time) {
	return pc.GetPeriodForDate(now)
}

func (pc PeriodConfig) GetNextPeriodStart(now time.Time) time.Time {
	_, currentEnd := pc.GetPeriodForDate(now)
	return currentEnd.Add(1 * time.Second)
}

func (pc PeriodConfig) GetPreviousPeriodStart(now time.Time) time.Time {
	start, _ := pc.GetPeriodForDate(now)
	return start.AddDate(0, -1, 0)
}

func DateRange(from, to time.Time) (start, end time.Time) {
	if from.Before(to) {
		return from, to
	}
	return to, from
}
