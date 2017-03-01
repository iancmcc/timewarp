package timewarp

import "time"

// Year returns a query that returns a range of time that exists in the
// provided year.
func Year(year int) Query {
	return func(input TimeRange) *TimeRange {
		var start, end time.Time
		if startYear := input.Start.Year(); startYear == year {
			start = input.Start
		} else if startYear < year {
			start = time.Date(year, 1, 1, 0, 0, 0, 0, input.Start.Location())
		} else {
			return nil
		}

		if endYear := input.End.Year(); endYear == year {
			end = input.End
		} else if endYear > year {
			end = time.Date(year+1, 1, 1, 0, 0, 0, 0, input.Start.Location())
		} else {
			return nil
		}

		return &TimeRange{Start: start, End: end}
	}
}

// Month returns a query that returns a range of time that exists in the
// provided month.
func Month(month time.Month) Query {
	return func(input TimeRange) *TimeRange {
		var (
			start time.Time
			end   time.Time
			delta int
		)

		if month > 0 {
			delta = getMonthDelta(input.Start.Month(), month)
		}
		if delta > 0 {
			start = input.Start.AddDate(0, delta, 1-input.Start.Day()).Truncate(24 * time.Hour)
		} else {
			start = input.Start
		}
		if !start.Before(input.End) {
			return nil
		}

		if end = start.AddDate(0, 1, 1-start.Day()).Truncate(24 * time.Hour); end.After(input.End) {
			end = input.End
		}

		return &TimeRange{Start: start, End: end}
	}
}

// TheMonth returns the full month that matches the range.  If zero, it
// returns the current month.
func TheMonth(month time.Month) Query {
	return func(input TimeRange) *TimeRange {
		var (
			start time.Time
			end   time.Time
			delta int
		)

		if month > 0 {
			delta = getMonthDelta(input.Start.Month(), month)
		}
		start = input.Start.AddDate(0, delta, 1-input.Start.Day()).Truncate(24 * time.Hour)
		end = start.AddDate(0, 1, 0)
		if start.Before(input.End) && end.After(input.Start) {
			return &TimeRange{Start: start, End: end}
		}
		return nil
	}
}

// Week finds the time ranges for n consecutive days that start on the given
// day.
func Week(weekday time.Weekday, days int) Query {
	return func(input TimeRange) *TimeRange {
		var (
			start time.Time
			end   time.Time
			delta int
		)

		if weekday >= 0 {
			delta = getWeekdayDelta(input.Start.Weekday(), weekday)
		}
		if delta > 0 {
			start = input.Start.AddDate(0, 0, delta).Truncate(24 * time.Hour)
		} else {
			start = input.Start
		}
		if !start.Before(input.End) {
			return nil
		}

		if end = start.AddDate(0, 0, days).Truncate(24 * time.Hour); end.After(input.End) {
			end = input.End
		}
		return &TimeRange{Start: start, End: end}
	}
}

// TheWeek finds the time ranges for n consecutive days that start on the
// given day from the provided offset.  If Weekday < 0.  The start day is the
// current day.
func TheWeek(weekday time.Weekday, days, offset, n int) Query {
	return func(input TimeRange) *TimeRange {
		var (
			start time.Time
			end   time.Time
			delta int
		)

		if weekday >= 0 {
			delta = getWeekdayDelta(input.Start.Weekday(), weekday)
		}
		start = input.Start.AddDate(0, 0, delta+days*offset).Truncate(24 * time.Hour)
		end = start.AddDate(0, 0, days*n)
		if !start.Before(input.End) || !end.After(input.Start) {
			return nil
		}
		return &TimeRange{Start: start, End: end}
	}
}

// Days finds the time ranges for n consecutive days from the provided input.
func Days(offset, n int) Query {
	return func(input TimeRange) *TimeRange {
		var start, end time.Time
		if offset > 0 {
			start = input.Start.AddDate(0, 0, offset).Truncate(24 * time.Hour)
		} else {
			start = input.Start
			n += offset
		}

		if !start.Before(input.End) {
			return nil
		}

		if end = start.AddDate(0, 0, n).Truncate(24 * time.Hour); !end.After(input.Start) {
			return nil
		} else if end.After(input.End) {
			end = input.End
		}

		return &TimeRange{Start: start, End: end}
	}
}

// TheDays returns the full time range from the given input day.
func TheDays(offset, n int) Query {
	return func(input TimeRange) *TimeRange {
		var (
			start = input.Start.AddDate(0, 0, offset).Truncate(24 * time.Hour)
			end   = start.AddDate(0, 0, n)
		)

		if !start.Before(input.End) || !end.After(input.Start) {
			return nil
		}
		return &TimeRange{Start: start, End: end}
	}
}

// Times returns the time that suits the timerange
func Times(format, from, to string) Query {
	var (
		fromTime, _ = time.Parse(format, from)
		toTime, _   = time.Parse(format, to)
		deltaTime   = time.Duration(mod(int(toTime.Sub(fromTime)), int(24*time.Hour)))
	)

	return func(input TimeRange) *TimeRange {
		var start, end time.Time

		delta := time.Duration(fromTime.Hour()-input.Start.Hour())*time.Hour + time.Duration(fromTime.Minute()-input.Start.Minute())*time.Minute

		start = input.Start.Add(delta).Add(-24 * time.Hour)
		end = start.Add(deltaTime)

		for !end.After(input.Start) {
			start = start.Add(24 * time.Hour)
			end = start.Add(deltaTime)
		}

		if start.Before(input.Start) {
			start = input.Start
		}

		if !start.Before(input.End) {
			return nil
		}

		if end.After(input.End) {
			end = input.End
		}

		return &TimeRange{Start: start, End: end}
	}
}

// mod performs modulo, but adds the divisor value if negative
func mod(x, y int) int {
	return (x%y + y) % y
}

// getMonthDelta returns the number of months between two months of the year
func getMonthDelta(from, to time.Month) int {
	return mod(int(to-from), 12)
}

// getWeekdayDelta returns the number of days between two days of the week
func getWeekdayDelta(from, to time.Weekday) int {
	return mod(int(to-from), 7)
}
