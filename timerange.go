// Copyright © 2016 Great Ideas Club. All rights reserved.

package timerange

import (
	"sort"
	"time"
)

// TimeRange describes a start and end time as [start, end)
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Parse creates a new time range given the format and start end times
func Parse(format, start, end string) (TimeRange, error) {
	startTime, err := time.Parse(format, start)
	if err != nil {
		return TimeRange{}, err
	}
	endTime, err := time.Parse(format, end)
	if err != nil {
		return TimeRange{}, err
	}
	return TimeRange{Start: startTime, End: endTime}, nil
}

// Duration returns the difference between the start and end time
func (tr *TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// timeRangeSlice is a helper type for sorting a slice of time ranges
type timeRangeSlice []TimeRange

// Len returns the number of TimeRange objects in the slice
func (trs timeRangeSlice) Len() int {
	return len(trs)
}

// Less returns true if index i Start precedes index j Start or the indices
// have the same start time but index i End precedes index j End
func (trs timeRangeSlice) Less(i, j int) bool {
	if trs[i].Start.Equal(trs[j].Start) {
		return trs[i].End.Before(trs[j].End)
	}
	return trs[i].Start.Before(trs[j].Start)
}

// Swap swaps the values between the two indices
func (trs timeRangeSlice) Swap(i, j int) {
	trs[i], trs[j] = trs[j], trs[i]
}

// Sort sorts a slice of TimeRange objects by start and end time.
func Sort(trs []TimeRange) {
	sort.Sort(timeRangeSlice(trs))
}

// Merge merges overlapping TimeRange objects
func Merge(trs *[]TimeRange) {
	var ranges = *trs

	Sort(ranges)

	var index = 0
	for index < len(ranges)-1 {
		if ranges[index].End.Before(ranges[index+1].Start) {
			index++
		} else {
			if ranges[index].End.Before(ranges[index+1].End) {
				ranges[index].End = ranges[index+1].End
			}
			ranges = append(ranges[:index+1], ranges[index+2:]...)
		}
	}

	*trs = ranges
}

// SearchIndex returns the index of the first time range that satisfies the
// boundaries of a provided time range.  Returns -1 if not found.
func SearchIndex(trs []TimeRange, v TimeRange) int {
	for i, tr := range trs {
		if !v.Start.Before(tr.Start) && !v.End.After(tr.End) {
			return i
		}
	}

	return -1
}
