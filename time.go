package main

import "time"

const TimeOfDay = "3:04:05 PM"

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

func ParseClock(t, tz string) (time.Time, error) {
	now := time.Now()

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return now, err
	}

	x, err := time.ParseInLocation(TimeOfDay, t, loc)
	if err != nil {
		return now, err
	}

	return x.AddDate(now.Year(), int(now.Month())-1, now.Day()-1), nil
}
