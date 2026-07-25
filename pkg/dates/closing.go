package dates

import "time"

func NextClosingDate(purchaseDate time.Time, closingDay int) time.Time {
	purchaseYear, purchaseMonth, purchaseDay := purchaseDate.Date()

	var candidate time.Time
	if purchaseDay <= closingDay {
		candidate = time.Date(purchaseYear, purchaseMonth, closingDay, 0, 0, 0, 0, purchaseDate.Location())
	} else {
		candidate = time.Date(purchaseYear, purchaseMonth+1, closingDay, 0, 0, 0, 0, purchaseDate.Location())
	}

	return candidate
}

func ClosingDates(purchaseDate time.Time, closingDay int, months int) []time.Time {
	dates := make([]time.Time, months)
	first := NextClosingDate(purchaseDate, closingDay)
	for i := 0; i < months; i++ {
		dates[i] = first.AddDate(0, i, 0)
	}
	return dates
}

func SettlementForMonth(closingDay int, year int, month time.Month) (start, end time.Time) {
	end = time.Date(year, month, closingDay, 23, 59, 59, 0, time.UTC)

	var startMonth time.Month
	startYear := year
	if month == time.January {
		startMonth = time.December
		startYear = year - 1
	} else {
		startMonth = month - 1
	}

	startDay := closingDay + 1
	if startDay > 28 {
		startDay = 1
		startMonth++
		if startMonth > time.December {
			startMonth = time.January
			startYear++
		}
	}

	start = time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, time.UTC)
	return start, end
}
