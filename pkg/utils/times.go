package utils

import "time"

func GetFirstDayOfNextMonth(now time.Time) int64 {
	var nextMonthFirstDay int64 = 0
	currentYear, currentMonth, _ := now.Date()
	location := now.Location()

	if currentMonth < 12 {
		nextMonthFirstDay = time.Date(currentYear, currentMonth+1, 1, 0, 0, 1, 0, location).Unix()
	} else {
		nextMonthFirstDay = time.Date(currentYear+1, 1, 1, 0, 0, 1, 0, location).Unix()
	}

	return nextMonthFirstDay
}

func GetFirstDayOfPreviousMonth(now time.Time) int64 {
	var previousMonthFirstDay int64 = 0
	currentYear, currentMonth, _ := now.Date()
	location := now.Location()

	if currentMonth == 1 {
		previousMonthFirstDay = time.Date(currentYear-1, 12, 1, 0, 0, 1, 0, location).Unix()
	} else {
		previousMonthFirstDay = time.Date(currentYear, currentMonth-1, 1, 0, 0, 1, 0, location).Unix()
	}

	return previousMonthFirstDay
}

func GetLastDayOfMonth(now time.Time) int64 {
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1).Unix()

	return lastOfMonth
}

func GetFirstDayOfMonth(now time.Time) int64 {
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation).Unix()

	return firstOfMonth
}
