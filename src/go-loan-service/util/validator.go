package util

import "time"

func ValidateDate(date time.Time) bool {
	// Check if the year is valid (between 0001 and 9999)
	if date.Year() < 1 || date.Year() > 9999 {
		return false
	}

	// Check if the month is valid (between 1 and 12)
	if date.Month() < 1 || date.Month() > 12 {
		return false
	}

	// Check if the day is valid for the given month and year
	if date.Day() < 1 || date.Day() > daysInMonth(date.Year(), int(date.Month())) {
		return false
	}

	// The time is always valid
	return true
}

func daysInMonth(year int, month int) int {
	// Create a time object in the UTC timezone for the first day of the next month
	nextMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	// Subtract one day to get the last day of the current month
	lastDayOfMonth := nextMonth.Add(-24 * time.Hour).Day()
	return lastDayOfMonth
}
