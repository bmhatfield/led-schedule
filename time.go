package main

import "time"

func Next(t time.Time) time.Time {
	if time.Now().After(t) {
		// Tomorrow
		return t.AddDate(0, 0, 1)
	}

	return t
}

func Clock(hour, min, sec int) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
}
